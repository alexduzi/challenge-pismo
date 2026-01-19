//go:build wireinject
// +build wireinject

package main

import (
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/alexduzi/challengepismo/internal/infrastructure/http/handler"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/alexduzi/challengepismo/internal/infrastructure/persistence/postgres"
	"github.com/alexduzi/challengepismo/internal/repository"
	"github.com/alexduzi/challengepismo/internal/usecase"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

var setAccountRepositoryDependency = wire.NewSet(
	postgres.NewAccountRepository,
	wire.Bind(new(repository.AccountRepository), new(*postgres.AccountRepositoryImpl)),
)

var setTransactionRepositoryDependency = wire.NewSet(
	postgres.NewTransactionRepository,
	wire.Bind(new(repository.TransactionRepository), new(*postgres.TransactionRepositoryImpl)),
)

var setAccountUseCaseDependency = wire.NewSet(
	usecase.NewAccountUseCase,
	wire.Bind(new(usecase.AccountUseCase), new(*usecase.AccountUseCaseImpl)),
)

var setTransactionUseCaseDependency = wire.NewSet(
	usecase.NewTransactionUseCase,
	wire.Bind(new(usecase.TransactionUseCase), new(*usecase.TransactionUseCaseImpl)),
)

var setHandlerDependency = wire.NewSet(
	handler.NewAccountHandler,
	handler.NewTransactionHandler,
	handler.NewHealthHandler,
)

var setRepositoryDependency = wire.NewSet(
	setAccountRepositoryDependency,
	setTransactionRepositoryDependency,
)

var setUseCaseDependency = wire.NewSet(
	setAccountUseCaseDependency,
	setTransactionUseCaseDependency,
)

func InitializeApplication(cfg *config.Config, db *sqlx.DB, log logger.Logger) *Application {
	wire.Build(
		setRepositoryDependency,
		setUseCaseDependency,
		setHandlerDependency,
		NewApplication,
	)
	return nil
}
