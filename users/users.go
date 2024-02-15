package users

import (
	"encoding/json"
	"net/http"
	L "userTransactions/logging"
)

type Service struct {
	mux  *http.ServeMux
	repo Repo
	port string
}

const PORT = ":8765"

func Init() *Service {
	L.Logger.Println("Service is being initialized")

	var srv Service

	srv.port = PORT

	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/createUser", http.HandlerFunc(srv.createUserHandler))
	mux.Handle("/balance", http.HandlerFunc(getUserBalance))

	srv.mux = mux

	L.Logger.Println("Service initialized")
	return &srv

}

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

func (s *Service) Stop() {
	L.Logger.Print("Closing DB connection")
	s.repo.db.Close()
	L.Logger.Print("DB connection closed")
}

func (s Service) createUserHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var userData createUserRequest
	err := d.Decode(&userData)
	if err != nil {
		L.Logger.Printf("Failed decoding the request. %s", err.Error())
		L.Logger.Printf("Payload: %+v", r.Body)
		return
	}

	panic("createUserHandler unimplemented")
}

func getUserBalance(w http.ResponseWriter, r *http.Request) {
	panic("getUserBalance unimplemented")
}
