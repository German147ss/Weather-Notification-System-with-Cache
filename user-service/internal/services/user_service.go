package services

import (
	"user-service/internal/entities"
	"user-service/internal/ports"
)

type UserService struct {
	UserRepository       ports.UserRepository
	WeatherService       ports.WeatherService
	NotificationConsumer ports.NotificationConsumer
}

func NewUserService(repo ports.UserRepository, weatherService ports.WeatherService, consumer ports.NotificationConsumer) *UserService {
	return &UserService{
		UserRepository:       repo,
		WeatherService:       weatherService,
		NotificationConsumer: consumer,
	}
}

func (s *UserService) RegisterUser(user entities.User) (*entities.WeatherAndWaves, error) {
	err := s.UserRepository.InsertUser(user)
	if err != nil {
		return nil, err
	}
	return s.WeatherService.GetWeatherAndWaves(user.LocationCode)
}

func (s *UserService) DeactivateNotifications(userID string) error {
	return s.UserRepository.DeactivateUserNotifications(userID)
}
