package models

import (
	"time"

	"gorm.io/datatypes"
)

// Ride represents a trip/ride request
type Ride struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Ride identifiers
	RideCode string `gorm:"uniqueIndex;type:varchar(50)" json:"ride_code"`

	// Passenger and Driver
	PassengerID *uint   `json:"passenger_id"`
	Passenger   *User   `gorm:"foreignKey:PassengerID;constraint:OnDelete:CASCADE" json:"passenger,omitempty"`
	DriverID    *uint   `json:"driver_id"`
	Driver      *Driver `gorm:"foreignKey:DriverID;constraint:OnDelete:SET NULL" json:"driver,omitempty"`

	// Pickup location
	PickupLatitude  float64 `json:"pickup_latitude"`
	PickupLongitude float64 `json:"pickup_longitude"`
	PickupAddress   string  `json:"pickup_address"`

	// Dropoff location
	DropoffLatitude  float64 `json:"dropoff_latitude"`
	DropoffLongitude float64 `json:"dropoff_longitude"`
	DropoffAddress   string  `json:"dropoff_address"`

	// Ride status
	Status RideStatus `gorm:"type:varchar(50);default:'requested'" json:"status"` // requested, assigned, started, completed, cancelled

	// Estimated and actual times
	RequestedAt       time.Time  `json:"requested_at"`
	AssignedAt        *time.Time `json:"assigned_at,omitempty"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	CancelledAt       *time.Time `json:"cancelled_at,omitempty"`
	EstimatedDuration int        `json:"estimated_duration"` // in seconds

	// Pricing
	EstimatedFare float64 `json:"estimated_fare"`
	ActualFare    float64 `json:"actual_fare,omitempty"`
	Tip           float64 `gorm:"default:0" json:"tip"`
	DiscountCode  string  `json:"discount_code,omitempty"`

	// Distance
	EstimatedDistance float64 `json:"estimated_distance"` // in km
	ActualDistance    float64 `json:"actual_distance,omitempty"`

	// Ride details
	NumberOfPassengers int    `gorm:"default:1" json:"number_of_passengers"`
	SpecialRequests    string `json:"special_requests,omitempty"`

	// Ratings
	DriverRating    int `json:"driver_rating,omitempty"`    // 1-5
	PassengerRating int `json:"passenger_rating,omitempty"` // 1-5

	// JSON metadata for flexible storage
	Metadata datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// RideStatus defines the status of a ride
type RideStatus string

const (
	RideRequested RideStatus = "requested"
	RideAssigned  RideStatus = "assigned"
	RideStarted   RideStatus = "started"
	RideCompleted RideStatus = "completed"
	RideCancelled RideStatus = "cancelled"
)

// CalculateEstimate calculates the estimated fare and duration based on distance
func (r *Ride) CalculateEstimate() {
	// Base fare in cents
	const baseFare float64 = 250  // $2.50
	const perKmRate float64 = 100 // $1.00 per km
	const perMinRate float64 = 35 // $0.35 per minute

	// Estimate fare (simplified)
	r.EstimatedFare = baseFare + (r.EstimatedDistance * perKmRate)

	// Estimate duration (simplified) - average 30 km/h = 2 min per km
	r.EstimatedDuration = int(r.EstimatedDistance * 120)
}
