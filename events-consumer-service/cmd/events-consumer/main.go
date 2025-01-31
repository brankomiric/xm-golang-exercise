package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/segmentio/kafka-go"
	"github.com/subosito/gotenv"
)

const COMPANIES_TOPIC = "companies-events-topic"

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalf("gotenv.Load() error: %s\n", err.Error())
	}
}

func main() {
	kafkaBrokersStr := os.Getenv("KAFKA_BROKER_URLS")
	if kafkaBrokersStr == "" {
		log.Fatal("Missing Kafka addresses")
	}
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")

	// Set up the Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   COMPANIES_TOPIC,
		GroupID: "default-consumer-group",
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, os.Interrupt)

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a goroutine to handle shutdown signals
	go func() {
		sig := <-stop
		log.Printf("Received signal: %v, shutting down...", sig)
		cancel()
	}()

	for {
		// Check if the context has been canceled (shutdown signal received)
		select {
		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			if err := reader.Close(); err != nil {
				log.Fatalf("could not close reader: %v", err)
			}
			return
		default:
			// Read a message from Kafka
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Only logging messages for now
			fmt.Printf("Received message with key: %s, value: %s\n", string(message.Key), string(message.Value))
		}
	}
}
