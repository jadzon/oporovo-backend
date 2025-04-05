package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"vibely-backend/src/app"
	"vibely-backend/src/models"
)

// TutorHandler handles endpoints related to tutors.
type TutorHandler struct {
	App *app.Application
}

// NewTutorHandler initializes a new TutorHandler.
func NewTutorHandler(app *app.Application) *TutorHandler {
	return &TutorHandler{App: app}
}

// TutorFilters is used to pass filter criteria to the service layer.
type TutorFilters struct {
	Page    int
	Limit   int
	Subject string
	Level   string
}

// GetTutors handles GET /api/tutors?subject=...&level=...&page=...&limit=...
func (h *TutorHandler) GetTutors(c *gin.Context) {
	// Parse pagination parameters.
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	// Enforce a maximum of 20 tutors per page.
	if err != nil || limit < 1 || limit > 20 {
		limit = 20
	}

	// Parse filter parameters.
	subject := c.Query("subject")
	level := c.Query("level")

	// Create a filters struct.
	filters := models.TutorFilters{
		Page:    page,
		Limit:   limit,
		Subject: subject,
		Level:   level,
	}

	// Get tutors from the service layer.
	tutors, total, err := h.App.UserService.GetTutors(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tutors"})
		return
	}

	// Convert each tutor model to a TutorDTO.
	var tutorDTOs []models.TutorDTO
	for _, tutor := range tutors {
		tutorDTOs = append(tutorDTOs, tutor.ToTutorDTO())
	}

	c.JSON(http.StatusOK, gin.H{
		"tutors": tutorDTOs,
		"page":   page,
		"limit":  limit,
		"total":  total,
	})
}
