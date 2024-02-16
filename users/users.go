package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	L "userTransactions/logging"
)

type Service struct {
	mux  *http.ServeMux
	repo Repo
	port string
}

const PORT = ":8765"

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "users"
)

// Init configures and initializes the service and DB connection
func Init() *Service {
	L.Logger.Println("Service is being initialized")

	var srv Service

	// initialize REST service
	srv.port = PORT

	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/createUser", http.HandlerFunc(srv.createUserHandler))
	mux.Handle("/balance", http.HandlerFunc(getUserBalance))

	srv.mux = mux

	// initialiye DB connection
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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

	L.Logger.Println("Service initialized")
	return &srv

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
	L.Logger.Print("Closing DB connection")
	s.repo.db.Close()
	L.Logger.Print("DB connection closed")
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
}

func getUserBalance(w http.ResponseWriter, r *http.Request) {
	panic("getUserBalance unimplemented")
}
