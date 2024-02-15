package main

import (
	"os"
	"os/signal"
	L "userTransactions/logging"
	"userTransactions/users"
)

func main() {
	L.Logger.Println("starting service user")

	// initialiye service
	srv := users.Init()

	// start service
	srv.Run()

	// shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	<-shutdownCh
	L.Logger.Println("Received interrupt signal. Shutting down gracefully...")

	srv.Stop()
}
