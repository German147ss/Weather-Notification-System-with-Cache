package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type WeatherResponse struct {
	Nome        string     `json:"nome"`
	UF          string     `json:"uf"`
	Atualizacao string     `json:"atualizacao"`
	Previsoes   []Previsao `json:"previsao"`
}

type Previsao struct {
	Dia    string  `json:"dia"`
	Tempo  string  `json:"tempo"`
	Maxima int     `json:"maxima"`
	Minima int     `json:"minima"`
	IUV    float32 `json:"iuv"`
}

func GetWeather(city string) (*WeatherResponse, error) {
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

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
