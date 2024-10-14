package adapters

import (
	"encoding/json"
	"fmt"
	"notification-service/internal/core/domain"
	"notification-service/internal/core/ports"

	"github.com/streadway/amqp"
)

type RabbitMQNotificationPublisher struct {
	Channel *amqp.Channel
}

func NewRabbitMQNotificationPublisher(channel *amqp.Channel) ports.NotificationPublisher {
	return &RabbitMQNotificationPublisher{Channel: channel}
}

func (p *RabbitMQNotificationPublisher) Publish(notification domain.WeatherAndWaves) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error al convertir notificación a JSON: %v", err)
	}

	err = p.Channel.Publish(
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

	return nil
}
