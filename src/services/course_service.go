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
	EnrollStudent(courseID uuid.UUID, student models.User) (models.Course, error)
}

type courseService struct {
	courseRepo    repositories.CourseRepository
	lessonService LessonService
}

// NewCourseService creates a new instance of CourseService.
func NewCourseService(courseRepo repositories.CourseRepository, lessonService LessonService) CourseService {
	return &courseService{courseRepo: courseRepo, lessonService: lessonService}
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

	if err := s.courseRepo.CreateCourse(&course); err != nil {
		return models.Course{}, err
	}
	return course, nil
}

// GetCourseByID retrieves a course without preloading associations.
func (s *courseService) GetCourseByID(courseID uuid.UUID) (models.Course, error) {
	return s.courseRepo.GetCourseByID(courseID)
}

// GetCourseWithParticipants retrieves a course along with tutor, students, and lessons.
func (s *courseService) GetCourseWithParticipants(courseID uuid.UUID) (models.Course, error) {
	return s.courseRepo.GetCourseWithParticipants(courseID)
}

// UpdateCourse updates an existing course.
func (s *courseService) UpdateCourse(course models.Course) (models.Course, error) {
	if err := s.courseRepo.UpdateCourse(&course); err != nil {
		return models.Course{}, err
	}
	return course, nil
}
func (s *courseService) GetCourses(subject, level string, page, limit int) ([]models.Course, error) {
	return s.courseRepo.GetCourses(subject, level, page, limit)
}
func (s *courseService) GetCoursesForUser(userID uuid.UUID) ([]models.Course, error) {
	return s.courseRepo.GetCoursesForUser(userID)
}
func (s *courseService) EnrollStudent(courseID uuid.UUID, student models.User) (models.Course, error) {
	// Retrieve the course with its participants and lessons.
	course, err := s.courseRepo.GetCourseWithParticipants(courseID)
	if err != nil {
		return models.Course{}, err
	}

	// Check if student is already enrolled in the course.
	for _, enrolled := range course.Students {
		if enrolled.ID == student.ID {
			return models.Course{}, errors.New("student already enrolled in this course")
		}
	}

	// Enroll student in course.
	if err := s.courseRepo.EnrollStudent(courseID, student); err != nil {
		return models.Course{}, err
	}

	// Enroll student in every lesson of the course.
	// Here we assume s.courseRepo is an instance of courseRepository injected in the course service.
	for _, lesson := range course.Lessons {
		if err := s.lessonService.EnrollStudent(lesson.ID, student); err != nil {
			// Optionally, you might want to roll back the course enrollment if lesson enrollment fails.
			return models.Course{}, err
		}
	}

	// Return the updated course with its associated tutor, students, and lessons.
	return s.courseRepo.GetCourseWithParticipants(courseID)
}
