package ports

import (
	"app/internal/core/domain"
	"context"
	"time"
)

// WeatherService interface
type WeatherService interface {
	GetWeather(city string) (*domain.CityWeather, error)
	SearchIdByName(cityName string) (string, error)
	GetWaves(city string) (*domain.CityWaves, error)
}

// CacheService interface
type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}
