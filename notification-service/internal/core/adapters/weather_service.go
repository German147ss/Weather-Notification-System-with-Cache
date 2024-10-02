package adapters

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"notification-service/internal/core/domain"
	"notification-service/internal/core/ports"
	"os"
)

type WeatherAPIService struct{}

func NewWeatherAPIService() ports.WeatherService {
	return &WeatherAPIService{}
}

func (s *WeatherAPIService) GetWeather(city string) (*domain.CityWeather, error) {
	baseURL := os.Getenv("WEATHER_API_BASE_URL")
	if baseURL == "" {
		return nil, errors.New("WEATHER_API_BASE_URL environment variable not set")
	}

	url := baseURL + "/weather/" + city
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error fetching weather data")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather domain.CityWeather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

func (s *WeatherAPIService) GetWeatherAndWaves(city string) (*domain.WeatherAndWaves, error) {
	baseURL := os.Getenv("WEATHER_API_BASE_URL")
	if baseURL == "" {
		return nil, errors.New("WEATHER_API_BASE_URL environment variable not set")
	}

	url := baseURL + "/weather/waves/" + city
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error fetching weather data")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather domain.WeatherAndWaves
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
