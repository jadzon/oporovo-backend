package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"vibely-backend/src/app"
	"vibely-backend/src/models"
)

// LessonHandler will hold references to the application context (App),
// which includes the LessonService, UserRepository, etc.
type LessonHandler struct {
	App *app.Application
}

// NewLessonHandler constructs a handler with references to your app’s services/repositories.
func NewLessonHandler(app *app.Application) *LessonHandler {
	return &LessonHandler{
		App: app,
	}
}

type createLessonRequest struct {
	TutorID     string   `json:"tutor_id"`
	StudentIDs  []string `json:"student_ids"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Subject     string   `json:"subject"` // Added subject field
	Level       string   `json:"level"`   // Added level field
	StartTime   string   `json:"start_time"`
	EndTime     string   `json:"end_time"`
}

// CreateLesson schedules a new lesson in "scheduled" state.
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	var req createLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Parse tutor_id
	tutorUUID, err := uuid.Parse(req.TutorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tutor_id"})
		return
	}

	// Convert student_ids to slice of user objects
	var studentUsers []models.User
	for _, sid := range req.StudentIDs {
		studentUUID, err := uuid.Parse(sid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student_id: " + sid})
			return
		}

		student, err := h.App.UserService.GetUserByID(studentUUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found: " + sid})
			return
		}
		studentUsers = append(studentUsers, student)
	}

	// Parse start_time and end_time
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"})
		return
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"})
		return
	}

	// Build the Lesson model
	lesson := models.Lesson{
		TutorID:     tutorUUID,
		Students:    studentUsers,
		Title:       req.Title,
		Description: req.Description,
		Subject:     req.Subject, // Assign subject field
		Level:       req.Level,   // Assign level field
		StartTime:   start,
		EndTime:     end,
		// We'll set Status in the service (to "scheduled").
	}

	// Call service to schedule the lesson
	scheduledLesson, err := h.App.LessonService.ScheduleLesson(lesson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return as DTO
	c.JSON(http.StatusCreated, scheduledLesson.ToDTO())
}
func (h *LessonHandler) GetLesson(c *gin.Context) {
	lessonIDStr := c.Param("lessonID")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson ID"})
		return
	}

	// Retrieve lesson with participants
	lesson, err := h.App.LessonService.GetLessonWithParticipants(lessonID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		return
	}

	c.JSON(http.StatusOK, lesson.ToDTO())
}
func (h *LessonHandler) ConfirmLesson(c *gin.Context) {
	lessonIDStr := c.Param("lessonID")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson ID"})
		return
	}

	lesson, err := h.App.LessonService.ConfirmLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson.ToDTO())
}
func (h *LessonHandler) StartLesson(c *gin.Context) {
	lessonIDStr := c.Param("lessonID")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson ID"})
		return
	}

	lesson, err := h.App.LessonService.StartLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson.ToDTO())
}
func (h *LessonHandler) CompleteLesson(c *gin.Context) {
	lessonIDStr := c.Param("lessonID")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson ID"})
		return
	}

	lesson, err := h.App.LessonService.CompleteLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson.ToDTO())
}
func (h *LessonHandler) FailLesson(c *gin.Context) {
	lessonIDStr := c.Param("lessonID")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson ID"})
		return
	}

	lesson, err := h.App.LessonService.FailLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson.ToDTO())
}
func (h *LessonHandler) CancelLesson(c *gin.Context) {
	lessonIDStr := c.Param("lessonID")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson ID"})
		return
	}

	lesson, err := h.App.LessonService.CancelLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson.ToDTO())
}

type postponeLessonRequest struct {
	NewStartTime string `json:"new_start_time"` // RFC3339 format
	NewEndTime   string `json:"new_end_time"`
}

// PostponeLesson updates the lesson times and sets status back to "scheduled" (or pending).
func (h *LessonHandler) PostponeLesson(c *gin.Context) {
	lessonIDStr := c.Param("lessonID")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson ID"})
		return
	}

	var req postponeLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	newStart, err := time.Parse(time.RFC3339, req.NewStartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid new_start_time"})
		return
	}
	newEnd, err := time.Parse(time.RFC3339, req.NewEndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid new_end_time"})
		return
	}

	lesson, err := h.App.LessonService.PostponeLesson(lessonID, newStart, newEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson.ToDTO())
}
func (h *LessonHandler) GetLessonsForUser(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	lessons, err := h.App.LessonService.GetLessonsForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert lessons to DTOs
	var dtos []models.LessonDTO
	for _, lesson := range lessons {
		dtos = append(dtos, lesson.ToDTO())
	}

	c.JSON(http.StatusOK, dtos)
}

// GetTutorsForUser returns a unique list of tutors for lessons where the user is a student.
func (h *LessonHandler) GetTutorsForUser(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	tutors, err := h.App.LessonService.GetTutorsForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert tutors to UserDTOs
	var dtos []models.TutorDTO
	for _, tutor := range tutors {
		dtos = append(dtos, tutor.ToTutorDTO())
	}

	c.JSON(http.StatusOK, dtos)
}
