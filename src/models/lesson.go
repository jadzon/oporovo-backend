package models

import (
	"time"

	"github.com/google/uuid"
)

// Lesson represents a tutoring session between a tutor and a student
type Lesson struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Foreign keys
	TutorID   uuid.UUID `json:"tutor_id" gorm:"type:uuid;not null"`
	StudentID uuid.UUID `json:"student_id" gorm:"type:uuid;not null"`

	// Relationships
	Tutor   User `json:"tutor" gorm:"foreignKey:TutorID"`
	Student User `json:"student" gorm:"foreignKey:StudentID"`

	// Lesson details
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	Status      string    `json:"status" gorm:"type:varchar(20);not null;default:'scheduled'"`
}
