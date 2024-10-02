package adapters

import (
	"database/sql"
	"user-service/internal/entities"
	"user-service/internal/ports"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) ports.UserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) InsertUser(user entities.User) error {
	sqlInsert := `INSERT INTO user_preferences (location_code, notification_schedule) VALUES ($1, $2) RETURNING id`
	err := r.DB.QueryRow(sqlInsert, user.LocationCode, user.NotificationSchedule).Scan(&user.ID)
	return err
}

func (r *PostgresUserRepository) DeactivateUserNotifications(id string) error {
	sqlUpdate := `UPDATE user_preferences SET is_enabled = false WHERE id = $1`
	_, err := r.DB.Exec(sqlUpdate, id)
	return err
}
