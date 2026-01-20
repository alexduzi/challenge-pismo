package router

import (
	"testing"

	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/handler"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/alexduzi/challengepismo/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestRouter_CanCreateRoutes(t *testing.T) {
	cfg, _ := config.LoadConfig()
	logger := logger.NewMockLogger(t)

	accHandler := handler.NewAccountHandler(usecase.NewMockAccountUseCase(t), logger)
	tranHandler := handler.NewTransactionHandler(usecase.NewMockTransactionUseCase(t), logger)
	healthHandler := handler.NewHealthHandler(cfg)

	router := NewRouter(cfg, logger, accHandler, tranHandler, healthHandler)

	engine := router.Setup()

	assert.NotNil(t, engine)
}
