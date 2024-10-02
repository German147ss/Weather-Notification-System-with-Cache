package adapters

import (
	"database/sql"
	"fmt"
	"notification-service/internal/core/domain"
	"notification-service/internal/core/ports"
)

type PostgresNotificationRepository struct {
	DB *sql.DB
}

func NewPostgresNotificationRepository(db *sql.DB) ports.NotificationRepository {
	return &PostgresNotificationRepository{DB: db}
}

func (r *PostgresNotificationRepository) FetchScheduledNotifications(begin, end int) ([]domain.Notification, error) {
	var notifications []domain.Notification

	query := `
        SELECT id, location_code, notification_schedule
        FROM user_preferences
        WHERE is_enabled = true AND notification_schedule >= $1 AND notification_schedule <= $2
    `

	rows, err := r.DB.Query(query, begin, end)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No hay notificaciones programadas")
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener notificaciones programadas: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notification domain.Notification
		err := rows.Scan(&notification.ID, &notification.LocationCode, &notification.NotificationSchedule)
		if err != nil {
			fmt.Printf("Error al escanear notificación: %v\n", err)
			continue
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
