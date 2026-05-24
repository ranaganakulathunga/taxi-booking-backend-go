package models

import (
	"time"
)

// Vehicle represents a taxi vehicle
type Vehicle struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Basic details
	Make  string `json:"make"`  // e.g., "Toyota"
	Model string `json:"model"` // e.g., "Camry"
	Year  int    `json:"year"`

	// Identification
	RegistrationNumber string `gorm:"uniqueIndex" json:"registration_number"`
	VINNumber          string `gorm:"uniqueIndex" json:"vin_number"`
	Color              string `json:"color"`

	// Capacity
	MaxPassengers int `gorm:"default:4" json:"max_passengers"`

	// Status
	IsActive bool `gorm:"default:true" json:"is_active"`

	// Relationships
	Drivers []Driver `gorm:"foreignKey:VehicleID;constraint:OnDelete:SET NULL" json:"drivers,omitempty"`
}
