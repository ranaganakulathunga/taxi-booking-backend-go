package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a passenger/customer in the system
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `gorm:"uniqueIndex" json:"email"`
	PhoneNumber string `gorm:"uniqueIndex" json:"phone_number"`
	Password    string `json:"-"` // Hashed password, never exposed in JSON

	ProfilePicture string  `json:"profile_picture,omitempty"`
	Rating         float64 `gorm:"default:0" json:"rating"`

	// Address fields for default location
	Address string `json:"address,omitempty"`
	City    string `json:"city,omitempty"`

	// Relationships
	Rides []Ride `gorm:"foreignKey:PassengerID;constraint:OnDelete:CASCADE" json:"rides,omitempty"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}

// BeforeSave hooks
func (u *User) BeforeSave(tx *gorm.DB) error {
	return nil
}
