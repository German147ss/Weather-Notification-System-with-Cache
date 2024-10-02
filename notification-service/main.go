package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

var DB *sql.DB

// Estructura para representar una notificación
type Notification struct {
	ID                   int
	LocationCode         string `json:"location_code"`
	NotificationSchedule int    `json:"notification_schedule"`
	State                string
}

// Función que obtiene notificaciones programadas para ser enviadas
func fetchScheduledNotifications(db *sql.DB) ([]Notification, error) {
	var notifications []Notification

	// Obtener la hora actual
	begin := getCurrentTimeInSeconds()

	// rango
	ends := begin + 59

	// Consulta SQL para obtener las notificaciones pendientes en el rango de tiempo
	query := `
        SELECT id, location_code, notification_schedule
        FROM user_preferences
        WHERE is_enabled = true AND notification_schedule >= $1 AND notification_schedule <= $2
    `

	// Ejecutar la consulta
	rows, err := db.Query(query, begin, ends)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No hay notificaciones programadas")
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener notificaciones programadas: %v", err)
	}
	defer rows.Close()

	// Iterar sobre los resultados
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.ID, &notification.LocationCode, &notification.NotificationSchedule)
		if err != nil {
			log.Printf("Error al escanear notificación: %v", err)
			continue
		}
		notifications = append(notifications, notification)
	}

	// Retornar la lista de notificaciones programadas
	return notifications, nil
}

// Función para obtener el horario actual del día en segundos
func getCurrentTimeInSeconds() int {
	now := time.Now()
	hours := now.Hour()
	minutes := now.Minute()
	seconds := now.Second()

	// Convertir el horario actual a segundos
	return hours*3600 + minutes*60 + seconds
}

// Función para verificar si se debe enviar la notificación
func shouldSendNotification(currentTimeInSeconds, notificationSchedule int) bool {
	return currentTimeInSeconds >= notificationSchedule && currentTimeInSeconds < (notificationSchedule+59)
}

func main() {

	rabbitConn, rabbitChannel, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	fmt.Println("Conectado a RabbitMQ y cola declarada")

	// Supongamos que db es la conexión a PostgreSQL
	db := initDB()

	// Crear el objeto cron
	c := cron.New()

	// Definir un job que publique notificaciones cada minuto
	c.AddFunc("@every 1m", func() {
		currentTimeInSeconds := getCurrentTimeInSeconds()
		fmt.Println("Horario actual:", currentTimeInSeconds)
		fmt.Println("Ejecutando el job para publicar notificaciones...")
		// Obtener notificaciones programadas para el momento actual
		notifications, err := fetchScheduledNotifications(db)
		if err != nil {
			log.Fatalf("Error obteniendo notificaciones: %v", err)
		}

		// Procesar las notificaciones
		for _, notification := range notifications {
			if shouldSendNotification(currentTimeInSeconds, notification.NotificationSchedule) {
				weatherResponse, err := GetWeather(notification.LocationCode)
				if err != nil {
					log.Fatalf("Error obteniendo clima para la ciudad: %+v", err)
				}
				log.Printf("Clima obtenido para la ciudad: %+v", weatherResponse)
				err = publishNotification(rabbitChannel, *weatherResponse)
				if err != nil {
					log.Fatalf("Error al publicar notificación: %v", err)
				}
			} else {
				fmt.Printf("Notificación %d no está dentro del horario actual\n", notification.ID)
			}
		}
	})

	// Iniciar el cron scheduler
	c.Start()

	// Mantener el programa en ejecución
	select {}
}

func initDB() *sql.DB {
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

	//make ping
	_, err = DB.Exec("SELECT 1")
	if err != nil {
		fmt.Println("Error pinging the database")
	}

	return DB
}

// func to create table if not exists
func createTableIfNotExists(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user_preferences (
		id SERIAL PRIMARY KEY,
		location_code VARCHAR(255) NOT NULL,
		notification_schedule INT NOT NULL DEFAULT 28800,
		is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
		state VARCHAR(255) NOT NULL DEFAULT 'pending'
	);`)
	if err != nil {
		fmt.Println("Error al crear tabla:", err)
		panic(err)
	}
}
