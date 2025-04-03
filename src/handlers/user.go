package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"os"
	"vibely-backend/src/app"
	"vibely-backend/src/models"
)

type UserHandler struct {
	App          *app.Application
	oauthConfig  *oauth2.Config
	frontendURL  string
	backendURL   string
	cookieDomain string
}
type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
}

func NewUserHandler(app *app.Application) *UserHandler {
	// Discord OAuth configuration
	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("DISCORD_REDIRECT_URL"),
		Scopes:       []string{"identify", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
	return &UserHandler{
		App:          app,
		oauthConfig:  oauthConfig,
		frontendURL:  app.Config.FrontendUrl,
		backendURL:   app.Config.BackendUrl,
		cookieDomain: app.Config.CookieDomain,
	}
}
func (h *UserHandler) DiscordLogin(c *gin.Context) {
	//TODO
	// Generate a random state to prevent CSRF attacks
	state := h.App.AuthService.GenerateRandomState()

	// Store state in a cookie for verification in the callback
	c.SetCookie(
		"oauth_state",
		state,
		600, // 10 minutes
		"/",
		h.cookieDomain,
		false,
		true,
	)

	// Redirect to Discord authorization URL
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}
func (h *UserHandler) DiscordCallback(c *gin.Context) {
	// Get the state and code from the query parameters
	queryState := c.Query("state")
	code := c.Query("code")

	// Verify the state to prevent CSRF attacks
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || cookieState != queryState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	// Clear the oauth_state cookie
	c.SetCookie("oauth_state", "", -1, "/", h.cookieDomain, false, true)

	// Exchange the authorization code for a token
	token, err := h.oauthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Use the token to fetch the user's information from Discord
	discordUser, err := h.getDiscordUser(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info from Discord"})
		return
	}

	// Create or update user in the database
	user, err := h.App.UserService.UserFromDiscord(discordUser.ID, discordUser.Username, discordUser.Email, discordUser.Avatar, discordUser.Discriminator)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user"})
		return
	}
	// Generate JWT tokens
	accessToken, err := h.App.AuthService.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}
	refreshToken, err := h.App.AuthService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Set the tokens as cookies
	setAccessTokenCookie(c, accessToken, h.cookieDomain)
	setRefreshTokenCookie(c, refreshToken, h.cookieDomain)

	// Redirect to the frontend
	//c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/success")
	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL)
}
func (h *UserHandler) getDiscordUser(accessToken string) (*DiscordUser, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("discord API returned non-200 status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var discordUser DiscordUser
	if err := json.Unmarshal(body, &discordUser); err != nil {
		return nil, err
	}

	return &discordUser, nil
}

// REST

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
func (h *UserHandler) HelloAuthorized(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello Authorized"})
}
func (H *UserHandler) GetUser(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found in context"})
		log.Println("USER NOT FOUND IN CONTEXT")
		return
	}
	user, ok := userInterface.(models.User) // Replace with your actual User model
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to cast user"})
		log.Println("FAILED TO CAST")
		return
	}
	userDTO := user.ToUserDTO()
	c.JSON(http.StatusOK, userDTO)
}
func (h *UserHandler) RefreshTokens(c *gin.Context) {
	refreshToken, err := getRefreshTokenFromCookie(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No token"})
		return
	}
	user, err := h.App.UserService.GetUserFromRefreshToken(refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}
	newRefreshToken, err := h.App.AuthService.GenerateAccessToken(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}
	newAccessToken, err := h.App.AuthService.GenerateAccessToken(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}
	setAccessTokenCookie(c, newAccessToken, h.cookieDomain)
	setRefreshTokenCookie(c, newRefreshToken, h.cookieDomain)
	c.JSON(http.StatusOK, gin.H{"message": "Refreshed tokens"})
}
