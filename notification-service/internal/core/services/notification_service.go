// services/notification_service.go
package services

import (
	"fmt"
	"notification-service/internal/core/ports"
	"time"
)

type NotificationService struct {
	NotificationRepo ports.NotificationRepository
	WeatherService   ports.WeatherService
	NotificationPubl ports.NotificationPublisher
}

func NewNotificationService(repo ports.NotificationRepository, weatherService ports.WeatherService, publ ports.NotificationPublisher) *NotificationService {
	return &NotificationService{
		NotificationRepo: repo,
		WeatherService:   weatherService,
		NotificationPubl: publ,
	}
}

func (s *NotificationService) SendScheduledNotifications() {
	currentTimeInSeconds := getCurrentTimeInSeconds()
	fmt.Println("Horario actual:", currentTimeInSeconds)

	begin := currentTimeInSeconds
	end := begin + 59

	notifications, err := s.NotificationRepo.FetchScheduledNotifications(begin, end)
	if err != nil {
		fmt.Printf("Error obteniendo notificaciones: %v\n", err)
		return
	}
	fmt.Printf("Notificaciones obtenidas: %+v\n", notifications)

	for _, notification := range notifications {
		weatherResponse, err := s.WeatherService.GetWeatherAndWaves(notification.LocationCode)
		if err != nil {
			fmt.Printf("Error obteniendo clima para la ciudad: %+v\n", err)
			continue
		}
		fmt.Printf("Clima obtenido para la ciudad: %+v\n", weatherResponse)
		err = s.NotificationPubl.Publish(*weatherResponse)
		if err != nil {
			fmt.Printf("Error al publicar notificación: %v\n", err)
		}

	}
}

// Función para obtener el horario actual del día en segundos
func getCurrentTimeInSeconds() int {
	now := time.Now()
	return now.Hour()*3600 + now.Minute()*60 + now.Second()
}
