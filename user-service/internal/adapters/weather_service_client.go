package adapters

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"user-service/internal/entities"
	"user-service/internal/ports"
)

type WeatherServiceClient struct {
	client *http.Client
}

func NewWeatherServiceClient() ports.WeatherService {
	return &WeatherServiceClient{
		client: &http.Client{},
	}
}

func (s *WeatherServiceClient) GetWeatherAndWaves(city string) (*entities.WeatherAndWaves, error) {
	baseURL := os.Getenv("WEATHER_API_BASE_URL")
	if baseURL == "" {
		return nil, errors.New("WEATHER_API_BASE_URL environment variable not set")
	}

	url := baseURL + "/weather/waves/" + city
	resp, err := s.client.Get(url)
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

	var weather entities.WeatherAndWaves
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
