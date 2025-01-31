package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gookit/event"
	"github.com/segmentio/kafka-go"
)

// Register event for each mutating operation on Companies
func RegisterEventListeners(w *kafka.Writer) {
	log.Println("Registering event listeners...")

	event.On("company_created", event.ListenerFunc(func(e event.Event) error {
		fmt.Printf("handling event: %s\n", e.Name())
		message, err := prepareKafkaMessage("company_created", e.Data())
		if err == nil {
			err := w.WriteMessages(context.Background(), *message)
			if err != nil {
				log.Fatalf("could not write message: %v", err)
			}
		}
		return nil
	}), event.Normal)

	event.On("company_updated", event.ListenerFunc(func(e event.Event) error {
		fmt.Printf("handling event: %s\n", e.Name())
		message, err := prepareKafkaMessage("company_updated", e.Data())
		if err == nil {
			err := w.WriteMessages(context.Background(), *message)
			if err != nil {
				log.Fatalf("could not write message: %v", err)
			}
		}
		return nil
	}), event.Normal)

	event.On("company_deleted", event.ListenerFunc(func(e event.Event) error {
		fmt.Printf("handling event: %s\n", e.Name())
		message, err := prepareKafkaMessage("company_deleted", e.Data())
		if err == nil {
			err := w.WriteMessages(context.Background(), *message)
			if err != nil {
				log.Fatalf("could not write message: %v", err)
			}
		}
		return nil
	}), event.Normal)
}

func prepareKafkaMessage(key string, value map[string]any) (*kafka.Message, error) {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &kafka.Message{
		Key:   []byte(key),
		Value: valueBytes,
	}, nil
}
