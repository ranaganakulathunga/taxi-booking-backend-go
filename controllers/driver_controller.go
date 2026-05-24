package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/models"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/services"
)

type DriverController struct {
	driverService services.DriverService
}

func NewDriverController(driverService services.DriverService) *DriverController {
	return &DriverController{
		driverService: driverService,
	}
}

// GetDriverDetails handles GET /api/v1/drivers/:id
func (dc *DriverController) GetDriverDetails(c *gin.Context) {
	driverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	driver, err := dc.driverService.GetDriverByID(uint(driverID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, driver)
}

// UpdateLocationRequest is the request body for updating driver location
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// UpdateLocation handles PUT /api/v1/drivers/:id/location
func (dc *DriverController) UpdateLocation(c *gin.Context) {
	driverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = dc.driverService.UpdateDriverLocation(uint(driverID), req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Location updated successfully",
	})
}

// UpdateStatusRequest is the request body for updating driver status
type UpdateStatusRequest struct {
	Status models.DriverStatus `json:"status" binding:"required"`
}

// UpdateStatus handles PUT /api/v1/drivers/:id/status
func (dc *DriverController) UpdateStatus(c *gin.Context) {
	driverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = dc.driverService.UpdateDriverStatus(uint(driverID), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status updated successfully",
	})
}

// GetAvailableDrivers handles GET /api/v1/drivers/available
func (dc *DriverController) GetAvailableDrivers(c *gin.Context) {
	drivers, err := dc.driverService.GetAvailableDrivers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"drivers": drivers,
	})
}

// GetOnlineDrivers handles GET /api/v1/drivers/online
func (dc *DriverController) GetOnlineDrivers(c *gin.Context) {
	drivers, err := dc.driverService.GetOnlineDrivers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"drivers": drivers,
	})
}

// GetDriverRides handles GET /api/v1/drivers/:id/rides
func (dc *DriverController) GetDriverRides(c *gin.Context) {
	driverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	rides, err := dc.driverService.GetDriverRides(uint(driverID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rides": rides,
	})
}
