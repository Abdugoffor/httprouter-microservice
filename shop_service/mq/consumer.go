package mq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	{
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}
	}
	ch, err := conn.Channel()
	{
		if err != nil {
			log.Fatalf("Failed to open a channel: %v", err)
		}
	}

	msgs, err := ch.Consume(
		"user_events",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	fmt.Println("🎧 RabbitMQ listening...")

	for msg := range msgs {
		fmt.Println("RabbitMQ message:", string(msg.Body))
	}

}
