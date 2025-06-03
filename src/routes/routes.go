package routes

import (
	"github.com/gin-gonic/gin"
	"vibely-backend/src/app"
	"vibely-backend/src/handlers"
	"vibely-backend/src/middleware"
)

func Setup(app *app.Application) *gin.Engine {
	router := gin.Default()

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(app)
	lessonHandler := handlers.NewLessonHandler(app)
	tutorHandler := handlers.NewTutorHandler(app)
	courseHandler := handlers.NewCourseHandler(app)

	// Apply Global Middleware
	router.Use(middleware.EnableCORS)

	// Public API Routes
	publicAPI := router.Group("/api")
	{
		// Authentication Routes
		//publicAPI.GET("/hello", userHandler.Hello)
		publicAPI.GET("/tutors", tutorHandler.GetTutors)
		router.GET("/api/auth/discord", userHandler.DiscordLogin)
		router.GET("/api/auth/discord/callback", userHandler.DiscordCallback)
		//TODO:
		router.POST("/api/token/refresh-token", userHandler.RefreshTokens)

		publicAPI.GET("/tutors/:tutorID/availability", tutorHandler.GetAvailability)
		//router.GET("/api/auth/logout", userHandler.Logout)
	}
	// Protected API Routes
	authorized := router.Group("/api")
	authorized.Use(userHandler.ExtractJWTMiddleware())
	{
		authorized.GET("/user/me", userHandler.GetUser)
		authorized.GET("/user/:userID", userHandler.GetUserById)
		authorized.POST("/lessons", lessonHandler.CreateLesson)
		authorized.GET("/lessons/:lessonID", lessonHandler.GetLesson)
		authorized.PATCH("/lessons/:lessonID/confirm", lessonHandler.ConfirmLesson)
		authorized.PATCH("/lessons/:lessonID/start", lessonHandler.StartLesson)
		authorized.PATCH("/lessons/:lessonID/complete", lessonHandler.CompleteLesson)
		authorized.PATCH("/lessons/:lessonID/fail", lessonHandler.FailLesson)
		authorized.PATCH("/lessons/:lessonID/cancel", lessonHandler.CancelLesson)
		authorized.PATCH("/lessons/:lessonID/postpone", lessonHandler.PostponeLesson)

		authorized.GET("/user/:userID/lessons", lessonHandler.GetLessonsForUser)
		authorized.GET("/user/:userID/tutors", lessonHandler.GetTutorsForUser)

		authorized.POST("/courses", courseHandler.CreateCourse)
		authorized.GET("/courses/:courseID", courseHandler.GetCourse)
		authorized.GET("/courses", courseHandler.GetCourses)
		authorized.GET("/user/:userID/courses", courseHandler.GetCoursesForUser)
		authorized.POST("/courses/:courseID/enroll", courseHandler.EnrollInCourse)

		authorized.POST("/tutors/:tutorID/weekly-schedules", tutorHandler.AddWeeklySchedule)
		authorized.GET("/tutors/:tutorID/weekly-schedules", tutorHandler.GetWeeklySchedule)
		authorized.PUT("/weekly-schedules/:scheduleID", tutorHandler.UpdateWeeklySchedule)
		authorized.DELETE("/weekly-schedules/:scheduleID", tutorHandler.DeleteWeeklySchedule)

		authorized.POST("/tutors/:tutorID/exceptions", tutorHandler.AddException)
		authorized.GET("/tutors/:tutorID/exceptions", tutorHandler.GetExceptions)

		authorized.GET("/tutors/:tutorID/students", lessonHandler.GetStudentsForTutor)

	}
	return router
}
