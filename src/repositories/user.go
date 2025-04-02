package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"vibely-backend/src/models"
)

type UserRepository interface {
	GetUserByUsername(username string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByID(userID uuid.UUID) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return models.User{}, nil
	}
	user, _ = r.GetUserByUsername(user.Username)
	return user, nil
}
func (r *userRepository) GetUserByID(userID uuid.UUID) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (r *userRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err // Return the error directly
	}
	return user, nil
}
func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, err // Return the error directly
	}
	return user, nil
}
