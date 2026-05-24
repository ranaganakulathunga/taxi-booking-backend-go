package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Driver represents a taxi driver in the system
type Driver struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `gorm:"uniqueIndex" json:"email"`
	PhoneNumber string `gorm:"uniqueIndex" json:"phone_number"`
	Password    string `json:"-"` // Hashed password, never exposed in JSON

	ProfilePicture string  `json:"profile_picture,omitempty"`
	Rating         float64 `gorm:"default:5.0" json:"rating"`

	// License and verification
	LicenseNumber string    `gorm:"uniqueIndex" json:"license_number"`
	LicenseExpiry time.Time `json:"license_expiry"`
	IsVerified    bool      `gorm:"default:false" json:"is_verified"`

	// Driver status
	Status    DriverStatus `gorm:"type:varchar(50);default:'offline'" json:"status"` // offline, online, on_trip
	IsActive  bool         `gorm:"default:true" json:"is_active"`
	IsBlocked bool         `gorm:"default:false" json:"is_blocked"`

	// Current location (last known)
	CurrentLatitude  float64   `json:"current_latitude"`
	CurrentLongitude float64   `json:"current_longitude"`
	LastLocationTime time.Time `json:"last_location_time,omitempty"`

	// Statistics
	TotalRides     int     `gorm:"default:0" json:"total_rides"`
	TotalEarnings  float64 `gorm:"default:0" json:"total_earnings"`
	AcceptanceRate float64 `gorm:"default:0" json:"acceptance_rate"`

	// Relationships
	VehicleID uint    `json:"vehicle_id,omitempty"`
	Vehicle   Vehicle `gorm:"foreignKey:VehicleID;constraint:OnDelete:SET NULL" json:"vehicle,omitempty"`
	Rides     []Ride  `gorm:"foreignKey:DriverID;constraint:OnDelete:SET NULL" json:"rides,omitempty"`

	// JSON metadata for flexible storage
	Metadata datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// DriverStatus defines the status of a driver
type DriverStatus string

const (
	DriverOffline DriverStatus = "offline"
	DriverOnline  DriverStatus = "online"
	DriverOnTrip  DriverStatus = "on_trip"
)

func (d *Driver) BeforeSave(tx *gorm.DB) error {
	return nil
}
