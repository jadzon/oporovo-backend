package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"vibely-backend/src/models"
)

type TutorAvailabilityRepository interface {
	// Weekly schedule methods
	CreateWeeklySchedule(schedule models.TutorWeeklySchedule) (models.TutorWeeklySchedule, error)
	GetWeeklyScheduleByID(scheduleID uuid.UUID) (models.TutorWeeklySchedule, error)
	GetWeeklySchedulesByTutorID(tutorID uuid.UUID) ([]models.TutorWeeklySchedule, error)
	UpdateWeeklySchedule(schedule models.TutorWeeklySchedule) error
	DeleteWeeklySchedule(scheduleID uuid.UUID) error

	// Exception methods
	CreateException(exception models.TutorScheduleException) (models.TutorScheduleException, error)
	GetExceptionByID(exceptionID uuid.UUID) (models.TutorScheduleException, error)
	GetExceptionsByTutorID(tutorID uuid.UUID, startDate, endDate time.Time) ([]models.TutorScheduleException, error)
	UpdateException(exception models.TutorScheduleException) error
	DeleteException(exceptionID uuid.UUID) error
}

type tutorAvailabilityRepository struct {
	db *gorm.DB
}

func NewTutorAvailabilityRepository(db *gorm.DB) TutorAvailabilityRepository {
	return &tutorAvailabilityRepository{db: db}
}

// Weekly schedule implementations
func (r *tutorAvailabilityRepository) CreateWeeklySchedule(schedule models.TutorWeeklySchedule) (models.TutorWeeklySchedule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(&schedule).Error; err != nil {
		return models.TutorWeeklySchedule{}, err
	}

	return schedule, nil
}

func (r *tutorAvailabilityRepository) GetWeeklyScheduleByID(scheduleID uuid.UUID) (models.TutorWeeklySchedule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var schedule models.TutorWeeklySchedule
	if err := r.db.WithContext(ctx).Where("id = ?", scheduleID).First(&schedule).Error; err != nil {
		return models.TutorWeeklySchedule{}, err
	}

	return schedule, nil
}

func (r *tutorAvailabilityRepository) GetWeeklySchedulesByTutorID(tutorID uuid.UUID) ([]models.TutorWeeklySchedule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var schedules []models.TutorWeeklySchedule
	if err := r.db.WithContext(ctx).Where("tutor_id = ?", tutorID).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *tutorAvailabilityRepository) UpdateWeeklySchedule(schedule models.TutorWeeklySchedule) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.WithContext(ctx).Save(&schedule).Error
}

func (r *tutorAvailabilityRepository) DeleteWeeklySchedule(scheduleID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.WithContext(ctx).Delete(&models.TutorWeeklySchedule{}, "id = ?", scheduleID).Error
}

// Exception implementations
func (r *tutorAvailabilityRepository) CreateException(exception models.TutorScheduleException) (models.TutorScheduleException, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(&exception).Error; err != nil {
		return models.TutorScheduleException{}, err
	}

	return exception, nil
}

func (r *tutorAvailabilityRepository) GetExceptionByID(exceptionID uuid.UUID) (models.TutorScheduleException, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exception models.TutorScheduleException
	if err := r.db.WithContext(ctx).Where("id = ?", exceptionID).First(&exception).Error; err != nil {
		return models.TutorScheduleException{}, err
	}

	return exception, nil
}

func (r *tutorAvailabilityRepository) GetExceptionsByTutorID(tutorID uuid.UUID, startDate, endDate time.Time) ([]models.TutorScheduleException, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exceptions []models.TutorScheduleException
	query := r.db.WithContext(ctx).Where("tutor_id = ?", tutorID)

	// Apply date range filter if provided
	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Find(&exceptions).Error; err != nil {
		return nil, err
	}

	return exceptions, nil
}

func (r *tutorAvailabilityRepository) UpdateException(exception models.TutorScheduleException) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.WithContext(ctx).Save(&exception).Error
}

func (r *tutorAvailabilityRepository) DeleteException(exceptionID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.WithContext(ctx).Delete(&models.TutorScheduleException{}, "id = ?", exceptionID).Error
}
