package users

import (
	"net/http"
	L "userTransactions/logging"
)

type Service struct {
	mux  *http.ServeMux
	port string
}

const PORT = ":8765"

func Init() *Service {
	L.Logger.Println("Service is being initialized")

	var srv Service

	srv.port = PORT

	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/createUser", http.HandlerFunc(createUserHandler))
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

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	panic("createUserHandler unimplemented")
}

func getUserBalance(w http.ResponseWriter, r *http.Request) {
	panic("getUserBalance unimplemented")
}
