//

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	if dsn == "" {
		log.Fatal("No database connection string provided. Set DATABASE_URL or DB_HOST/DB_USER/DB_PASSWORD/DB_NAME/DB_PORT")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	DB = database

	log.Println("Database connected successfully")

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.Driver{},
		&models.Vehicle{},
		&models.Location{},
		&models.Ride{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database")
	}

	log.Println("Database migrated successfully")
}
