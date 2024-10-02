package main

import (
	"app/internal/core/adapters"
	"app/internal/core/ports"
	"app/internal/core/services"
	"app/internal/handlers"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func main() {
	cfg := NewConfig()

	var cacheService ports.CacheService

	if cfg.CacheType == "redis" {
		rdb := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			DB:   0,
		})

		cacheService = adapters.NewRedisCacheService(rdb)
	} else {
		cacheService = adapters.NewMemoryCacheService()
	}

	weatherService := adapters.NewCPTECWeatherService()

	service := services.New(weatherService, cacheService)

	// Initialize HTTP handler
	handler := handlers.NewWeatherHandler(service)

	// Set up HTTP routes
	http.HandleFunc("/weather/", handler.GetWeather)
	http.HandleFunc("/weather/city/", handler.SearchIdByName)
	http.HandleFunc("/weather/waves/", handler.GetWeatherAndWaves)

	// Start HTTP server
	fmt.Printf("Server started on port %s\n", cfg.AppPort)
	http.ListenAndServe(":"+cfg.AppPort, nil)
}
