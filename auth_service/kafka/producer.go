package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func InitProducer() {
	writer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "auth_events",
		Balancer: &kafka.LeastBytes{},
	}
}

func Publish(data interface{}) {
	body, _ := json.Marshal(data)

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: body,
		},
	)

	if err != nil {
		log.Println("Kafka error:", err)
	}
}
