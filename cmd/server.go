package main

import (
	"log"
	"os"
	"vibely-backend/src/app"
	"vibely-backend/src/config"
	"vibely-backend/src/routes"
)

func main() {
	// Load configurations
	cfg := config.NewConfig()
	wd, err := os.Getwd() // Get current working directory
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	log.Printf("Current working directory: %s", wd)
	log.Printf("Config Loaded: USER=%s, PASSWORD=%s, DB=%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	// Pass config to application initialization
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
		return
	}
	//seed
	//if err := seeds.Seed(application.DB); err != nil {
	//	log.Fatalf("Seeding error: %v", err)
	//}
	router := routes.Setup(application)
	port := ":" + cfg.Port
	log.Printf("Starting server on port %s...", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
