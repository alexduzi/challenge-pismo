package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/handler"
	httpRouter "github.com/alexduzi/challengepismo/internal/infrastructure/http/router"
	"github.com/alexduzi/challengepismo/internal/infrastructure/persistence/postgres"
	"github.com/alexduzi/challengepismo/internal/usecase"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := postgres.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}
	defer func() {
		if err := postgres.CloseDB(db); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	accountRepo := postgres.NewAccountRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	accountUseCase := usecase.NewAccountUseCase(accountRepo)
	transactionUseCase := usecase.NewTransactionUseCase(accountRepo, transactionRepo)

	accountHandler := handler.NewAccountHandler(accountUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	healthHandler := handler.NewHealthHandler(cfg)

	router := httpRouter.NewRouter(
		cfg,
		accountHandler,
		transactionHandler,
		healthHandler,
	)
	engine := router.Setup()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: engine.Handler(),
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
