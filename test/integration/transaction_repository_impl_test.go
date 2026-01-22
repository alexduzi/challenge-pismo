package integration

import (
	"context"
	"testing"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	conn "github.com/alexduzi/challengepismo/internal/infrastructure/persistence/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	ctx                   context.Context
	postgresContainer     *postgres.PostgresContainer
	db                    *sqlx.DB
	logger                logger.Logger
	accountRepository     *conn.AccountRepositoryImpl
	transactionRepository *conn.TransactionRepositoryImpl
	account               *domain.Account
	transaction           *domain.Transaction
}

func (suite *TransactionRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	conf, _ := config.LoadConfig()

	db, container, err := setupTestDatabase(suite.ctx, conf)
	if err != nil {
		panic(err)
	}

	suite.postgresContainer = container

	suite.account = &domain.Account{
		DocumentNumber: "12345678918",
		FullName:       "John Doe",
		Email:          "test10@test.com",
		Phone:          "9999999999",
		AccountType:    "savings",
	}

	suite.logger = logger.NewLogger(conf)

	suite.accountRepository = conn.NewAccountRepository(db, suite.logger)
	suite.transactionRepository = conn.NewTransactionRepository(db, suite.logger)
}

func (suite *TransactionRepositoryTestSuite) TearDownSuite() {
	if err := suite.postgresContainer.Terminate(suite.ctx); err != nil {
		suite.logger.Error("Error terminating postgres container: ", err)
	}
	if err := conn.CloseDB(suite.db); err != nil {
		suite.logger.Error("Error terminating sqlx connection: ", err)
	}
}

func (suite *TransactionRepositoryTestSuite) Test01_TransactionRepository_CanCreateTransaction() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test")
	}

	account, err := suite.accountRepository.Save(suite.ctx, *suite.account)
	suite.account = account

	assert.NotNil(suite.T(), account)
	assert.NoError(suite.T(), err)

	suite.transaction = &domain.Transaction{
		AccountID:       account.AccountID,
		OperationTypeID: domain.NormalPurchase,
		Amount:          -50.0,
	}

	transaction, err := suite.transactionRepository.Save(suite.ctx, *suite.transaction)
	assert.NotNil(suite.T(), transaction)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), account.AccountID, transaction.AccountID)
	assert.Equal(suite.T(), domain.NormalPurchase, transaction.OperationTypeID)
	assert.Equal(suite.T(), float64(-50.0), transaction.Amount)
}

func (suite *TransactionRepositoryTestSuite) Test02_TransactionRepository_CanCreateTransaction_ErrAccountNotExistsForTransaction() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test")
	}

	suite.transaction = &domain.Transaction{
		AccountID:       99,
		OperationTypeID: domain.NormalPurchase,
		Amount:          -50.0,
	}

	transaction, err := suite.transactionRepository.Save(suite.ctx, *suite.transaction)
	assert.Nil(suite.T(), transaction)
	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, exception.ErrAccountNotExistsForTransaction)
}

func (suite *TransactionRepositoryTestSuite) Test03_TransactionRepository_CanCreateTransaction_ErrInvalidAmountForOperationType() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test")
	}

	suite.transaction = &domain.Transaction{
		AccountID:       suite.account.AccountID,
		OperationTypeID: domain.NormalPurchase,
		Amount:          150.0,
	}

	transaction, err := suite.transactionRepository.Save(suite.ctx, *suite.transaction)
	assert.Nil(suite.T(), transaction)
	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, exception.ErrInvalidAmountForOperationType)
}

func TestTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
