package repositories

import (
	"errors"

	"github.com/ranaganakulathunga/taxi-booking-backend-go/models"
	"gorm.io/gorm"
)

type DriverRepository interface {
	Create(driver *models.Driver) error
	GetByID(id uint) (*models.Driver, error)
	GetByEmail(email string) (*models.Driver, error)
	Update(driver *models.Driver) error
	Delete(id uint) error
	GetAllActive() ([]models.Driver, error)
	GetOnlineDrivers() ([]models.Driver, error)
	GetAvailableDrivers() ([]models.Driver, error)
}

type driverRepository struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) Create(driver *models.Driver) error {
	return r.db.Create(driver).Error
}

func (r *driverRepository) GetByID(id uint) (*models.Driver, error) {
	var driver models.Driver
	err := r.db.Preload("Vehicle").First(&driver, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("driver not found")
	}
	return &driver, err
}

func (r *driverRepository) GetByEmail(email string) (*models.Driver, error) {
	var driver models.Driver
	err := r.db.Where("email = ?", email).Preload("Vehicle").First(&driver).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("driver not found")
	}
	return &driver, err
}

func (r *driverRepository) Update(driver *models.Driver) error {
	return r.db.Save(driver).Error
}

func (r *driverRepository) Delete(id uint) error {
	return r.db.Delete(&models.Driver{}, id).Error
}

func (r *driverRepository) GetAllActive() ([]models.Driver, error) {
	var drivers []models.Driver
	err := r.db.Where("is_active = ? AND is_blocked = ?", true, false).
		Preload("Vehicle").
		Order("rating DESC").
		Find(&drivers).Error
	return drivers, err
}

func (r *driverRepository) GetOnlineDrivers() ([]models.Driver, error) {
	var drivers []models.Driver
	err := r.db.Where("status = ? AND is_active = ? AND is_blocked = ?", models.DriverOnline, true, false).
		Preload("Vehicle").
		Order("rating DESC").
		Find(&drivers).Error
	return drivers, err
}

func (r *driverRepository) GetAvailableDrivers() ([]models.Driver, error) {
	var drivers []models.Driver
	err := r.db.Where("status = ? AND is_active = ? AND is_blocked = ? AND is_verified = ?",
		models.DriverOnline, true, false, true).
		Preload("Vehicle").
		Order("rating DESC").
		Find(&drivers).Error
	return drivers, err
}
