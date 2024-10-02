package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Estructura del usuario
type Usuario struct {
	ID                   int    `json:"id"`
	LocationCode         string `json:"location_code"`
	NotificationSchedule int    `json:"notification_schedule"`
	IsEnabled            bool   `json:"is_enabled"`
	State                string `json:"state"`
}

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

// Variable de la conexión a la base de datos
var DB *sql.DB

// Función para registrar un usuario en la base de datos
func registrarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario

	// Leer el cuerpo de la solicitud y decodificarlo
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		fmt.Println("Error al decodificar el cuerpo de la solicitud:", err)
		http.Error(w, "Datos de usuario inválidos", http.StatusBadRequest)
		return
	}

	// Preparar la consulta SQL para insertar el nuevo usuario
	sqlInsert := `INSERT INTO user_preferences (location_code, notification_schedule) VALUES ($1, $2) RETURNING id`

	// Ejecutar la consulta e insertar el usuario
	err = DB.QueryRow(sqlInsert, usuario.LocationCode, usuario.NotificationSchedule).Scan(&usuario.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al insertar usuario: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar el usuario registrado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// Funcion para dar de baja las notificaciones de un usuario
func desactivarNotificaciones(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario desde la URL
	id := r.URL.Path[len("/desactivar/"):]
	fmt.Println("ID del usuario:", id)

	// Desactivar las notificaciones para el usuario
	desactivarNotificacionesForUser(id)

	// Retornar un mensaje de respuesta
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notificaciones desactivadas para el usuario"))
}

func desactivarNotificacionesForUser(id string) {
	// Preparar la consulta SQL para desactivar las notificaciones para el usuario
	sqlUpdate := `UPDATE user_preferences SET is_enabled = false WHERE id = $1`

	// Ejecutar la consulta para desactivar las notificaciones
	_, err := DB.Exec(sqlUpdate, id)
	if err != nil {
		fmt.Println("Error al desactivar las notificaciones:", err)
		return
	}
}

func main() {
	// Conectar a la base de datos
	DB = initDB()

	defer DB.Close()

	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	go consumeUserNotifications(ch)

	// Ruta para registrar un nuevo usuario
	http.HandleFunc("/usuarios", registrarUsuario)

	// Ruta para desactivar las notificaciones de un usuario
	http.HandleFunc("/desactivar/", desactivarNotificaciones)

	// Iniciar el servidor HTTP
	fmt.Println("Servidor iniciado en el puerto 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
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

// Función para consumir notificaciones de la cola de RabbitMQ
func consumeUserNotifications(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		"user_notifications", // Nombre de la cola
		"",                   // Consumer
		true,                 // Auto-Acknowledge
		false,                // Exclusivo
		false,                // No-local
		false,                // No-wait
		nil,                  // Args
	)
	if err != nil {
		log.Fatalf("Error al consumir notificaciones: %v", err)
	}

	// Procesar los mensajes en segundo plano
	for d := range msgs { // Este es el bucle que sigue escuchando mensajes
		var notification WeatherResponse //TODO CHANGE THIS STRUCT IN NOTIFICATION SERVICE
		err := json.Unmarshal(d.Body, &notification)
		if err != nil {
			log.Printf("Error al deserializar notificación: %v", err)
			continue
		}

		// Aquí procesas la notificación (ej. enviar al usuario)
		log.Printf("Notificación recibida por el User Service: %+v", notification)
	}
}

func connectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {

	port := os.Getenv("RABBITMQ_PORT")
	if port == "" {
		port = "5672"
	}
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")
	if err != nil {
		return nil, nil, fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("error creando canal en RabbitMQ: %v", err)
	}

	// Declarar una cola
	_, err = ch.QueueDeclare(
		"notificaciones", // Nombre de la cola
		false,            // Durable
		false,            // Auto delete cuando no se esté usando
		false,            // Exclusiva
		false,            // No esperar
		nil,              // Argumentos adicionales
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error declarando cola: %v", err)
	}

	fmt.Println("Conectado a RabbitMQ y cola declarada")
	return conn, ch, nil
}
