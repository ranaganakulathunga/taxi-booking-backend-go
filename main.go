package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/config"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/models"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/realtime"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Database
	config.ConnectDatabase()
	log.Println("Running database migrations...")
	migrateDatabase()

	// Initialize Gin Engine
	router := gin.Default()

	// Initialize WebSocket Hub
	wsHub := realtime.NewHub()
	go wsHub.Run()

	// Setup Routes
	routes.SetupRoutes(router)

	// WebSocket endpoint
	router.GET("/api/v1/ws/ride-updates", realtime.WebSocketHandler(wsHub))

	// Get Port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("WebSocket server available at ws://localhost:%s/api/v1/ws/ride-updates", port)
	router.Run(":" + port)
}

func migrateDatabase() {
	migrations := []interface{}{
		&models.User{},
		&models.Driver{},
		&models.Vehicle{},
		&models.Ride{},
		&models.Location{},
	}

	for _, migration := range migrations {
		if err := config.DB.AutoMigrate(migration); err != nil {
			log.Fatalf("Failed to migrate %T: %v", migration, err)
		}
	}

	log.Println("Database migrations completed successfully")
}
