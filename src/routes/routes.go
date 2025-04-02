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

	// Apply Global Middleware
	router.Use(middleware.EnableCORS)

	// Public API Routes
	publicAPI := router.Group("/api")
	{
		// Authentication Routes
		publicAPI.POST("/login", userHandler.Login)
		publicAPI.GET("/hello", userHandler.Hello)
		router.GET("/api/auth/discord", userHandler.DiscordLogin)
		router.GET("/api/auth/discord/callback", userHandler.DiscordCallback)
		//TODO:
		//router.GET("/api/auth/refresh", userHandler.RefreshToken)
		//router.GET("/api/auth/logout", userHandler.Logout)
	}
	// Protected API Routes
	authorized := router.Group("/api")
	authorized.Use(userHandler.ExtractJWTMiddleware())
	{
		authorized.GET("/user", userHandler.HelloAuthorized)
		// Add more protected routes here
	}
	//{
	//	// User Profile Route
	//	protectedAPI.GET("/user/profile", userHandler.GetUserProfile)
	//
	//	// Video Interaction Routes
	//	videoGroup := protectedAPI.Group("/video")
	//	{
	//		videoGroup.POST("/:id/like", videoHandler.LikeVideo)
	//		videoGroup.POST("/:id/dislike", videoHandler.DislikeVideo)
	//		videoGroup.GET("/:id", videoHandler.GetVideo)
	//	}
	//}

	return router
}
