package services

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"vibely-backend/src/models"
	"vibely-backend/src/repositories"
)

type UserService interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByID(userID uuid.UUID) (models.User, error)
	GetUserFromAccessToken(token string) (models.User, error)
	GetUserFromRefreshToken(refreshToken string) (models.User, error)
	UserFromDiscord(discordID, username, email string) (models.User, error)
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
	// Check if user with this email exists
	_, err := s.userRepository.GetUserByEmail(user.Email)
	if err == nil {
		return models.User{}, errors.New("user with that email already exists")
	}

	// Check if user with this username exists
	_, err = s.userRepository.GetUserByUsername(user.Username)
	if err == nil {
		return models.User{}, errors.New("user with that username already exists")
	}

	// Create the user in repository
	return s.userRepository.CreateUser(user)
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
func (s *userService) GetUserByEmail(email string) (models.User, error) {
	// You need to make sure your repository has this method implemented
	return s.userRepository.GetUserByEmail(email)
}
func (s *userService) UserFromDiscord(discordID, username, email string) (models.User, error) {
	// Try to find the user by email first
	user, err := s.GetUserByEmail(email)

	// If user doesn't exist, create a new one
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := models.User{
				Username:  username,
				Email:     email,
				DiscordID: discordID,
				Role:      models.UserRoleStudent,
			}

			return s.CreateUser(newUser)
		}

		return models.User{}, err
	}
	//TODO: UPDATE USER
	// Update existing user with Discord info if needed
	//if user.DiscordID != discordID {
	//	user.DiscordID = discordID
	//	return s.UpdateUser(user)
	//}

	return user, nil
}
