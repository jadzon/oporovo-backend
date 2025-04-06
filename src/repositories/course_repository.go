package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"vibely-backend/src/models"
)

// CourseRepository defines the methods to interact with the Course datastore.
type CourseRepository interface {
	CreateCourse(course *models.Course) error
	GetCourseByID(courseID uuid.UUID) (models.Course, error)
	GetCourseWithParticipants(courseID uuid.UUID) (models.Course, error)
	UpdateCourse(course *models.Course) error
	GetCourses(subject, level string, page, limit int) ([]models.Course, error)
	GetCoursesForUser(userID uuid.UUID) ([]models.Course, error)
}

type courseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new instance of CourseRepository.
func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

// CreateCourse inserts a new Course record along with its associations.
func (r *courseRepository) CreateCourse(course *models.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.WithContext(ctx).Create(course).Error
}

// GetCourseByID retrieves a Course by its ID.
func (r *courseRepository) GetCourseByID(courseID uuid.UUID) (models.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var course models.Course
	if err := r.db.WithContext(ctx).First(&course, "id = ?", courseID).Error; err != nil {
		return models.Course{}, err
	}
	return course, nil
}

// GetCourseWithParticipants preloads the Tutor, Students, and Lessons associations.
func (r *courseRepository) GetCourseWithParticipants(courseID uuid.UUID) (models.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var course models.Course
	if err := r.db.WithContext(ctx).
		Preload("Tutor").
		Preload("Students").
		Preload("Lessons").
		First(&course, "id = ?", courseID).Error; err != nil {
		return models.Course{}, err
	}
	return course, nil
}

// UpdateCourse saves updates to an existing Course.
func (r *courseRepository) UpdateCourse(course *models.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.WithContext(ctx).Save(course).Error
}
func (r *courseRepository) GetCourses(subject, level string, page, limit int) ([]models.Course, error) {
	if limit > 20 {
		limit = 20
	}
	offset := (page - 1) * limit
	var courses []models.Course
	query := r.db.Model(&models.Course{})
	if subject != "" {
		query = query.Where("subject = ?", subject)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}
	err := query.Offset(offset).Limit(limit).
		Preload("Tutor").
		Preload("Students").
		Preload("Lessons").
		Find(&courses).Error
	return courses, err
}
func (r *courseRepository) GetCoursesForUser(userID uuid.UUID) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.
		Preload("Tutor").
		Preload("Students").
		Preload("Lessons").
		Where("tutor_id = ?", userID).
		Or("id IN (SELECT course_id FROM course_students WHERE user_id = ?)", userID).
		Find(&courses).Error
	return courses, err
}
