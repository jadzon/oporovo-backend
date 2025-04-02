package app

import (
	"gorm.io/gorm"
	"vibely-backend/src/config"
	"vibely-backend/src/database"
	"vibely-backend/src/repositories"
	"vibely-backend/src/services"
)

type Application struct {
	Config      config.Config
	DB          *gorm.DB
	UserService services.UserService
	AuthService services.AuthService
}

func New(cfg config.Config) (*Application, error) {
	// Initialize the database using config
	db, err := database.GetDB(cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	if err != nil {
		return nil, err
	}
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(cfg.AccessJWTSecretKey, cfg.RefreshJWTSecretKey)
	userService := services.NewUserService(userRepository, authService)
	return &Application{
		Config:      cfg,
		DB:          db,
		UserService: userService,
		AuthService: authService,
	}, nil
}
