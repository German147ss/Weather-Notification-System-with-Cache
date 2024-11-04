package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"user-service/internal/adapters"
	"user-service/internal/config"
	"user-service/internal/entities"
	"user-service/internal/services"
)

func main() {
	// Initialize database and repositories
	dbConn := config.InitDB()
	defer dbConn.Close()
	userRepository := adapters.NewPostgresUserRepository(dbConn)

	// Initialize HTTP client for weather service
	weatherService := adapters.NewWeatherServiceClient()

	rabbitConn, rabbitChannel, err := adapters.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	defer rabbitChannel.Close()
	notificationConsumer := adapters.NewRabbitMQNotificationConsumer(rabbitChannel)
	go notificationConsumer.ConsumeUserNotifications()

	userService := services.NewUserService(userRepository, weatherService, notificationConsumer)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var user entities.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid user data", http.StatusBadRequest)
			return
		}
		weatherAndWaves, err := userService.RegisterUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weatherAndWaves)
	})

	http.HandleFunc("/opt-out/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Path[len("/opt-out/"):]
		if err := userService.DeactivateNotifications(userID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notificaciones desactivadas para el usuario"))
	})

	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		adapters.Mu.Lock()
		defer adapters.Mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(adapters.Notifications)
	})

	// Start HTTP server
	fmt.Println("Servidor iniciado en el puerto 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
