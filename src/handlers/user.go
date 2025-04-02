package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"vibely-backend/src/app"
)

type UserHandler struct {
	App *app.Application
}

func NewUserHandler(app *app.Application) *UserHandler {
	return &UserHandler{App: app}
}
func (h *UserHandler) Login(c *gin.Context) {}

func (h *UserHandler) ExtractJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getAccessTokenFromCookie(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		user, err := h.App.UserService.GetUserFromAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
func setAccessTokenCookie(c *gin.Context, token, domain string) {
	c.SetCookie(
		"authToken", // Name
		token,       // Value
		900,         // Max-Age in seconds (15 minutes)
		"/",         // Path
		domain,      // Domain
		false,       // Secure (set to true in production with HTTPS)
		true,        // HttpOnly
	)
}

// setRefreshTokenCookie sets the refreshToken as an HTTP-only cookie
func setRefreshTokenCookie(c *gin.Context, token, domain string) {
	c.SetCookie(
		"refreshToken",
		token,
		604800,
		"/",
		domain,
		false,
		true,
	)
}
func getAccessTokenFromCookie(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("authToken")
	if err != nil {
		return "", errors.New("JWT cookie not found")
	}
	return cookie, nil
}

func getRefreshTokenFromCookie(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		return "", errors.New("refresh token not found")
	}
	return cookie, nil
}
func (h *UserHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello"})
}
