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
type Lesson struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	TutorID uuid.UUID `json:"tutor_id"`
	Tutor   User      `gorm:"foreignKey:TutorID"`

	Students []User `gorm:"many2many:lesson_students" json:"students"`

	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// LessonDTO is the shape you might return via API.
// It includes the tutor’s info and the list of students’ info in a minimal form.
type LessonDTO struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
	Status      string       `json:"status"`
	Tutor       TutorDTO     `json:"tutor"`
	Students    []StudentDTO `json:"students"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Convert a Lesson model to LessonDTO.
func (l Lesson) ToDTO() LessonDTO {
	tutorDTO := l.Tutor.ToTutorDTO()

	// Convert all students to user DTO
	var students []StudentDTO
	for _, s := range l.Students {
		students = append(students, s.ToStudentDTO())
	}

	return LessonDTO{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		StartTime:   l.StartTime,
		EndTime:     l.EndTime,
		Status:      l.Status,
		Tutor:       tutorDTO,
		Students:    students,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}
