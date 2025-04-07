package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"

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

// AddWeeklySchedule adds a new weekly availability slot for a tutor
func (h *TutorHandler) AddWeeklySchedule(c *gin.Context) {
	// Get tutor ID from URL parameter
	tutorIDStr := c.Param("tutorID")
	tutorID, err := uuid.Parse(tutorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tutor ID format"})
		return
	}

	// Parse request body
	type ScheduleRequest struct {
		DayOfWeek int    `json:"day_of_week" binding:"required,min=0,max=6"`
		StartTime string `json:"start_time" binding:"required"`
		EndTime   string `json:"end_time" binding:"required"`
	}

	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	schedule, err := h.App.TAService.CreateWeeklySchedule(
		tutorID,
		req.DayOfWeek,
		req.StartTime,
		req.EndTime,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Weekly schedule added successfully",
		"schedule": schedule,
	})
}

// GetWeeklySchedule retrieves all weekly schedule slots for a tutor
func (h *TutorHandler) GetWeeklySchedule(c *gin.Context) {
	// Get tutor ID from URL parameter
	tutorIDStr := c.Param("tutorID")
	tutorID, err := uuid.Parse(tutorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tutor ID format"})
		return
	}

	// Call service
	schedules, err := h.App.TAService.GetWeeklySchedulesByTutorID(tutorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schedules": schedules,
	})
}

// UpdateWeeklySchedule updates an existing weekly schedule slot
func (h *TutorHandler) UpdateWeeklySchedule(c *gin.Context) {
	// Get schedule ID from URL parameter
	scheduleIDStr := c.Param("scheduleID")
	scheduleID, err := uuid.Parse(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID format"})
		return
	}

	// Parse request body
	type ScheduleRequest struct {
		DayOfWeek int    `json:"day_of_week" binding:"required,min=0,max=6"`
		StartTime string `json:"start_time" binding:"required"`
		EndTime   string `json:"end_time" binding:"required"`
	}

	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	schedule, err := h.App.TAService.UpdateWeeklySchedule(
		scheduleID,
		req.DayOfWeek,
		req.StartTime,
		req.EndTime,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Weekly schedule updated successfully",
		"schedule": schedule,
	})
}

// DeleteWeeklySchedule deletes a weekly schedule slot
func (h *TutorHandler) DeleteWeeklySchedule(c *gin.Context) {
	// Get schedule ID from URL parameter
	scheduleIDStr := c.Param("scheduleID")
	scheduleID, err := uuid.Parse(scheduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID format"})
		return
	}

	// Call service
	err = h.App.TAService.DeleteWeeklySchedule(scheduleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Weekly schedule deleted successfully",
	})
}

// AddException adds a new exception to the tutor's schedule
func (h *TutorHandler) AddException(c *gin.Context) {
	// Get tutor ID from URL parameter
	tutorIDStr := c.Param("tutorID")
	tutorID, err := uuid.Parse(tutorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tutor ID format"})
		return
	}

	// Parse request body
	type ExceptionRequest struct {
		Date      string `json:"date" binding:"required"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		IsRemoval bool   `json:"is_removal"`
	}

	var req ExceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Call service
	exception, err := h.App.TAService.AddException(
		tutorID,
		date,
		req.StartTime,
		req.EndTime,
		req.IsRemoval,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Schedule exception added successfully",
		"exception": exception,
	})
}

// GetExceptions retrieves all exceptions for a tutor in a date range
func (h *TutorHandler) GetExceptions(c *gin.Context) {
	// Get tutor ID from URL parameter
	tutorIDStr := c.Param("tutorID")
	tutorID, err := uuid.Parse(tutorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tutor ID format"})
		return
	}

	// Parse date range from query parameters
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	var startDate, endDate time.Time

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
	}

	// Call service
	exceptions, err := h.App.TAService.GetExceptionsByTutorID(tutorID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exceptions": exceptions,
	})
}

// GetAvailability retrieves a tutor's availability for a date range
// GetAvailability retrieves a tutor's availability for a date range
func (h *TutorHandler) GetAvailability(c *gin.Context) {
	// Get tutor ID from URL parameter
	tutorIDStr := c.Param("tutorID")
	tutorID, err := uuid.Parse(tutorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tutor ID format"})
		return
	}

	// Parse date range from query parameters (required)
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	// Call service
	slots, err := h.App.TAService.GetAvailabilityForDateRange(tutorID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("GetAvailabilityForDateRange called with:", tutorID)
	fmt.Println("Date range:", startDate.Format("2006-01-02"), "to", endDate.Format("2006-01-02"))

	// [... existing code ...]

	// Iterate through each day in the date range
	for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
		dayOfWeek := int(date.Weekday())
		fmt.Printf("Processing date: %s, day of week: %d\n", date.Format("2006-01-02"), dayOfWeek)

		// [... rest of your code ...]
	}
	
	// IMPORTANT FIX: Ensure slots is never null in the response
	if slots == nil {
		slots = []models.AvailabilitySlot{} // Return empty array instead of null
	}
	fmt.Println("AVAILABLE SLOTS")
	fmt.Println(slots)
	c.JSON(http.StatusOK, gin.H{
		"tutor_id":        tutorID,
		"available_slots": slots,
	})
}
