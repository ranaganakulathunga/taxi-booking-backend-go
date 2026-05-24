package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/config"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/controllers"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/repositories"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/services"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", controllers.HealthCheck)

		// Initialize repositories
		rideRepo := repositories.NewRideRepository(config.DB)
		driverRepo := repositories.NewDriverRepository(config.DB)

		// Initialize services
		rideService := services.NewRideService(rideRepo, driverRepo)
		driverService := services.NewDriverService(driverRepo, rideRepo)

		// Initialize controllers
		rideController := controllers.NewRideController(rideService, driverService)
		driverController := controllers.NewDriverController(driverService)

		// Ride endpoints
		rides := api.Group("/rides")
		{
			rides.POST("", rideController.RequestRide)                        // POST /api/v1/rides - request a ride
			rides.GET("/:id", rideController.GetRideDetails)                  // GET /api/v1/rides/:id - ride details
			rides.GET("/passenger/:id", rideController.GetPassengerRides)     // GET /api/v1/rides/passenger/:id - passenger rides
			rides.GET("/pending", rideController.GetPendingRides)             // GET /api/v1/rides/pending - pending rides
			rides.PUT("/:id/cancel", rideController.CancelRide)               // PUT /api/v1/rides/:id/cancel - cancel ride
			rides.PUT("/:id/start", rideController.StartRide)                 // PUT /api/v1/rides/:id/start - start ride
			rides.PUT("/:id/complete", rideController.CompleteRide)           // PUT /api/v1/rides/:id/complete - complete ride
			rides.POST("/:id/allocate", rideController.AllocateDriver)        // POST /api/v1/rides/:id/allocate - allocate driver
			rides.POST("/drivers/nearest", rideController.FindNearestDrivers) // POST /api/v1/rides/drivers/nearest - find nearest drivers
		}

		// Driver endpoints
		drivers := api.Group("/drivers")
		{
			drivers.GET("/:id", driverController.GetDriverDetails)          // GET /api/v1/drivers/:id - driver details
			drivers.GET("/available", driverController.GetAvailableDrivers) // GET /api/v1/drivers/available - available drivers
			drivers.GET("/online", driverController.GetOnlineDrivers)       // GET /api/v1/drivers/online - online drivers
			drivers.GET("/:id/rides", driverController.GetDriverRides)      // GET /api/v1/drivers/:id/rides - driver rides
			drivers.PUT("/:id/location", driverController.UpdateLocation)   // PUT /api/v1/drivers/:id/location - update location
			drivers.PUT("/:id/status", driverController.UpdateStatus)       // PUT /api/v1/drivers/:id/status - update status
		}

		// WebSocket endpoints (placeholder for now)
		// ws := api.Group("/ws")
		// {
		//     ws.GET("/ride-updates", websocketHandler)
		// }
	}
}
