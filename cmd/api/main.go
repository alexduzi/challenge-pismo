package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	httpRouter "github.com/alexduzi/challengepismo/internal/infrastructure/http/router"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/alexduzi/challengepismo/internal/infrastructure/persistence/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg)

	db, err := postgres.ConnectDB(cfg)
	if err != nil {
		log.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		log.Info("Closing database connection")
		if err := postgres.CloseDB(db); err != nil {
			log.Error("Error closing database", slog.String("error", err.Error()))
		}
	}()

	app := InitializeApplication(cfg, db, log)

	router := httpRouter.NewRouter(
		app.Config,
		app.Logger,
		app.AccountHandler,
		app.TransactionHandler,
		app.HealthHandler,
	)
	engine := router.Setup()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engine.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Server starting",
			slog.String("port", cfg.Port),
			slog.String("environment", cfg.GinMode),
			slog.String("health_check", fmt.Sprintf("http://localhost:%s/health", cfg.Port)),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("Server stopped gracefully")
}
