package main

import (
	"os"
	"os/signal"
	L "userTransactions/logging"
	"userTransactions/users"
)

func main() {
	L.Logger.Println("starting service user")

	srv := users.Init()
	srv.Run()

	// shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	<-shutdownCh

	L.Logger.Println("\nReceived interrupt signal. Shutting down gracefully...")
}
