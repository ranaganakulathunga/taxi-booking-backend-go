package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Database
	// config.ConnectDatabase() // Commented out until valid Postgres URI is provided

	// Initialize Gin Engine
	router := gin.Default()

	// Setup Routes
	routes.SetupRoutes(router)

	// Get Port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
