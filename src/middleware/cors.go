package middleware

import (
	"github.com/gin-gonic/gin"
)

func EnableCORS(c *gin.Context) {
	origin := "http://localhost:5173" // Set to the React app's origin

	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204) // Respond to preflight requests
		return
	}

	c.Next()
}
