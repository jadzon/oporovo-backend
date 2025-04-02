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
	GetUserFromAccessToken(token string) (models.User, error)
	GetUserFromRefreshToken(refreshToken string) (models.User, error)
}
type userService struct {
	userRepository repositories.UserRepository
	authService    AuthService
}

func NewUserService(repo repositories.UserRepository, authService AuthService) UserService {
	return &userService{
		userRepository: repo,
		authService:    authService,
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
func (s *userService) GetUserFromAccessToken(token string) (models.User, error) {
	userID, err := s.authService.ExtractUserIDfromAccessToken(token)
	if err != nil {
		return models.User{}, err
	}
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (s *userService) GetUserFromRefreshToken(refreshToken string) (models.User, error) {
	userID, err := s.authService.ExtractUserIDfromRefreshToken(refreshToken)
	if err != nil {
		return models.User{}, err
	}
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
