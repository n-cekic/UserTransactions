package main

import (
	"flag"
	"os"
	"os/signal"
	L "userTransactions/logging"
	"userTransactions/users"
)

var (
	dbHost     = flag.String("db.host", "localhost", "Database Host")
	dbPort     = flag.String("db.port", "5432", "Database Port")
	dbUser     = flag.String("db.user", "postgres", "Database User")
	dbPassword = flag.String("db.password", "postgres", "Database Password")
	dbName     = flag.String("db.name", "users", "Database Name")
	port       = flag.String("service.port", ":8080", "Port for service to run on")
	broker     = flag.String("kafka.brokers", "localhost:9092", "Kafka broker address")
	topic      = flag.String("kafka.topic", "new.user", "Kafka topic to publish messages to")
	subject    = flag.String("nats.subject", "get.balance.", "NATS subject pattern to publish to")
)

func main() {
	L.Logger.Info("starting service user")

	// initialiye service
	srv := users.Init(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName, *subject, *broker, *topic)

	// start service
	srv.Run(*port)

	// shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	<-shutdownCh
	L.Logger.Info("Received interrupt signal. Shutting down gracefully...")

	srv.Stop()
}
