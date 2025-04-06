package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
	"vibely-backend/src/models"
)

func GetDB(user, password, name string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Warsaw",
		user,
		password,
		name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.New("failed to open gorm connection")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to configure database connection: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Add all models for migration
	err = db.AutoMigrate(
		&models.User{},
		&models.Lesson{},
		&models.Course{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate: %w", err)
	}

	log.Println("Database migration completed successfully.")
	return db, nil
}
