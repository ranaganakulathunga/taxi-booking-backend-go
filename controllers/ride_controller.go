package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/services"
)

type RideController struct {
	rideService   services.RideService
	driverService services.DriverService
}

func NewRideController(rideService services.RideService, driverService services.DriverService) *RideController {
	return &RideController{
		rideService:   rideService,
		driverService: driverService,
	}
}

// RequestRideRequest is the request body for requesting a ride
type RequestRideRequest struct {
	PassengerID      uint    `json:"passenger_id" binding:"required"`
	PickupLatitude   float64 `json:"pickup_latitude" binding:"required"`
	PickupLongitude  float64 `json:"pickup_longitude" binding:"required"`
	DropoffLatitude  float64 `json:"dropoff_latitude" binding:"required"`
	DropoffLongitude float64 `json:"dropoff_longitude" binding:"required"`
	DropoffAddress   string  `json:"dropoff_address" binding:"required"`
}

// RequestRide handles POST /api/v1/rides
func (rc *RideController) RequestRide(c *gin.Context) {
	var req RequestRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ride, err := rc.rideService.RequestRide(
		req.PassengerID,
		req.PickupLatitude,
		req.PickupLongitude,
		req.DropoffLatitude,
		req.DropoffLongitude,
		req.DropoffAddress,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ride requested successfully",
		"ride":    ride,
	})
}

// GetRideDetails handles GET /api/v1/rides/:id
func (rc *RideController) GetRideDetails(c *gin.Context) {
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	ride, err := rc.rideService.GetRideDetails(uint(rideID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ride)
}

// CancelRideRequest is the request for cancelling a ride
type CancelRideRequest struct {
	Reason string `json:"reason,omitempty"`
}

// CancelRide handles PUT /api/v1/rides/:id/cancel
func (rc *RideController) CancelRide(c *gin.Context) {
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	var req CancelRideRequest
	c.ShouldBindJSON(&req)

	err = rc.rideService.CancelRide(uint(rideID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ride cancelled successfully",
	})
}

// StartRideRequest is the request for starting a ride
type StartRideRequest struct {
	DriverID uint `json:"driver_id" binding:"required"`
}

// StartRide handles PUT /api/v1/rides/:id/start
func (rc *RideController) StartRide(c *gin.Context) {
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	var req StartRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = rc.rideService.StartRide(uint(rideID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ride started successfully",
	})
}

// CompleteRideRequest is the request for completing a ride
type CompleteRideRequest struct {
	ActualDistance float64 `json:"actual_distance" binding:"required"`
	ActualFare     float64 `json:"actual_fare" binding:"required"`
}

// CompleteRide handles PUT /api/v1/rides/:id/complete
func (rc *RideController) CompleteRide(c *gin.Context) {
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	var req CompleteRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = rc.rideService.CompleteRide(uint(rideID), req.ActualDistance, req.ActualFare)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ride completed successfully",
	})
}

// GetPassengerRides handles GET /api/v1/rides/passenger/:id
func (rc *RideController) GetPassengerRides(c *gin.Context) {
	passengerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passenger ID"})
		return
	}

	rides, err := rc.rideService.GetPassengerRideHistory(uint(passengerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rides": rides,
	})
}

// GetPendingRides handles GET /api/v1/rides/pending
func (rc *RideController) GetPendingRides(c *gin.Context) {
	rides, err := rc.rideService.GetPendingRides()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rides": rides,
	})
}

// AllocateDriverRequest is the request for allocating a driver
type AllocateDriverRequest struct {
	DriverID uint `json:"driver_id" binding:"required"`
}

// AllocateDriver handles POST /api/v1/rides/:id/allocate
func (rc *RideController) AllocateDriver(c *gin.Context) {
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	var req AllocateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ride, err := rc.rideService.AllocateDriver(uint(rideID), req.DriverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Driver allocated successfully",
		"ride":    ride,
	})
}

// FindNearestDriversRequest is the request for finding nearest drivers
type FindNearestDriversRequest struct {
	PickupLatitude  float64 `json:"pickup_latitude" binding:"required"`
	PickupLongitude float64 `json:"pickup_longitude" binding:"required"`
	Limit           int     `json:"limit" binding:"required"`
}

// FindNearestDrivers handles POST /api/v1/rides/drivers/nearest
func (rc *RideController) FindNearestDrivers(c *gin.Context) {
	var req FindNearestDriversRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// drivers, err := rc.rideService.FindNearestDrivers(req.PickupLatitude, req.PickupLongitude, req.Limit)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"drivers": drivers,
	// }
	//)
}
