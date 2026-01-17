package rest

import (
	"github.com/gin-gonic/gin"
)

func (h HttpHandler) SetupRouter() *gin.Engine {
	// Set Gin mode based on configuration
	// gin.SetMode(h.config.GinMode)

	router := gin.Default()

	// router.Use(middleware.ErrorHandlerMiddleware())

	// Swagger documentation
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health and Readiness endpoints
	// router.GET("/health", h.HealthCheck)
	// router.GET("/readiness", h.ReadinessCheck)

	// Weather endpoint
	v1 := router.Group("/api/v1")
	v1.POST("/accounts", h.CreateAccount)
	v1.GET("/accounts/:accountId", h.GetAccount)

	v1.POST("/transactions", h.CreateTransaction)

	return router
}
