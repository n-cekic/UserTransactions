package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	L "userTransactions/logging"

	"github.com/IBM/sarama"
	"github.com/nats-io/nats.go"
)

type Service struct {
	mux      *http.ServeMux
	repo     Repo
	port     string
	producer sarama.SyncProducer
	nc       *nats.Conn
}

// service port
const PORT = ":8765"

// DB connection params
const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "users"
)

// kafka connection params
const (
	broker = "localhost:9092"
	topic  = "newuser"
)

// NATS connection
const (
	subject = "get.balance"
)

// Init configures and initializes the service and DB connection
func Init() *Service {
	L.Logger.Println("Service is being initialized")

	var srv Service

	// initialize REST service
	srv.muxSetup()

	// initialiye DB connection
	srv.dbSetup()

	// initialize kafka connection
	srv.kafkaSetup()

	// initialize NATS connection
	srv.natsSetup()

	L.Logger.Println("Service initialized")
	return &srv

}

func (srv *Service) muxSetup() {
	srv.port = PORT

	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/createUser", http.HandlerFunc(srv.createUserHandler))
	mux.Handle("/balance", http.HandlerFunc(srv.getUserBalance))

	srv.mux = mux
}

func (srv *Service) dbSetup() {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		L.Logger.Fatalf("Failed initializin DB connection. Connection string: %s. Error: %s", connectionString, err.Error())
	}
	srv.repo.db = db

	err = srv.repo.db.Ping()
	if err != nil {
		L.Logger.Fatalf("Failed to ping the DB. Connection string: %s. Error: %s", connectionString, err.Error())
	} else {
		L.Logger.Println("DB connection established")
	}
}

func (srv *Service) kafkaSetup() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(strings.Split(broker, ","), config)
	if err != nil {
		L.Logger.Fatal("Error creating Kafka producer: ", err)
	}
	srv.producer = producer
}

func (srv *Service) natsSetup() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		L.Logger.Fatal("Failet connecting to NATS: ", err)
	}

	srv.nc = nc
}

// Run is used to start service
func (s *Service) Run() {
	L.Logger.Println("Running the service")
	go func() {
		err := http.ListenAndServe(s.port, s.mux)
		if err != nil {
			L.Logger.Fatalf("failed startign the service on port %s :%s", s.port, err.Error())
		} else {
			L.Logger.Printf("The service is up and running on port %s", s.port)
		}
	}()
}

// Stop is stopping the service and closes DB connection
func (s *Service) Stop() {
	L.Logger.Print("Closing DB connection...")
	s.repo.db.Close()

	L.Logger.Print("Closing Kafka producer...")
	if err := s.producer.Close(); err != nil {
		L.Logger.Print("Error closing Kafka producer: ", err)
	}

	L.Logger.Print("Closing NATS connection...")
	s.nc.Close()
}

func (s *Service) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		L.Logger.Printf("Request method unallowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Request method unallowed: %s", r.Method)))
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var userData createUserRequest
	err := d.Decode(&userData)
	if err != nil {
		L.Logger.Printf("Failed decoding the request. %s", err.Error())
		L.Logger.Printf("Payload: %+v", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed decoding the request. %s", err.Error())))
		return
	}

	if userData.Email == "" {
		L.Logger.Print("Email is mandatory in request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email is mandatory in request"))
		return
	}

	err = s.repo.createUser(userData.Email)
	if err != nil {
		L.Logger.Printf("Failed creating new user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed creating new user. %s", err.Error())))
		return
	}

	// KAFKA notify
	newUserNotification := newUserKafkaMessage{
		Email:     userData.Email,
		CreatedAt: time.Now(),
	}

	msgJSON, err := json.Marshal(newUserNotification)
	if err != nil {
		L.Logger.Print("Error encoding newUserKafkaMessage struct to JSON: ", err)
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msgJSON),
	}

	_, _, err = s.producer.SendMessage(&msg)
	if err != nil {
		L.Logger.Printf("Failed sending kafka message %s : ", msg.Value, err)
	}

	w.Write([]byte("New user created"))
}

func (s *Service) getUserBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		L.Logger.Printf("Request method unallowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Request method unallowed: %s", r.Method)))
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var userData getBalanceRequest
	err := d.Decode(&userData)
	if err != nil {
		L.Logger.Printf("Failed decoding the request. %s", err.Error())
		L.Logger.Printf("Payload: %+v", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed decoding the request. %s", err.Error())))
		return
	}

	id, err := s.repo.getUserIDFromEmail(userData.Email)

	if err != nil {
		L.Logger.Print("Failed getting userID based on the given e-mail: ", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Failed getting userID based on the given e-mail: %s", err.Error())))
		return
	}

	/*
		NATS request to transactions service
		transctions service processing...
		NATS response from transactions service
	*/

	natsResp, err := s.nc.Request(subject, []byte(fmt.Sprintf("%d", id)), 10*time.Millisecond)
	if err != nil {
		L.Logger.Printf("Failed getting balance for user %s: %s", userData.Email, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed getting balance for user %s: %s", userData.Email, err.Error())))
		return
	}

	if natsResp == nil {
		L.Logger.Printf("NATS request timed out")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("NATS request timed out"))
		return
	}

	var balance float64

	err = json.Unmarshal(natsResp.Data, &balance)
	if err != nil {
		L.Logger.Printf("Error decoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error decoding NATS response: %v", err)))
		return
	}

	L.Logger.Print("Balance retreived ", balance)

	resp := getBalanceResponse{
		Email:   userData.Email,
		Balance: balance,
	}

	msgJSON, err := json.Marshal(resp)
	if err != nil {
		L.Logger.Print("Error encoding getBalanceResponse struct to JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error encoding getBalanceResponse struct to JSON: ", err)))
		return
	}

	w.Write([]byte(msgJSON))
}
