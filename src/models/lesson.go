package models

import (
	"time"

	"github.com/google/uuid"
)

// Lesson Status Constants
const (
	LessonStatusScheduled  = "scheduled"
	LessonStatusConfirmed  = "confirmed"
	LessonStatusInProgress = "in_progress"
	LessonStatusDone       = "done"
	LessonStatusFailed     = "failed"
	LessonStatusCancelled  = "cancelled"
)

// Lesson model
//   - Single Tutor (TutorID / Tutor field)
//   - Many Students (Students field via a pivot table)
//   - Optionally associated with a Course
type Lesson struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	TutorID uuid.UUID `json:"tutor_id"`
	Tutor   User      `gorm:"foreignKey:TutorID"`

	Students []User `gorm:"many2many:lesson_students" json:"students"`

	Title       string    `json:"title"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	Level       string    `json:"level"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`

	// Optional association with a Course.
	// CourseID is a pointer so it can be nil when there is no associated course.
	CourseID *uuid.UUID `json:"course_id,omitempty"`
	// Course is also optional and omitted from JSON if nil.
	Course *Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// CourseSummaryDTO contains minimal course details.
type CourseSummaryDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// LessonDTO is the shape you might return via API.
// It includes the tutor’s info and the list of students’ info in a minimal form,
// and optionally a summary of the associated course.
type LessonDTO struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Subject     string       `json:"subject"`
	Level       string       `json:"level"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
	Status      string       `json:"status"`
	Tutor       TutorDTO     `json:"tutor"`
	Students    []StudentDTO `json:"students"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`

	Course *CourseSummaryDTO `json:"course,omitempty"`
}

// ToDTO converts a Lesson model to LessonDTO.
// If the lesson has a non-nil CourseID and the Course field is preloaded with a name,
// it returns a CourseSummaryDTO.
func (l Lesson) ToDTO() LessonDTO {
	tutorDTO := l.Tutor.ToTutorDTO()

	// Convert all students to DTO.
	var students []StudentDTO
	for _, s := range l.Students {
		students = append(students, s.ToStudentDTO())
	}

	// Prepare CourseSummaryDTO only if CourseID is set, non-nil, and Course data is available.
	var courseSummary *CourseSummaryDTO
	if l.CourseID != nil && *l.CourseID != uuid.Nil && l.Course != nil && l.Course.ID != uuid.Nil && l.Course.Name != "" {
		courseSummary = &CourseSummaryDTO{
			ID:   l.Course.ID,
			Name: l.Course.Name,
		}
	}

	return LessonDTO{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		Subject:     l.Subject,
		Level:       l.Level,
		StartTime:   l.StartTime,
		EndTime:     l.EndTime,
		Status:      l.Status,
		Tutor:       tutorDTO,
		Students:    students,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
		Course:      courseSummary,
	}
}
