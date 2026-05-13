package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ranaganakulathunga/taxi-booking-backend-go/controllers"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.GET("/health", controllers.HealthCheck)
	}
}
