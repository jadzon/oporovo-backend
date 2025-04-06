package models

import (
	"time"

	"github.com/google/uuid"
)

// Course represents a course with a tutor, enrolled students, lessons, and course details.
type Course struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TutorID     uuid.UUID `json:"tutor_id"`
	Tutor       User      `gorm:"foreignKey:TutorID" json:"tutor"`
	Banner      string    `json:"banner"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`   // New price field
	Subject     string    `json:"subject"` // Exposed in the DTO, useful for filtering
	Level       string    `json:"level"`   // Exposed in the DTO, useful for filtering
	Lessons     []Lesson  `gorm:"foreignKey:CourseID" json:"lessons"`
	Students    []User    `gorm:"many2many:course_students" json:"students"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// CourseDTO is the shape returned via API.
type CourseDTO struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Banner      string       `json:"banner"`
	Price       float64      `json:"price"` // Exposed price
	Subject     string       `json:"subject"`
	Level       string       `json:"level"`
	Tutor       TutorDTO     `json:"tutor"`
	Students    []StudentDTO `json:"students"`
	Lessons     []LessonDTO  `json:"lessons"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// ToDTO converts a Course model to a CourseDTO.
func (c *Course) ToDTO() CourseDTO {
	tutorDTO := c.Tutor.ToTutorDTO()

	var studentDTOs []StudentDTO
	for _, student := range c.Students {
		studentDTOs = append(studentDTOs, student.ToStudentDTO())
	}

	var lessonDTOs []LessonDTO
	for _, lesson := range c.Lessons {
		lessonDTOs = append(lessonDTOs, lesson.ToDTO())
	}

	return CourseDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Banner:      c.Banner,
		Price:       c.Price, // map the price
		Subject:     c.Subject,
		Level:       c.Level,
		Tutor:       tutorDTO,
		Students:    studentDTOs,
		Lessons:     lessonDTOs,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
