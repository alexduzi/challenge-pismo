package handler

import (
	"net/http"
	"time"

	dto "github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	cfg *config.Config
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		cfg: cfg,
	}
}

func (h *HealthHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   h.cfg.AppName,
	})
}
