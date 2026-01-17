package http

import (
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/handler"
	"github.com/gin-gonic/gin"
)

type Router struct {
	accountHandler     *handler.AccountHandler
	transactionHandler *handler.TransactionHandler
	healthHandler      *handler.HealthHandler
	engine             *gin.Engine
}

func NewRouter(
	accountHandler *handler.AccountHandler,
	transactionHandler *handler.TransactionHandler,
	healthHandler *handler.HealthHandler,
) *Router {
	return &Router{
		accountHandler:     accountHandler,
		transactionHandler: transactionHandler,
		healthHandler:      healthHandler,
		engine:             gin.Default(),
	}
}

func (r *Router) Setup() *gin.Engine {
	// r.engine.Use(middleware.ErrorHandlerMiddleware())

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
