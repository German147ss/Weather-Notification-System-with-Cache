package services

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type WeatherService struct {
	weatherService ports.WeatherService
	cacheService   ports.CacheService
}

func New(weatherService ports.WeatherService, cacheService ports.CacheService) *WeatherService {
	return &WeatherService{
		weatherService: weatherService,
		cacheService:   cacheService,
	}
}

func (s *WeatherService) GetWeather(city string, ctx context.Context) (*domain.CityWeather, error) {

	cacheKey := "weather:" + city

	val, err := s.cacheService.Get(ctx, cacheKey)
	if err == redis.Nil {
		fmt.Println("Cache miss. Getting data from CPTEC...")
		weather, err := s.weatherService.GetWeather(city)
		if err != nil {
			return nil, err
		}

		jsonData, err := json.Marshal(weather)
		if err != nil {
			return nil, err
		}

		err = s.cacheService.Set(ctx, cacheKey, string(jsonData), time.Hour)
		if err != nil {
			return nil, err
		}

		fmt.Println("Data cached for city:", city)
		return weather, nil
	} else if err != nil {
		return nil, err
	}

	fmt.Printf("cache hit for city: %s", city)
	var weather domain.CityWeather
	err = json.Unmarshal([]byte(val), &weather)
	if err != nil {
		return nil, err
	}
	return &weather, nil
}

// SearchIdByName
func (s *WeatherService) SearchIdByName(name string, ctx context.Context) (string, error) {
	// Generate the key for Redis
	cacheKey := "city-id:" + name

	// Try to get the city ID from the Redis cache
	val, err := s.cacheService.Get(ctx, cacheKey)
	if err == redis.Nil {
		fmt.Println("Cache miss. Getting city ID from CPTEC...")
		cityId, err := s.weatherService.SearchIdByName(name)
		if err != nil {
			return "", err
		}

		err = s.cacheService.Set(ctx, cacheKey, cityId, 24*time.Hour)
		if err != nil {
			return "", err
		}

		fmt.Println("City ID cached for name:", name)
		return cityId, nil
	} else if err != nil {
		return "", err
	}

	// City ID found in cache, return it
	fmt.Printf("Cache hit for city name: %s", name)
	return val, nil
}

// func to get WeatherAndWaves
func (s *WeatherService) GetWeatherAndWaves(city string, ctx context.Context) (*domain.WeatherAndWaves, error) {
	weather, err := s.GetWeather(city, ctx)
	if err != nil {
		fmt.Println("Error getting weather:", err)
		return nil, err
	}
	waves, err := s.weatherService.GetWaves(city)
	if err != nil {
		fmt.Println("Error getting waves:", err)
		return nil, err
	}
	return &domain.WeatherAndWaves{
		Weather: weather,
		Waves:   waves,
	}, nil
}
