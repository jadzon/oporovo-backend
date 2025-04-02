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
	}
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
