package adapters

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"user-service/internal/entities"
	"user-service/internal/ports"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQNotificationConsumer struct {
	Channel *amqp.Channel
}

func NewRabbitMQNotificationConsumer(channel *amqp.Channel) ports.NotificationConsumer {
	return &RabbitMQNotificationConsumer{
		Channel: channel,
	}
}

func (c *RabbitMQNotificationConsumer) ConsumeUserNotifications() error {
	msgs, err := c.Channel.Consume(
		"user_notifications", // queue name
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		return fmt.Errorf("error consuming notifications: %v", err)
	}

	go func() {
		for d := range msgs {
			var notification entities.WeatherResponse
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				log.Printf("error deserializing notification: %v", err)
				continue
			}
			log.Printf("received notification: %+v", notification)
		}
	}()
	return nil
}

func ConnectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	port := os.Getenv("RABBITMQ_PORT")
	if port == "" {
		port = "5672"
	}
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("error creating channel in RabbitMQ: %v", err)
	}

	_, err = ch.QueueDeclare(
		"user_notifications", // queue name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error declaring queue: %v", err)
	}
	return conn, ch, nil
}
