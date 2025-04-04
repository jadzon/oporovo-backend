package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"vibely-backend/src/models"
)

type LessonRepository interface {
	CreateLesson(lesson *models.Lesson) error
	GetLessonByID(lessonID uuid.UUID) (models.Lesson, error)
	GetLessonWithParticipants(lessonID uuid.UUID) (models.Lesson, error)
	UpdateLesson(lesson *models.Lesson) error
	GetLessonsForUser(userID uuid.UUID) ([]models.Lesson, error)
	GetTutorsForUser(userID uuid.UUID) ([]models.User, error)
}

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db: db}
}

// CreateLesson inserts the Lesson (and any pivot records for Students).
func (r *lessonRepository) CreateLesson(lesson *models.Lesson) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// GORM will create the pivot table entries for `lesson_students` if
	// you have `gorm:"many2many:lesson_students"` on the Lesson model
	// and FullSaveAssociations is configured (by default in GORM v2).
	return r.db.WithContext(ctx).Create(lesson).Error
}

// GetLessonByID retrieves just the Lesson row (without Tutor/Students preloaded).
func (r *lessonRepository) GetLessonByID(lessonID uuid.UUID) (models.Lesson, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lesson models.Lesson
	if err := r.db.WithContext(ctx).
		Where("id = ?", lessonID).
		First(&lesson).Error; err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}

// GetLessonWithParticipants preloads the Tutor and Students associations.
func (r *lessonRepository) GetLessonWithParticipants(lessonID uuid.UUID) (models.Lesson, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lesson models.Lesson
	if err := r.db.WithContext(ctx).
		Preload("Tutor").
		Preload("Students").
		First(&lesson, "id = ?", lessonID).Error; err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}

// UpdateLesson updates the Lesson row and (if needed) pivot table changes.
func (r *lessonRepository) UpdateLesson(lesson *models.Lesson) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// db.Save() will update associated fields, but you may need to confirm
	// FullSaveAssociations is on or handle associations with .Association().
	return r.db.WithContext(ctx).Save(lesson).Error
}
func (r *lessonRepository) GetLessonsForUser(userID uuid.UUID) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.
		Preload("Tutor").
		Preload("Students").
		Where("tutor_id = ?", userID).
		Or("id IN (SELECT lesson_id FROM lesson_students WHERE user_id = ?)", userID).
		Find(&lessons).Error
	return lessons, err
}

// GetTutorsForUser retrieves unique tutors associated with lessons where the user is a student.
func (r *lessonRepository) GetTutorsForUser(userID uuid.UUID) ([]models.User, error) {
	var lessons []models.Lesson
	err := r.db.
		Model(&models.Lesson{}).
		Joins("JOIN lesson_students ON lesson_students.lesson_id = lessons.id").
		Where("lesson_students.user_id = ?", userID).
		Preload("Tutor").
		Find(&lessons).Error
	if err != nil {
		return nil, err
	}

	tutorMap := make(map[uuid.UUID]models.User)
	for _, lesson := range lessons {
		tutorMap[lesson.Tutor.ID] = lesson.Tutor
	}
	tutors := make([]models.User, 0, len(tutorMap))
	for _, tutor := range tutorMap {
		tutors = append(tutors, tutor)
	}
	return tutors, nil
}
