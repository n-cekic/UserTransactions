package users

import (
	"database/sql"
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
	subject  string
	topic    string
	producer sarama.SyncProducer
	nc       *nats.Conn
}

// Init configures and initializes the service
func Init(dbHost, dbPort, dbUser, dbPassword, dbName, subject, broker, topic string) *Service {
	L.Logger.Info("Service is being initialized")

	var srv Service
	srv.subject = subject
	srv.topic = topic

	// initialize REST service
	srv.muxSetup()

	// initialiye DB connection
	srv.dbSetup(dbHost, dbPort, dbUser, dbPassword, dbName)

	// initialize kafka connection
	srv.kafkaSetup(broker)

	// initialize NATS connection
	srv.natsSetup()

	L.Logger.Info("Service initialized")
	return &srv
}

func (srv *Service) muxSetup() {
	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/createuser", http.HandlerFunc(srv.createUserHandler))
	mux.Handle("/balance", http.HandlerFunc(srv.getUserBalance))

	srv.mux = mux
}

func (srv *Service) dbSetup(dbHost, dbPort, dbUser, dbPassword, dbName string) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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
		L.Logger.Info("DB connection established")
	}
}

func (srv *Service) kafkaSetup(broker string) {
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
		L.Logger.Fatal("Failed connecting to NATS: ", err)
	}

	srv.nc = nc
}

// Run is used to start service
func (s *Service) Run(port string) {
	L.Logger.Info("Running the service")
	go func() {
		err := http.ListenAndServe(port, s.mux)
		if err != nil {
			L.Logger.Fatalf("failed startign the service on port %s :%s", port, err.Error())
		} else {
			L.Logger.Infof("The service is up and running on port %s", port)
		}
	}()
}

// Stop is stopping the service and closes all connections
func (s *Service) Stop() {
	L.Logger.Info("Closing DB connection...")
	s.repo.db.Close()

	L.Logger.Info("Closing Kafka producer...")
	if err := s.producer.Close(); err != nil {
		L.Logger.Info("Error closing Kafka producer: ", err)
	}

	L.Logger.Info("Closing NATS connection...")
	s.nc.Close()
}
