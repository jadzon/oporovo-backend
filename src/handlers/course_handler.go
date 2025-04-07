package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"vibely-backend/src/app"
	"vibely-backend/src/models"
)

// CourseHandler holds the reference to the application services.
type CourseHandler struct {
	App *app.Application
}

// NewCourseHandler creates a new CourseHandler.
func NewCourseHandler(app *app.Application) *CourseHandler {
	return &CourseHandler{
		App: app,
	}
}

// createCourseRequest represents the expected payload for creating a course.
type createCourseRequest struct {
	TutorID     string   `json:"tutor_id"`
	StudentIDs  []string `json:"student_ids"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Subject     string   `json:"subject"`
	Level       string   `json:"level"`
}

// CreateCourse handles the creation of a new course.
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req createCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tutorUUID, err := uuid.Parse(req.TutorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tutor_id"})
		return
	}

	tutor, err := h.App.UserService.GetUserByID(tutorUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tutor not found"})
		return
	}

	var students []models.User
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
		students = append(students, student)
	}

	course := models.Course{
		TutorID:     tutorUUID,
		Tutor:       tutor,
		Name:        req.Name,
		Description: req.Description,
		Subject:     req.Subject,
		Level:       req.Level,
		Students:    students,
		CreatedAt:   time.Now(),
	}
	createdCourse, err := h.App.CourseService.CreateCourse(course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdCourse.ToDTO())
}

// GetCourse retrieves a course with its tutor, students, and lessons.
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseIDStr := c.Param("courseID")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}

	course, err := h.App.CourseService.GetCourseWithParticipants(courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	c.JSON(http.StatusOK, course.ToDTO())
}
func (h *CourseHandler) GetCourses(c *gin.Context) {
	subject := c.Query("subject")
	level := c.Query("level")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 20 {
		limit = 20
	}

	courses, err := h.App.CourseService.GetCourses(subject, level, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dtos []models.CourseDTO
	for _, course := range courses {
		dtos = append(dtos, course.ToDTO())
	}
	c.JSON(http.StatusOK, dtos)
}
func (h *CourseHandler) GetCoursesForUser(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	courses, err := h.App.CourseService.GetCoursesForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dtos []models.CourseDTO
	for _, course := range courses {
		dtos = append(dtos, course.ToDTO())
	}

	c.JSON(http.StatusOK, dtos)
}
func (h *CourseHandler) EnrollInCourse(c *gin.Context) {
	courseIDStr := c.Param("courseID")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}

	// Retrieve the authenticated user from context (set by middleware).
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUser, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user"})
		return
	}

	updatedCourse, err := h.App.CourseService.EnrollStudent(courseID, currentUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCourse.ToDTO())
}
