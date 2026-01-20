package handler

import (
	"net/http"
	"time"

	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints

type HealthHandler struct {
	cfg *config.Config
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		cfg: cfg,
	}
}

// GetStatus returns the health status of the service
// @Summary Health check
// @Description Returns the health status of the API service
// @Tags health
// @Produce json
// @Success 200 {object} response.HealthResponse "Service is healthy"
// @Router /health [get]
// @Router /readiness [get]
func (h *HealthHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, response.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   h.cfg.AppName,
	})
}
