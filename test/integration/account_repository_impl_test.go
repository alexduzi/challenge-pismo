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

type AccountRepositoryTestSuite struct {
	suite.Suite
	ctx               context.Context
	postgresContainer *postgres.PostgresContainer
	db                *sqlx.DB
	logger            logger.Logger
	accountRepository *conn.AccountRepositoryImpl
	account           *domain.Account
}

func (suite *AccountRepositoryTestSuite) SetupSuite() {
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
}

func (suite *AccountRepositoryTestSuite) TearDownSuite() {
	if err := suite.postgresContainer.Terminate(suite.ctx); err != nil {
		suite.logger.Error("Error terminating postgres container: ", err)
	}
	if err := conn.CloseDB(suite.db); err != nil {
		suite.logger.Error("Error terminating sqlx connection: ", err)
	}
}

func (suite *AccountRepositoryTestSuite) Test01_AccountRepository_CanCreateAccount() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test")
	}

	account, err := suite.accountRepository.Save(suite.ctx, *suite.account)

	assert.NotNil(suite.T(), account)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "12345678918", account.DocumentNumber)
	assert.Equal(suite.T(), "John Doe", account.FullName)
	assert.Equal(suite.T(), "test10@test.com", account.Email)
	assert.Equal(suite.T(), "savings", account.AccountType)
	assert.Equal(suite.T(), "9999999999", account.Phone)
	assert.Equal(suite.T(), float64(0), account.Balance)
	assert.True(suite.T(), !account.CreatedAt.IsZero())
	assert.True(suite.T(), !account.UpdatedAt.IsZero())

	suite.account = account
}

func (suite *AccountRepositoryTestSuite) Test02_AccountRepository_CanGetByID() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test")
	}

	account, err := suite.accountRepository.GetByID(suite.ctx, suite.account.AccountID)

	assert.NotNil(suite.T(), account)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), suite.account.DocumentNumber, account.DocumentNumber)
	assert.Equal(suite.T(), suite.account.FullName, account.FullName)
	assert.Equal(suite.T(), suite.account.Email, account.Email)
	assert.Equal(suite.T(), suite.account.AccountType, account.AccountType)
	assert.Equal(suite.T(), suite.account.Phone, account.Phone)
	assert.Equal(suite.T(), suite.account.Balance, account.Balance)
	assert.Equal(suite.T(), suite.account.CreatedAt, account.CreatedAt)
	assert.Equal(suite.T(), suite.account.UpdatedAt, account.UpdatedAt)
}

func (suite *AccountRepositoryTestSuite) Test03_AccountRepository_CanUpdateAccount() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test")
	}

	suite.account.FullName = "John Doe updated"

	account, err := suite.accountRepository.Update(suite.ctx, *suite.account)

	assert.NotNil(suite.T(), account)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), suite.account.DocumentNumber, account.DocumentNumber)
	assert.Equal(suite.T(), suite.account.FullName, account.FullName)
	assert.Equal(suite.T(), suite.account.Email, account.Email)
	assert.Equal(suite.T(), suite.account.AccountType, account.AccountType)
	assert.Equal(suite.T(), suite.account.Phone, account.Phone)
	assert.Equal(suite.T(), suite.account.Balance, account.Balance)
	assert.Equal(suite.T(), suite.account.CreatedAt, account.CreatedAt)
	assert.True(suite.T(), account.UpdatedAt.After(suite.account.UpdatedAt))

	suite.account = account
}

func (suite *AccountRepositoryTestSuite) Test04_AccountRepository_CanDeleteAccount() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test")
	}

	err := suite.accountRepository.Delete(suite.ctx, suite.account.AccountID)

	assert.NoError(suite.T(), err)

	account, err := suite.accountRepository.GetByID(suite.ctx, suite.account.AccountID)

	assert.Nil(suite.T(), account)
	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, exception.ErrNotFound)
}

func TestAccountRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}
