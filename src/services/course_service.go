package services

import (
	"errors"

	"github.com/google/uuid"
	"vibely-backend/src/models"
	"vibely-backend/src/repositories"
)

// CourseService defines business logic for courses.
type CourseService interface {
	CreateCourse(course models.Course) (models.Course, error)
	GetCourseByID(courseID uuid.UUID) (models.Course, error)
	GetCourseWithParticipants(courseID uuid.UUID) (models.Course, error)
	UpdateCourse(course models.Course) (models.Course, error)
	GetCourses(subject, level string, page, limit int) ([]models.Course, error)
	GetCoursesForUser(userID uuid.UUID) ([]models.Course, error)
}

type courseService struct {
	repo repositories.CourseRepository
}

// NewCourseService creates a new instance of CourseService.
func NewCourseService(repo repositories.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

// CreateCourse validates and creates a new course.
func (s *courseService) CreateCourse(course models.Course) (models.Course, error) {
	if course.TutorID == uuid.Nil {
		return models.Course{}, errors.New("tutor_id is required")
	}
	if course.Name == "" {
		return models.Course{}, errors.New("course name is required")
	}

	// Additional validations (e.g., check for students) can be added here.

	if err := s.repo.CreateCourse(&course); err != nil {
		return models.Course{}, err
	}
	return course, nil
}

// GetCourseByID retrieves a course without preloading associations.
func (s *courseService) GetCourseByID(courseID uuid.UUID) (models.Course, error) {
	return s.repo.GetCourseByID(courseID)
}

// GetCourseWithParticipants retrieves a course along with tutor, students, and lessons.
func (s *courseService) GetCourseWithParticipants(courseID uuid.UUID) (models.Course, error) {
	return s.repo.GetCourseWithParticipants(courseID)
}

// UpdateCourse updates an existing course.
func (s *courseService) UpdateCourse(course models.Course) (models.Course, error) {
	if err := s.repo.UpdateCourse(&course); err != nil {
		return models.Course{}, err
	}
	return course, nil
}
func (s *courseService) GetCourses(subject, level string, page, limit int) ([]models.Course, error) {
	return s.repo.GetCourses(subject, level, page, limit)
}
func (s *courseService) GetCoursesForUser(userID uuid.UUID) ([]models.Course, error) {
	return s.repo.GetCoursesForUser(userID)
}
