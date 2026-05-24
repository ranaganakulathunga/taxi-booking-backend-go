package repositories

import (
	"errors"

	"github.com/ranaganakulathunga/taxi-booking-backend-go/models"
	"gorm.io/gorm"
)

type RideRepository interface {
	Create(ride *models.Ride) error
	GetByID(id uint) (*models.Ride, error)
	GetByRideCode(code string) (*models.Ride, error)
	Update(ride *models.Ride) error
	Delete(id uint) error
	GetPassengerRides(passengerID uint) ([]models.Ride, error)
	GetDriverRides(driverID uint) ([]models.Ride, error)
	GetRidesByStatus(status models.RideStatus) ([]models.Ride, error)
	GetPendingRides() ([]models.Ride, error)
}

type rideRepository struct {
	db *gorm.DB
}

func NewRideRepository(db *gorm.DB) RideRepository {
	return &rideRepository{db: db}
}

func (r *rideRepository) Create(ride *models.Ride) error {
	return r.db.Create(ride).Error
}

func (r *rideRepository) GetByID(id uint) (*models.Ride, error) {
	var ride models.Ride
	err := r.db.Preload("Passenger").Preload("Driver").First(&ride, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("ride not found")
	}
	return &ride, err
}

func (r *rideRepository) GetByRideCode(code string) (*models.Ride, error) {
	var ride models.Ride
	err := r.db.Where("ride_code = ?", code).
		Preload("Passenger").
		Preload("Driver").
		First(&ride).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("ride not found")
	}
	return &ride, err
}

func (r *rideRepository) Update(ride *models.Ride) error {
	return r.db.Save(ride).Error
}

func (r *rideRepository) Delete(id uint) error {
	return r.db.Delete(&models.Ride{}, id).Error
}

func (r *rideRepository) GetPassengerRides(passengerID uint) ([]models.Ride, error) {
	var rides []models.Ride
	err := r.db.Where("passenger_id = ?", passengerID).
		Preload("Driver").
		Order("created_at DESC").
		Find(&rides).Error
	return rides, err
}

func (r *rideRepository) GetDriverRides(driverID uint) ([]models.Ride, error) {
	var rides []models.Ride
	err := r.db.Where("driver_id = ?", driverID).
		Preload("Passenger").
		Order("created_at DESC").
		Find(&rides).Error
	return rides, err
}

func (r *rideRepository) GetRidesByStatus(status models.RideStatus) ([]models.Ride, error) {
	var rides []models.Ride
	err := r.db.Where("status = ?", status).
		Preload("Passenger").
		Preload("Driver").
		Order("created_at DESC").
		Find(&rides).Error
	return rides, err
}

func (r *rideRepository) GetPendingRides() ([]models.Ride, error) {
	var rides []models.Ride
	err := r.db.Where("status IN ?", []models.RideStatus{models.RideRequested, models.RideAssigned}).
		Preload("Passenger").
		Preload("Driver").
		Order("created_at ASC").
		Find(&rides).Error
	return rides, err
}
