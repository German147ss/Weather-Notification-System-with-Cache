package handlers

import (
	"app/internal/core/services"
	"context"
	"encoding/json"
	"net/http"
)

func NewWeatherHandler(weatherService *services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

type WeatherHandler struct {
	weatherService *services.WeatherService
}

// Get weather handler
func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	city := r.URL.Path[len("/weather/"):]
	weather, err := h.weatherService.GetWeather(city, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(weather)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

// Search id by name handler
func (h *WeatherHandler) SearchIdByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/weather/city/"):]
	id, err := h.weatherService.SearchIdByName(name, context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(id))
}

func (h *WeatherHandler) GetWeatherAndWaves(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	city := r.URL.Path[len("/weather/waves/"):]
	weatherAndWaves, err := h.weatherService.GetWeatherAndWaves(city, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(weatherAndWaves)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
