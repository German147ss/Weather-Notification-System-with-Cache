package ports

import "user-service/internal/entities"

type WeatherService interface {
	GetWeatherAndWaves(city string) (*entities.WeatherAndWaves, error)
}
