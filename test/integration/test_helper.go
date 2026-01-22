package integration

import (
	"context"
	"fmt"

	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/handler"
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/router"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	conn "github.com/alexduzi/challengepismo/internal/infrastructure/persistence/postgres"
	"github.com/alexduzi/challengepismo/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDatabase(ctx context.Context, conf *config.Config) (*sqlx.DB, *postgres.PostgresContainer, error) {
	// The postgres module provides a simplified way to run the container
	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(conf.Database.DBName),
		postgres.WithUsername(conf.Database.User),
		postgres.WithPassword(conf.Database.Password),
		postgres.WithInitScripts("../../migrations/000001_init_schema.up.sql"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		pgContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Establish the database connection
	db, err := conn.ConnectDBWithDSN(conf, connStr)
	if err != nil {
		pgContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, pgContainer, nil
}

type TestServer struct {
	Logger logger.Logger
	Engine *gin.Engine
	DB     *sqlx.DB
}

func setupHttpHandler(ctx context.Context, conf *config.Config) (*TestServer, *postgres.PostgresContainer) {
	logger := logger.NewLogger(conf)

	db, container, err := setupTestDatabase(ctx, conf)
	if err != nil {
		panic(err)
	}

	accountRepository := conn.NewAccountRepository(db, logger)
	transactionRepository := conn.NewTransactionRepository(db, logger)

	accountUseCase := usecase.NewAccountUseCase(accountRepository, logger)
	transactionUseCase := usecase.NewTransactionUseCase(accountRepository, transactionRepository, logger)

	accHandler := handler.NewAccountHandler(accountUseCase, logger)
	tranHandler := handler.NewTransactionHandler(transactionUseCase, logger)
	healthHandler := handler.NewHealthHandler(conf)

	router := router.NewRouter(conf, logger, accHandler, tranHandler, healthHandler)
	engine := router.Setup()

	return &TestServer{
		Logger: logger,
		Engine: engine,
		DB:     db,
	}, container
}
