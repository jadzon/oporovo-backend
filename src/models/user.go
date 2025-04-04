package models

import (
	"fmt"
	"github.com/lib/pq"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	UserRoleStudent = "student"
	UserRoleTutor   = "tutor"
)

// User model
type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	DiscordID     string    `json:"discord_id" gorm:"uniqueIndex;not null"`
	Discriminator string    `json:"discriminator"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Email         string    `json:"email" gorm:"uniqueIndex;not null"`
	Username      string    `json:"username" gorm:"uniqueIndex;not null"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	DateOfBirth   string    `json:"date_of_birth"`
	Role          string    `json:"role" gorm:"type:varchar(20);not null;default:'student'"`
	Description   string    `json:"description" gorm:"type:text"`
	Avatar        string    `json:"avatar"`

	// New fields:
	Rating float64        `json:"rating"` // For tutors, e.g. 4.5 out of 5
	Levels pq.StringArray `json:"levels" gorm:"type:text[]"`
}

// RetrieveAvatarURL constructs the URL to the user's Discord avatar.
func (u *User) RetrieveAvatarURL() string {
	if u.Avatar != "" {
		ext := "png"
		// If the avatar hash starts with "a_", it's animated.
		if strings.HasPrefix(u.Avatar, "a_") {
			ext = "gif"
		}
		return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.%s", u.DiscordID, u.Avatar, ext)
	}

	// If no custom avatar exists, use the default avatar.
	discriminator, err := strconv.Atoi(u.Discriminator)
	if err != nil {
		discriminator = 0
	}
	return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", discriminator%5)
}

// StudentDTO includes fields for both students and tutors.
// Now includes Rating, so even students have a rating.
type StudentDTO struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	DiscordID   string    `json:"discord_id"`
	DateOfBirth string    `json:"date_of_birth"`
	Rating      float64   `json:"rating"`
}

// TutorDTO extends StudentDTO with tutor-specific fields.
type TutorDTO struct {
	StudentDTO
	Levels []string `json:"levels"`
}

// ToStudentDTO converts a User model to a StudentDTO.
func (u *User) ToStudentDTO() StudentDTO {
	return StudentDTO{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		Avatar:      u.RetrieveAvatarURL(),
		Role:        u.Role,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Description: u.Description,
		CreatedAt:   u.CreatedAt,
		DiscordID:   u.DiscordID,
		DateOfBirth: u.DateOfBirth,
		Rating:      u.Rating,
	}
}

// ToTutorDTO converts a User model to a TutorDTO.
func (u *User) ToTutorDTO() TutorDTO {
	return TutorDTO{
		StudentDTO: u.ToStudentDTO(),
		Levels:     u.Levels,
	}
}
