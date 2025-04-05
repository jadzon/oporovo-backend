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
		//router.GET("/api/auth/logout", userHandler.Logout)
	}
	// Protected API Routes
	authorized := router.Group("/api")
	authorized.Use(userHandler.ExtractJWTMiddleware())
	{
		authorized.GET("/user", userHandler.HelloAuthorized)
		authorized.GET("/user/me", userHandler.GetUser)
		router.POST("/lessons", lessonHandler.CreateLesson)
		router.GET("/lessons/:lessonID", lessonHandler.GetLesson)
		router.PATCH("/lessons/:lessonID/confirm", lessonHandler.ConfirmLesson)
		router.PATCH("/lessons/:lessonID/start", lessonHandler.StartLesson)
		router.PATCH("/lessons/:lessonID/complete", lessonHandler.CompleteLesson)
		router.PATCH("/lessons/:lessonID/fail", lessonHandler.FailLesson)
		router.PATCH("/lessons/:lessonID/cancel", lessonHandler.CancelLesson)
		router.PATCH("/lessons/:lessonID/postpone", lessonHandler.PostponeLesson)

		authorized.GET("/user/:userID/lessons", lessonHandler.GetLessonsForUser)
		authorized.GET("/user/:userID/tutors", lessonHandler.GetTutorsForUser)
	}
	return router
}
