package services

import (
	"fmt"
	"time"

	"github.com/ranaganakulathunga/taxi-booking-backend-go/models"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/repositories"
)

type DriverService interface {
	GetDriverByID(driverID uint) (*models.Driver, error)
	UpdateDriverLocation(driverID uint, latitude, longitude float64) error
	UpdateDriverStatus(driverID uint, status models.DriverStatus) error
	GetAvailableDrivers() ([]models.Driver, error)
	GetOnlineDrivers() ([]models.Driver, error)
	GetDriverRides(driverID uint) ([]models.Ride, error)
}

type driverService struct {
	driverRepo repositories.DriverRepository
	rideRepo   repositories.RideRepository
}

func NewDriverService(driverRepo repositories.DriverRepository, rideRepo repositories.RideRepository) DriverService {
	return &driverService{
		driverRepo: driverRepo,
		rideRepo:   rideRepo,
	}
}

// GetDriverByID retrieves driver details
func (s *driverService) GetDriverByID(driverID uint) (*models.Driver, error) {
	return s.driverRepo.GetByID(driverID)
}

// UpdateDriverLocation updates driver's current location
func (s *driverService) UpdateDriverLocation(driverID uint, latitude, longitude float64) error {
	driver, err := s.driverRepo.GetByID(driverID)
	if err != nil {
		return err
	}

	driver.CurrentLatitude = latitude
	driver.CurrentLongitude = longitude
	driver.LastLocationTime = time.Now()

	return s.driverRepo.Update(driver)
}

// UpdateDriverStatus updates driver's status
func (s *driverService) UpdateDriverStatus(driverID uint, status models.DriverStatus) error {
	driver, err := s.driverRepo.GetByID(driverID)
	if err != nil {
		return err
	}

	// Validate status transitions
	validTransitions := map[models.DriverStatus][]models.DriverStatus{
		models.DriverOffline: {models.DriverOnline},
		models.DriverOnline:  {models.DriverOffline, models.DriverOnTrip},
		models.DriverOnTrip:  {models.DriverOnline},
	}

	if validStatuses, ok := validTransitions[driver.Status]; ok {
		isValid := false
		for _, validStatus := range validStatuses {
			if validStatus == status {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid status transition from %s to %s", driver.Status, status)
		}
	}

	driver.Status = status
	return s.driverRepo.Update(driver)
}

// GetAvailableDrivers gets all available drivers
func (s *driverService) GetAvailableDrivers() ([]models.Driver, error) {
	return s.driverRepo.GetAvailableDrivers()
}

// GetOnlineDrivers gets all online drivers
func (s *driverService) GetOnlineDrivers() ([]models.Driver, error) {
	return s.driverRepo.GetOnlineDrivers()
}

// GetDriverRides gets all rides for a driver
func (s *driverService) GetDriverRides(driverID uint) ([]models.Ride, error) {
	return s.rideRepo.GetDriverRides(driverID)
}
