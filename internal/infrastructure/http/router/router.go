package router

import (
	_ "github.com/alexduzi/challengepismo/docs" // swagger docs
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/handler"
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/middleware"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	cfg                *config.Config
	logger             logger.Logger
	accountHandler     *handler.AccountHandler
	transactionHandler *handler.TransactionHandler
	healthHandler      *handler.HealthHandler
	engine             *gin.Engine
}

func NewRouter(
	cfg *config.Config,
	logger logger.Logger,
	accountHandler *handler.AccountHandler,
	transactionHandler *handler.TransactionHandler,
	healthHandler *handler.HealthHandler,
) *Router {
	return &Router{
		cfg:                cfg,
		logger:             logger,
		accountHandler:     accountHandler,
		transactionHandler: transactionHandler,
		healthHandler:      healthHandler,
		engine:             gin.New(),
	}
}

func (r *Router) Setup() *gin.Engine {
	gin.SetMode(r.cfg.GinMode)

	r.engine.Use(middleware.RequestIDMiddleware())
	r.engine.Use(middleware.LoggerMiddleware(r.logger))
	r.engine.Use(gin.Recovery())
	r.engine.Use(middleware.ErrorHandlerMiddleware())

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health checks
	r.engine.GET("/health", r.healthHandler.GetStatus)
	r.engine.GET("/readiness", r.healthHandler.GetStatus)

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		// Account routes
		accounts := v1.Group("/accounts")
		{
			accounts.POST("", r.accountHandler.CreateAccount)
			accounts.GET("/:accountId", r.accountHandler.GetAccountByID)
		}

		// Transaction routes
		transactions := v1.Group("/transactions")
		{
			transactions.POST("", r.transactionHandler.CreateTransaction)
		}
	}

	return r.engine
}
