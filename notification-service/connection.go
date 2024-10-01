package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func connectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
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
		return nil, nil, fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("error creando canal en RabbitMQ: %v", err)
	}

	// Declarar la cola de notificaciones de usuarios
	_, err = ch.QueueDeclare(
		"user_notifications", // Nombre de la cola
		false,                // Durable
		false,                // Delete when unused
		false,                // Exclusive
		false,                // No wait
		nil,                  // Arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error declarando cola: %v", err)
	}

	return conn, ch, nil
}

// Función para publicar notificaciones en RabbitMQ
func publishNotification(ch *amqp.Channel, notification WeatherResponse) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error al convertir notificación a JSON: %v", err)
	}

	log.Printf("Notificación a enviar: %+v", notification)

	err = ch.Publish(
		"",                   // Exchange
		"user_notifications", // Routing key (nombre de la cola)
		false,                // Mandatory
		false,                // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("error al publicar notificación: %v", err)
	}

	log.Printf("Notificación publicada en RabbitMQ: %+v", notification)
	return nil
}
