package main

import (
	"log"

	"diploma/internal/api/routes"
	"diploma/internal/config"

	_ "diploma/docs" // Import the generated Swagger docs
)

// @title DiplomaAPI
// @version 1.0
// @description This is a REST API for managing medical organization app.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Gin router
	router := routes.SetupRouter()

	// Start the server
	log.Printf("Server starting on port %s...\n", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
