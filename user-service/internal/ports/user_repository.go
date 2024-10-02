package ports

import "user-service/internal/entities"

type UserRepository interface {
	InsertUser(user entities.User) error
	DeactivateUserNotifications(id string) error
}
