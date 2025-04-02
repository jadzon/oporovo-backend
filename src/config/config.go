package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port                string
	CookieDomain        string
	PostgresUser        string
	PostgresPassword    string
	PostgresDB          string
	AccessJWTSecretKey  string
	RefreshJWTSecretKey string
	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURL  string
	FrontendUrl         string
	BackendUrl          string
}

func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	return Config{
		Port:                getEnv("PORT", "-"),
		CookieDomain:        getEnv("COOKIE_DOMAIN", "-"),
		PostgresUser:        getEnv("POSTGRES_USER", "-"),
		PostgresPassword:    getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:          getEnv("POSTGRES_DB", "-"),
		AccessJWTSecretKey:  getEnv("ACCESS_JWT_SECRET_KEY", ""),
		RefreshJWTSecretKey: getEnv("REFRESH_JWT_SECRET_KEY", ""),
		DiscordClientID:     getEnv("DISCORD_CLIENT_ID", ""),
		DiscordClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
		DiscordRedirectURL:  getEnv("DISCORD_REDIRECT_URL", ""),
		FrontendUrl:         getEnv("FRONTEND_URL", ""),
		BackendUrl:          getEnv("BACKEND_URL", ""),
	}
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
