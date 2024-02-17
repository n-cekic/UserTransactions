package main

import (
	"flag"
	"os"
	"os/signal"
	L "userTransactions/logging"
	"userTransactions/transactions"
)

var (
	servicePort = flag.String("service.port", ":8081", "Port for service to run on")
	dbHost      = flag.String("db.host", "localhost", "Database Host")
	dbPort      = flag.String("db.port", "5432", "Database Port")
	dbUser      = flag.String("db.user", "postgres", "Database User")
	dbPassword  = flag.String("db.password", "postgres", "Database Password")
	dbName      = flag.String("db.name", "transactions", "Database Name")
	subject     = flag.String("nats.subject", "get.balance.*", "Subject to subscribe NATS to")
	brokers     = flag.String("kafka.brokers", "localhost:9092", "Kafka broker address")
	group       = flag.String("kafka.group", "new.users", "Kafka consumer group name")
	topic       = flag.String("kafka.topic", "new.user", "Kafka topic name")
)

func main() {
	L.Logger.Info("starting service transactions")

	flag.Parse()

	// initialiye service
	srv := transactions.Init(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName, *subject, *brokers, *group, *topic)

	// start service
	srv.Run(*servicePort)

	// shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	<-shutdownCh
	L.Logger.Info("Received interrupt signal. Shutting down gracefully...")

	srv.Stop()
}
