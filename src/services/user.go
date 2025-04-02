package services

import (
	"errors"
	"github.com/google/uuid"
	"vibely-backend/src/models"
	"vibely-backend/src/repositories"
)

type UserService interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByID(userID uuid.UUID) (models.User, error)
}
type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		userRepository: repo,
	}
}
func (s *userService) CreateUser(user models.User) (models.User, error) {
	_, err := s.userRepository.GetUserByID(user.ID)
	if err == nil {
		return models.User{}, errors.New("user with that email already exists")
	}
	_, err = s.userRepository.GetUserByUsername(user.Username)
	if err == nil {
		return models.User{}, errors.New("user with that username already exists")
	}
	return user, nil
}
func (s *userService) GetUserByID(userID uuid.UUID) (models.User, error) {
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return models.User{}, err
	}
	return user, err
}
