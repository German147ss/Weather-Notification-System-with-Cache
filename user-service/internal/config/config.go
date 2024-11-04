package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	var err error
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable password=%s port=%s", dbUser, dbName, dbHost, dbPassword, dbPort)
	fmt.Println("Connecting to database with:", connectionString)

	var db *sql.DB
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connectionString)
		if err == nil && db.Ping() == nil {
			fmt.Println("Err", err)
			fmt.Println("Successfully connected to the database.")
			break
		}
		fmt.Println("Error connecting to the database, retrying in 3 seconds...", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		fmt.Println("Failed to connect to the database after multiple attempts.")
		panic(err)
	}

	createTableIfNotExists(db)
	return db
}

func createTableIfNotExists(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user_preferences (
        id SERIAL PRIMARY KEY,
        location_code VARCHAR(255) NOT NULL,
        notification_schedule INT NOT NULL DEFAULT 28800,
        is_enabled BOOLEAN NOT NULL DEFAULT TRUE
    );`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		panic(err)
	}
}
