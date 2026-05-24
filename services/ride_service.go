package services

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/ranaganakulathunga/taxi-booking-backend-go/models"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/repositories"
)

type RideService interface {
	RequestRide(passengerID uint, pickupLat, pickupLng, dropoffLat, dropoffLng float64, address string) (*models.Ride, error)
	GetRideDetails(rideID uint) (*models.Ride, error)
	CancelRide(rideID uint) error
	StartRide(rideID uint) error
	CompleteRide(rideID uint, actualDistance float64, actualFare float64) error
	GetPassengerRideHistory(passengerID uint) ([]models.Ride, error)
	GetDriverRideHistory(driverID uint) ([]models.Ride, error)
	GetPendingRides() ([]models.Ride, error)
	AllocateDriver(rideID uint, driverID uint) (*models.Ride, error)
}

type rideService struct {
	rideRepo   repositories.RideRepository
	driverRepo repositories.DriverRepository
}

func NewRideService(rideRepo repositories.RideRepository, driverRepo repositories.DriverRepository) RideService {
	return &rideService{
		rideRepo:   rideRepo,
		driverRepo: driverRepo,
	}
}

// RequestRide creates a new ride request
func (s *rideService) RequestRide(passengerID uint, pickupLat, pickupLng, dropoffLat, dropoffLng float64, address string) (*models.Ride, error) {
	// Calculate distance using Haversine formula
	distance := s.haversineDistance(pickupLat, pickupLng, dropoffLat, dropoffLng)

	// Create ride
	ride := &models.Ride{
		RideCode:           s.generateRideCode(),
		PassengerID:        &passengerID,
		PickupLatitude:     pickupLat,
		PickupLongitude:    pickupLng,
		DropoffLatitude:    dropoffLat,
		DropoffLongitude:   dropoffLng,
		DropoffAddress:     address,
		Status:             models.RideRequested,
		RequestedAt:        time.Now(),
		EstimatedDistance:  distance,
		NumberOfPassengers: 1,
	}

	// Calculate estimate
	ride.CalculateEstimate()

	// Save to database
	err := s.rideRepo.Create(ride)
	if err != nil {
		return nil, err
	}

	return ride, nil
}

// GetRideDetails fetches ride details
func (s *rideService) GetRideDetails(rideID uint) (*models.Ride, error) {
	return s.rideRepo.GetByID(rideID)
}

// CancelRide cancels an active ride
func (s *rideService) CancelRide(rideID uint) error {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil {
		return err
	}

	// Can only cancel requested or assigned rides
	if ride.Status != models.RideRequested && ride.Status != models.RideAssigned {
		return fmt.Errorf("cannot cancel ride with status %s", ride.Status)
	}

	now := time.Now()
	ride.Status = models.RideCancelled
	ride.CancelledAt = &now

	return s.rideRepo.Update(ride)
}

// StartRide marks a ride as started
func (s *rideService) StartRide(rideID uint) error {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil {
		return err
	}

	if ride.Status != models.RideAssigned {
		return fmt.Errorf("ride must be assigned before starting")
	}

	now := time.Now()
	ride.Status = models.RideStarted
	ride.StartedAt = &now

	return s.rideRepo.Update(ride)
}

// CompleteRide marks a ride as completed
func (s *rideService) CompleteRide(rideID uint, actualDistance float64, actualFare float64) error {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil {
		return err
	}

	if ride.Status != models.RideStarted {
		return fmt.Errorf("ride must be started before completing")
	}

	now := time.Now()
	ride.Status = models.RideCompleted
	ride.CompletedAt = &now
	ride.ActualDistance = actualDistance
	ride.ActualFare = actualFare

	// Update driver statistics
	if ride.DriverID != nil {
		driver, err := s.driverRepo.GetByID(*ride.DriverID)
		if err == nil {
			driver.TotalRides++
			driver.TotalEarnings += actualFare
			s.driverRepo.Update(driver)
		}
	}

	return s.rideRepo.Update(ride)
}

// GetPassengerRideHistory gets all rides for a passenger
func (s *rideService) GetPassengerRideHistory(passengerID uint) ([]models.Ride, error) {
	return s.rideRepo.GetPassengerRides(passengerID)
}

// GetDriverRideHistory gets all rides for a driver
func (s *rideService) GetDriverRideHistory(driverID uint) ([]models.Ride, error) {
	return s.rideRepo.GetDriverRides(driverID)
}

// GetPendingRides gets all pending rides
func (s *rideService) GetPendingRides() ([]models.Ride, error) {
	return s.rideRepo.GetPendingRides()
}

// AllocateDriver assigns a driver to a ride
func (s *rideService) AllocateDriver(rideID uint, driverID uint) (*models.Ride, error) {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	if ride.Status != models.RideRequested {
		return nil, fmt.Errorf("ride must be in requested status to allocate driver")
	}

	// Verify driver exists and is available
	driver, err := s.driverRepo.GetByID(driverID)
	if err != nil {
		return nil, err
	}

	if driver.Status != models.DriverOnline {
		return nil, fmt.Errorf("driver is not online")
	}

	if !driver.IsVerified {
		return nil, fmt.Errorf("driver is not verified")
	}

	// Update ride
	now := time.Now()
	ride.DriverID = &driverID
	ride.Status = models.RideAssigned
	ride.AssignedAt = &now

	err = s.rideRepo.Update(ride)
	if err != nil {
		return nil, err
	}

	// Update driver status
	driver.Status = models.DriverOnTrip
	s.driverRepo.Update(driver)

	return ride, nil
}

// FindNearestDrivers finds nearest available drivers to a pickup location
func (s *rideService) FindNearestDrivers(pickupLat, pickupLng float64, limit int) ([]models.Driver, error) {
	drivers, err := s.driverRepo.GetAvailableDrivers()
	if err != nil {
		return nil, err
	}

	// Calculate distances and sort
	type driverDistance struct {
		driver   models.Driver
		distance float64
	}

	var driverDistances []driverDistance
	for _, driver := range drivers {
		distance := s.haversineDistance(pickupLat, pickupLng, driver.CurrentLatitude, driver.CurrentLongitude)
		driverDistances = append(driverDistances, driverDistance{
			driver:   driver,
			distance: distance,
		})
	}

	// Sort by distance
	for i := 0; i < len(driverDistances); i++ {
		for j := i + 1; j < len(driverDistances); j++ {
			if driverDistances[j].distance < driverDistances[i].distance {
				driverDistances[i], driverDistances[j] = driverDistances[j], driverDistances[i]
			}
		}
	}

	// Return top N drivers
	result := []models.Driver{}
	for i := 0; i < len(driverDistances) && i < limit; i++ {
		result = append(result, driverDistances[i].driver)
	}

	return result, nil
}

// Helper function to calculate distance between two coordinates (Haversine formula)
func (s *rideService) haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadiusKm * c

	return math.Round(distance*100) / 100 // Round to 2 decimal places
}

// Helper function to generate unique ride code
func (s *rideService) generateRideCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "RIDE" + string(b)
}
