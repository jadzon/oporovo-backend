package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"vibely-backend/src/models"
	"vibely-backend/src/repositories"
)

type LessonService interface {
	// Basic create & read
	ScheduleLesson(lesson models.Lesson) (models.Lesson, error)
	GetLessonByID(lessonID uuid.UUID) (models.Lesson, error)
	GetLessonWithParticipants(lessonID uuid.UUID) (models.Lesson, error)

	// Status transitions
	ConfirmLesson(lessonID uuid.UUID) (models.Lesson, error)
	StartLesson(lessonID uuid.UUID) (models.Lesson, error)
	CompleteLesson(lessonID uuid.UUID) (models.Lesson, error)
	FailLesson(lessonID uuid.UUID) (models.Lesson, error)
	CancelLesson(lessonID uuid.UUID) (models.Lesson, error)
	PostponeLesson(lessonID uuid.UUID, newStart, newEnd time.Time) (models.Lesson, error)
}

type lessonService struct {
	repo repositories.LessonRepository
}

func NewLessonService(repo repositories.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

// ScheduleLesson: create a new Lesson with status "scheduled".
func (s *lessonService) ScheduleLesson(lesson models.Lesson) (models.Lesson, error) {
	// Basic validations
	if lesson.TutorID == uuid.Nil {
		return models.Lesson{}, errors.New("tutor_id is required")
	}
	if len(lesson.Students) == 0 {
		return models.Lesson{}, errors.New("at least one student is required")
	}

	lesson.Status = models.LessonStatusScheduled
	if err := s.repo.CreateLesson(&lesson); err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}

// GetLessonByID: minimal retrieval, no preloads
func (s *lessonService) GetLessonByID(lessonID uuid.UUID) (models.Lesson, error) {
	return s.repo.GetLessonByID(lessonID)
}

// GetLessonWithParticipants: retrieve a lesson with the Tutor and Students.
func (s *lessonService) GetLessonWithParticipants(lessonID uuid.UUID) (models.Lesson, error) {
	return s.repo.GetLessonWithParticipants(lessonID)
}

// ConfirmLesson: sets the Lesson status to "confirmed".
func (s *lessonService) ConfirmLesson(lessonID uuid.UUID) (models.Lesson, error) {
	lesson, err := s.repo.GetLessonByID(lessonID)
	if err != nil {
		return models.Lesson{}, err
	}
	if lesson.Status != models.LessonStatusScheduled {
		return models.Lesson{}, errors.New("only scheduled lessons can be confirmed")
	}
	lesson.Status = models.LessonStatusConfirmed
	if err := s.repo.UpdateLesson(&lesson); err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}

// StartLesson: sets the Lesson status to "in_progress".
func (s *lessonService) StartLesson(lessonID uuid.UUID) (models.Lesson, error) {
	lesson, err := s.repo.GetLessonByID(lessonID)
	if err != nil {
		return models.Lesson{}, err
	}
	if lesson.Status != models.LessonStatusConfirmed {
		return models.Lesson{}, errors.New("lesson must be confirmed before starting")
	}
	lesson.Status = models.LessonStatusInProgress
	if err := s.repo.UpdateLesson(&lesson); err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}

// CompleteLesson: sets the Lesson status to "done".
func (s *lessonService) CompleteLesson(lessonID uuid.UUID) (models.Lesson, error) {
	lesson, err := s.repo.GetLessonByID(lessonID)
	if err != nil {
		return models.Lesson{}, err
	}
	if lesson.Status != models.LessonStatusInProgress {
		return models.Lesson{}, errors.New("lesson must be in progress to be completed")
	}
	lesson.Status = models.LessonStatusDone
	if err := s.repo.UpdateLesson(&lesson); err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}

// FailLesson: sets the Lesson status to "failed".
func (s *lessonService) FailLesson(lessonID uuid.UUID) (models.Lesson, error) {
	lesson, err := s.repo.GetLessonByID(lessonID)
	if err != nil {
		return models.Lesson{}, err
	}
	if lesson.Status != models.LessonStatusInProgress &&
		lesson.Status != models.LessonStatusConfirmed {
		return models.Lesson{}, errors.New("cannot fail a lesson from the current status")
	}
	lesson.Status = models.LessonStatusFailed
	if err := s.repo.UpdateLesson(&lesson); err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}

// CancelLesson: sets the Lesson status to "cancelled".
func (s *lessonService) CancelLesson(lessonID uuid.UUID) (models.Lesson, error) {
	lesson, err := s.repo.GetLessonByID(lessonID)
	if err != nil {
		return models.Lesson{}, err
	}
	lesson.Status = models.LessonStatusCancelled
	if err := s.repo.UpdateLesson(&lesson); err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}

// PostponeLesson: updates times, sets status back to "scheduled" (or custom flow).
func (s *lessonService) PostponeLesson(lessonID uuid.UUID, newStart, newEnd time.Time) (models.Lesson, error) {
	lesson, err := s.repo.GetLessonByID(lessonID)
	if err != nil {
		return models.Lesson{}, err
	}
	if lesson.Status != models.LessonStatusScheduled &&
		lesson.Status != models.LessonStatusConfirmed {
		return models.Lesson{}, errors.New("lesson can only be postponed if it's scheduled or confirmed")
	}

	lesson.StartTime = newStart
	lesson.EndTime = newEnd
	lesson.Status = models.LessonStatusScheduled // or "pending_confirmation"

	if err := s.repo.UpdateLesson(&lesson); err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}
