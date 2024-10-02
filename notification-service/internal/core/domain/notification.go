package domain

// Estructura para representar una notificación
type Notification struct {
	ID                   int
	LocationCode         string `json:"location_code"`
	NotificationSchedule int    `json:"notification_schedule"`
	State                string
}
