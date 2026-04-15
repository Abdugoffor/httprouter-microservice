package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func Consume() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "auth_events",
		GroupID: "shop-service",
	})

	fmt.Println("🎧 Kafka listening...")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println("Kafka message:", string(msg.Value))
	}
}
