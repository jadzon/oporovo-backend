package app

import (
	"gorm.io/gorm"
	"vibely-backend/src/config"
	"vibely-backend/src/database"
	"vibely-backend/src/repositories"
	"vibely-backend/src/services"
)

type Application struct {
	Config        config.Config
	DB            *gorm.DB
	UserService   services.UserService
	AuthService   services.AuthService
	LessonService services.LessonService
	CourseService services.CourseService
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
	lessonRepository := repositories.NewLessonRepository(db)
	lessonService := services.NewLessonService(lessonRepository)
	courseRepository := repositories.NewCourseRepository(db)
	courseService := services.NewCourseService(courseRepository)
	return &Application{
		Config:        cfg,
		DB:            db,
		UserService:   userService,
		AuthService:   authService,
		LessonService: lessonService,
		CourseService: courseService,
	}, nil
}
