package main

import (
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/handler"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
)

type Application struct {
	Config             *config.Config
	Logger             logger.Logger
	AccountHandler     *handler.AccountHandler
	TransactionHandler *handler.TransactionHandler
	HealthHandler      *handler.HealthHandler
}

func NewApplication(
	cfg *config.Config,
	log logger.Logger,
	accountHandler *handler.AccountHandler,
	transactionHandler *handler.TransactionHandler,
	healthHandler *handler.HealthHandler,
) *Application {
	return &Application{
		Config:             cfg,
		Logger:             log,
		AccountHandler:     accountHandler,
		TransactionHandler: transactionHandler,
		HealthHandler:      healthHandler,
	}
}
