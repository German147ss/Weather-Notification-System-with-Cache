// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"notification-service/internal/core/adapters"
	"notification-service/internal/core/services"

	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"github.com/streadway/amqp"
)

var (
	DB           *sql.DB
	RabbitMQConn *amqp.Connection
	RabbitMQChan *amqp.Channel
)

func main() {
	initDB()
	initRabbitMQ()

	repo := adapters.NewPostgresNotificationRepository(DB)
	weatherService := adapters.NewWeatherAPIService()
	publisher := adapters.NewRabbitMQNotificationPublisher(RabbitMQChan)

	service := services.NewNotificationService(repo, weatherService, publisher)

	// Crear el objeto cron
	c := cron.New()

	// Definir un job que publique notificaciones cada minuto
	c.AddFunc("@every 1m", service.SendScheduledNotifications)

	// Iniciar el cron scheduler
	c.Start()

	// Mantener el programa en ejecución
	select {}
}

func initDB() {
	var err error
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable password=%s port=%s", dbUser, dbName, dbHost, dbPassword, dbPort)
	fmt.Println("Connecting to database with:", connectionString)

	for i := 0; i < 5; i++ { // Retry logic
		DB, err = sql.Open("postgres", connectionString)
		if err == nil && DB.Ping() == nil {
			fmt.Println("Successfully connected to the database.")
			break
		}
		fmt.Println("Error connecting to the database, retrying in 5 seconds...")
		time.Sleep(5 * time.Second) // Wait before retrying
	}

	if err != nil {
		fmt.Println("Failed to connect to the database after multiple attempts.")
		panic(err)
	}

	createTableIfNotExists(DB)
}

func createTableIfNotExists(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user_preferences (
		id SERIAL PRIMARY KEY,
		location_code VARCHAR(255) NOT NULL,
		notification_schedule INT NOT NULL DEFAULT 28800,
		is_enabled BOOLEAN NOT NULL DEFAULT TRUE
	);`)
	if err != nil {
		fmt.Println("Error al crear tabla:", err)
		panic(err)
	}
}

func initRabbitMQ() {
	var err error
	port := os.Getenv("RABBITMQ_PORT")
	if port == "" {
		port = "5672"
	}
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}
	RabbitMQConn, err = amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}

	RabbitMQChan, err = RabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Error creando canal en RabbitMQ: %v", err)
	}

	_, err = RabbitMQChan.QueueDeclare(
		"user_notifications", // Nombre de la cola
		false,                // Durable
		false,                // Delete when unused
		false,                // Exclusive
		false,                // No wait
		nil,                  // Arguments
	)
	if err != nil {
		log.Fatalf("Error declarando cola: %v", err)
	}
}
