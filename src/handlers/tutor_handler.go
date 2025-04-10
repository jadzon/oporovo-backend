package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
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

	// Get current time
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Parse start date - default to today if not provided
	startDateStr := c.DefaultQuery("start_date", today.Format("2006-01-02"))
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	// If start date is in the past, set it to today
	if startDate.Before(today) {
		startDate = today
	}

	// Calculate max end date (1 month from start date)
	maxEndDate := startDate.AddDate(0, 1, 0)

	// Parse end date
	endDateStr := c.DefaultQuery("end_date", maxEndDate.Format("2006-01-02"))
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	// Limit end date to max 1 month from start date
	if endDate.After(maxEndDate) {
		endDate = maxEndDate
	}

	fmt.Printf("\n================= AVAILABILITY REQUEST =================\n")
	fmt.Printf("Tutor ID: %s\n", tutorID)
	fmt.Printf("Date range: %s to %s\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	fmt.Printf("Current time: %s\n\n", now.Format("2006-01-02 15:04:05"))

	// 1. Get the tutor's availability (which already accounts for exceptions)
	availableSlots, err := h.App.TAService.GetAvailabilityForDateRange(tutorID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log initial availability grouped by day
	fmt.Printf("INITIAL AVAILABILITY (%d slots):\n", len(availableSlots))
	// Group slots by date for easier reading
	slotsByDate := make(map[string][]string)
	for _, slot := range availableSlots {
		dateKey := slot.Date.Format("02.01") // DD.MM format
		timeSlot := fmt.Sprintf("%s-%s", slot.StartTime, slot.EndTime)
		slotsByDate[dateKey] = append(slotsByDate[dateKey], timeSlot)
	}

	// Print slots by date
	for date, slots := range slotsByDate {
		fmt.Printf("DAY %s: %s\n", date, strings.Join(slots, ", "))
	}
	fmt.Printf("\n")

	// Filter out slots that have already passed on the current day
	var filteredSlots []models.AvailabilitySlot
	fmt.Printf("FILTERING OUT PAST SLOTS (current time: %s):\n", now.Format("15:04"))

	for _, slot := range availableSlots {
		// If this slot is for today, check if it has already passed
		if slot.Date.Year() == now.Year() && slot.Date.Month() == now.Month() && slot.Date.Day() == now.Day() {
			// Parse the slot's start time
			startTimeParts := strings.Split(slot.StartTime, ":")
			if len(startTimeParts) != 2 {
				continue // Skip invalid time format
			}
			startHour, _ := strconv.Atoi(startTimeParts[0])
			startMin, _ := strconv.Atoi(startTimeParts[1])

			// Create a datetime for the slot's start time
			slotStartTime := time.Date(
				now.Year(), now.Month(), now.Day(),
				startHour, startMin, 0, 0, now.Location(),
			)

			// Skip if the slot has already passed
			if slotStartTime.Before(now) {
				fmt.Printf("  SKIPPING: %s %s-%s (already passed)\n",
					slot.Date.Format("02.01"), slot.StartTime, slot.EndTime)
				continue
			} else {
				fmt.Printf("  KEEPING: %s %s-%s (still in future)\n",
					slot.Date.Format("02.01"), slot.StartTime, slot.EndTime)
			}
		}

		filteredSlots = append(filteredSlots, slot)
	}

	// Log availability after filtering past slots by day
	fmt.Printf("\nAVAILABILITY AFTER FILTERING PAST SLOTS (%d slots):\n", len(filteredSlots))
	// Reset and rebuild slots by date
	slotsByDate = make(map[string][]string)
	for _, slot := range filteredSlots {
		dateKey := slot.Date.Format("02.01") // DD.MM format
		timeSlot := fmt.Sprintf("%s-%s", slot.StartTime, slot.EndTime)
		slotsByDate[dateKey] = append(slotsByDate[dateKey], timeSlot)
	}

	// Print slots by date again
	for date, slots := range slotsByDate {
		fmt.Printf("DAY %s: %s\n", date, strings.Join(slots, ", "))
	}
	fmt.Printf("\n")

	// 2. Get the tutor's scheduled lessons
	lessons, err := h.App.LessonService.GetLessonsByTutorIDAndDateRange(tutorID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("FOUND %d LESSONS IN DATE RANGE:\n", len(lessons))
	for _, lesson := range lessons {
		if lesson.Status == models.LessonStatusCancelled {
			fmt.Printf("  CANCELLED LESSON: %s on %s at %s-%s\n",
				lesson.ID,
				lesson.StartTime.Format("02.01"),
				lesson.StartTime.Format("15:04"),
				lesson.EndTime.Format("15:04"))
		} else {
			fmt.Printf("  ACTIVE LESSON: %s on %s at %s-%s\n",
				lesson.ID,
				lesson.StartTime.Format("02.01"),
				lesson.StartTime.Format("15:04"),
				lesson.EndTime.Format("15:04"))
		}
	}
	fmt.Printf("\n")

	// 3. Filter/split availability slots that overlap with lessons
	finalAvailability := filterAvailabilityWithLessons(filteredSlots, lessons)

	// Log final availability by day
	fmt.Printf("FINAL AVAILABILITY AFTER CONSIDERING LESSONS (%d slots):\n", len(finalAvailability))
	// Reset and rebuild slots by date
	slotsByDate = make(map[string][]string)
	for _, slot := range finalAvailability {
		dateKey := slot.Date.Format("02.01") // DD.MM format
		timeSlot := fmt.Sprintf("%s-%s", slot.StartTime, slot.EndTime)
		slotsByDate[dateKey] = append(slotsByDate[dateKey], timeSlot)
	}

	// Print slots by date one last time
	for date, slots := range slotsByDate {
		fmt.Printf("DAY %s: %s\n", date, strings.Join(slots, ", "))
	}
	fmt.Printf("\n================= END OF PROCESSING =================\n\n")

	// IMPORTANT: Ensure the response is never null
	if finalAvailability == nil {
		finalAvailability = []models.AvailabilitySlot{} // Return empty array instead of null
	}

	c.JSON(http.StatusOK, gin.H{
		"tutor_id":        tutorID,
		"available_slots": finalAvailability,
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

// filterAvailabilityWithLessons adjusts availability slots by removing or splitting them
// when they overlap with scheduled lessons
func filterAvailabilityWithLessons(slots []models.AvailabilitySlot, lessons []models.Lesson) []models.AvailabilitySlot {
	if len(lessons) == 0 {
		fmt.Println("No lessons to filter against, returning all slots")
		return slots
	}

	var finalSlots []models.AvailabilitySlot

	fmt.Printf("Processing %d original availability slots against %d lessons\n", len(slots), len(lessons))

	for _, slot := range slots {
		// Parse slot times to datetime objects
		slotDate := slot.Date

		// Parse start time
		startParts := strings.Split(slot.StartTime, ":")
		if len(startParts) != 2 {
			fmt.Printf("WARNING: Invalid start time format: %s, skipping slot\n", slot.StartTime)
			continue
		}
		startHour, _ := strconv.Atoi(startParts[0])
		startMin, _ := strconv.Atoi(startParts[1])
		slotStart := time.Date(
			slotDate.Year(), slotDate.Month(), slotDate.Day(),
			startHour, startMin, 0, 0, slotDate.Location(),
		)

		// Parse end time
		endParts := strings.Split(slot.EndTime, ":")
		if len(endParts) != 2 {
			fmt.Printf("WARNING: Invalid end time format: %s, skipping slot\n", slot.EndTime)
			continue
		}
		endHour, _ := strconv.Atoi(endParts[0])
		endMin, _ := strconv.Atoi(endParts[1])
		slotEnd := time.Date(
			slotDate.Year(), slotDate.Month(), slotDate.Day(),
			endHour, endMin, 0, 0, slotDate.Location(),
		)

		fmt.Printf("Processing slot: %s %s-%s\n", slotDate.Format("2006-01-02"), slot.StartTime, slot.EndTime)

		// Start with the full slot
		availableRanges := []struct {
			start time.Time
			end   time.Time
		}{
			{slotStart, slotEnd},
		}

		// For each lesson, adjust the available ranges
		for _, lesson := range lessons {
			if lesson.Status == models.LessonStatusCancelled {
				fmt.Printf("- Skipping cancelled lesson: %s\n", lesson.ID)
				continue
			}

			fmt.Printf("- Checking lesson: %s (%s - %s)\n",
				lesson.ID,
				lesson.StartTime.Format("2006-01-02 15:04"),
				lesson.EndTime.Format("2006-01-02 15:04"),
			)

			// Create a new list for ranges after processing this lesson
			var newRanges []struct {
				start time.Time
				end   time.Time
			}

			// Process each existing range against this lesson
			for _, r := range availableRanges {
				// No overlap case - keep range as is
				if lesson.EndTime.Before(r.start) || lesson.StartTime.After(r.end) {
					fmt.Printf("  - No overlap with range %s-%s\n",
						r.start.Format("15:04"), r.end.Format("15:04"))
					newRanges = append(newRanges, r)
					continue
				}

				// Handle overlap cases - up to 2 new ranges could be created
				fmt.Printf("  - Overlap detected with range %s-%s\n",
					r.start.Format("15:04"), r.end.Format("15:04"))

				// Part before lesson
				if r.start.Before(lesson.StartTime) {
					newRange := struct {
						start time.Time
						end   time.Time
					}{r.start, lesson.StartTime}
					fmt.Printf("  - Adding range before lesson: %s-%s\n",
						newRange.start.Format("15:04"), newRange.end.Format("15:04"))
					newRanges = append(newRanges, newRange)
				}

				// Part after lesson
				if r.end.After(lesson.EndTime) {
					newRange := struct {
						start time.Time
						end   time.Time
					}{lesson.EndTime, r.end}
					fmt.Printf("  - Adding range after lesson: %s-%s\n",
						newRange.start.Format("15:04"), newRange.end.Format("15:04"))
					newRanges = append(newRanges, newRange)
				}
			}

			// Update availableRanges for next lesson
			availableRanges = newRanges

			// If no ranges left, we can exit early
			if len(availableRanges) == 0 {
				fmt.Printf("  - No available ranges left after processing this lesson\n")
				break
			}
		}

		// Convert remaining time ranges back to AvailabilitySlot format
		for _, r := range availableRanges {
			// Only add slots that are at least 15 minutes long (to avoid tiny gaps)
			minDuration := 15 * time.Minute
			if r.end.Sub(r.start) < minDuration {
				fmt.Printf("Skipping too short range: %s-%s (less than 15 minutes)\n",
					r.start.Format("15:04"), r.end.Format("15:04"))
				continue
			}

			// Format the times back to HH:MM strings
			newSlot := models.AvailabilitySlot{
				Date:      slot.Date,
				StartTime: fmt.Sprintf("%02d:%02d", r.start.Hour(), r.start.Minute()),
				EndTime:   fmt.Sprintf("%02d:%02d", r.end.Hour(), r.end.Minute()),
			}
			fmt.Printf("Adding final slot: %s %s-%s\n",
				newSlot.Date.Format("2006-01-02"), newSlot.StartTime, newSlot.EndTime)
			finalSlots = append(finalSlots, newSlot)
		}
	}

	fmt.Printf("Final slot count after filtering: %d\n", len(finalSlots))
	return finalSlots
}

// GetStudentsForTutor returns all students who have taken lessons with this tutor
func (h *LessonHandler) GetStudentsForTutor(c *gin.Context) {
	tutorIDStr := c.Param("tutorID")
	tutorID, err := uuid.Parse(tutorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tutor ID"})
		return
	}

	students, err := h.App.LessonService.GetStudentsForTutor(tutorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert students to DTOs
	var dtos []models.StudentDTO
	for _, student := range students {
		dtos = append(dtos, student.ToStudentDTO())
	}

	c.JSON(http.StatusOK, dtos)
}
