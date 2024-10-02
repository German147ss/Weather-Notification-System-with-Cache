// ports/notification_repo.go
package ports

import "notification-service/internal/core/domain"

type NotificationRepository interface {
	FetchScheduledNotifications(begin, end int) ([]domain.Notification, error)
}

type WeatherService interface {
	GetWeather(city string) (*domain.CityWeather, error)
	GetWeatherAndWaves(city string) (*domain.WeatherAndWaves, error)
}

type NotificationPublisher interface {
	Publish(notification domain.CityWeather) error
}
