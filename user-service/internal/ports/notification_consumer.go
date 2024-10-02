package ports

type NotificationConsumer interface {
	ConsumeUserNotifications() error
}
