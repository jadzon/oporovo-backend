package models

import "time"
import "github.com/google/uuid"

const (
	UserRoleStudent = "student"
	UserRoleTutor   = "tutor"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	Username    string    `json:"username" gorm:"uniqueIndex;not null"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth string    `json:"date_of_birth"`
	Role        string    `json:"role" gorm:"type:varchar(20);not null;default:'student'"`
	Description string    `json:"description" gorm:"type:text"`
}
