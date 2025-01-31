package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"xm-company/internal/database"
	"xm-company/internal/events"
	"xm-company/server"

	"github.com/segmentio/kafka-go"
	"github.com/subosito/gotenv"
)

const DEV_ENV string = "dev"

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalf("gotenv.Load() error: %s\n", err.Error())
	}
}

func main() {
	log.Println("Starting company service...")

	env := os.Getenv("ENV")
	if env == "" {
		env = DEV_ENV
	}
	isDevMode := env == DEV_ENV

	port := os.Getenv("PORT")

	connParams, err := database.ReadConnectionStringParams()
	if err != nil {
		log.Fatalf("DB initialization error: %s", err.Error())
	}

	connectionStr := database.CreateConnectionString(connParams.Host, connParams.Port, connParams.User, connParams.Password, connParams.DBName)
	db, err := database.Initialize(connectionStr)
	if err != nil {
		log.Fatalf("DB initialization error: %s", err.Error())
	}

	kafkaBrokersStr := os.Getenv("KAFKA_BROKER_URLS")
	if kafkaBrokersStr == "" {
		log.Fatal("Missing Kafka addresses")
	}
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")
	// Kafka writer setup
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kafkaBrokers,
		Topic:    "companies-events-topic",
		Balancer: &kafka.LeastBytes{},
	})

	server := &http.Server{
		Addr:              port,
		Handler:           server.SetupRouter(db, isDevMode),
		ReadHeaderTimeout: 1 * time.Second,
	}

	go func() {
		log.Printf("Server listening on port %s\n", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("server.ListenAndServe() error: %s\n", err.Error())
		}
	}()

	// Register event listeners
	go events.RegisterEventListeners(writer)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, os.Interrupt)

	<-stop

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Println("Shutting down company service")

	// Closing Kafka writer
	if err := writer.Close(); err != nil {
		log.Printf("could not close writer: %v", err.Error())
	}

	err = server.Shutdown(stopCtx)
	if err != nil {
		log.Printf("server.Shutdown() error: %s\n", err.Error())
	}
}
