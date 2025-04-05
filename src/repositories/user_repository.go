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
	GetTutors(filters models.TutorFilters) ([]models.User, int64, error)
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

// GetTutors retrieves tutor users based on the provided filters.
func (r *userRepository) GetTutors(filters models.TutorFilters) ([]models.User, int64, error) {
	var tutors []models.User
	var total int64

	// Base query: only consider users with the role "tutor".
	query := r.db.Model(&models.User{}).Where("role = ?", models.UserRoleTutor)

	// Apply subject filter if provided.
	if filters.Subject != "" {
		// Assuming Subjects is a PostgreSQL array column.
		query = query.Where("? = ANY(subjects)", filters.Subject)
	}

	// Apply level filter if provided.
	if filters.Level != "" {
		// Assuming Levels is a PostgreSQL array column.
		query = query.Where("? = ANY(levels)", filters.Level)
	}

	// Count total records matching the filters.
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset for pagination.
	offset := (filters.Page - 1) * filters.Limit
	if err := query.Offset(offset).Limit(filters.Limit).Find(&tutors).Error; err != nil {
		return nil, 0, err
	}

	return tutors, total, nil
}
