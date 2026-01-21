package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/alexduzi/challengepismo/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AccountUseCaseTestSuite struct {
	suite.Suite
	logger            *logger.MockLogger
	accountRepository *repository.MockAccountRepository
	accountUseCase    *AccountUseCaseImpl
}

func (s *AccountUseCaseTestSuite) SetupTest() {
	s.logger = logger.NewMockLogger(s.T())
	s.accountRepository = repository.NewMockAccountRepository(s.T())
	s.accountUseCase = NewAccountUseCase(s.accountRepository, s.logger)
}

func (s *AccountUseCaseTestSuite) TestCreateAccount_Success() {
	ctx := context.Background()

	accRequest := request.CreateAccountRequest{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	expectedAccount := &domain.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating account", mock.Anything).Once()
	s.logger.EXPECT().Info("Use case: Account created", mock.Anything).Once()

	s.accountRepository.EXPECT().Save(ctx, mock.MatchedBy(func(acc domain.Account) bool {
		return acc.DocumentNumber == accRequest.DocumentNumber &&
			acc.FullName == accRequest.FullName &&
			acc.Email == accRequest.Email &&
			acc.Phone == accRequest.Phone &&
			acc.AccountType == accRequest.AccountType
	})).Return(expectedAccount, nil).Once()

	accResponse, err := s.accountUseCase.CreateAccount(ctx, accRequest)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), accResponse)
	assert.Equal(s.T(), expectedAccount.AccountID, accResponse.AccountID)
	assert.Equal(s.T(), expectedAccount.DocumentNumber, accResponse.DocumentNumber)
	assert.Equal(s.T(), expectedAccount.FullName, accResponse.FullName)
	assert.Equal(s.T(), expectedAccount.Email, accResponse.Email)
}

func (s *AccountUseCaseTestSuite) TestCreateAccount_RepositoryError() {
	ctx := context.Background()

	accRequest := request.CreateAccountRequest{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	expectedError := exception.ErrDatabaseError

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating account", mock.Anything).Once()
	s.logger.EXPECT().Error("Use case: Failed to save account", mock.Anything).Once()

	s.accountRepository.EXPECT().Save(ctx, mock.Anything).Return(nil, expectedError).Once()

	accResponse, err := s.accountUseCase.CreateAccount(ctx, accRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), accResponse)
	assert.Equal(s.T(), expectedError, err)
}

func (s *AccountUseCaseTestSuite) TestCreateAccount_DuplicateDocumentError() {
	ctx := context.Background()

	accRequest := request.CreateAccountRequest{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating account", mock.Anything).Once()
	s.logger.EXPECT().Error("Use case: Failed to save account", mock.Anything).Once()

	s.accountRepository.EXPECT().Save(ctx, mock.Anything).Return(nil, exception.ErrDuplicateDocument).Once()

	accResponse, err := s.accountUseCase.CreateAccount(ctx, accRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), accResponse)
	assert.ErrorIs(s.T(), err, exception.ErrDuplicateDocument)
}

func (s *AccountUseCaseTestSuite) TestCreateAccount_DuplicateEmailError() {
	ctx := context.Background()

	accRequest := request.CreateAccountRequest{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating account", mock.Anything).Once()
	s.logger.EXPECT().Error("Use case: Failed to save account", mock.Anything).Once()

	s.accountRepository.EXPECT().Save(ctx, mock.Anything).Return(nil, exception.ErrDuplicateEmail).Once()

	accResponse, err := s.accountUseCase.CreateAccount(ctx, accRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), accResponse)
	assert.ErrorIs(s.T(), err, exception.ErrDuplicateEmail)
}

func (s *AccountUseCaseTestSuite) TestGetAccountByID_Success() {
	ctx := context.Background()
	accountID := int64(1)

	expectedAccount := &domain.Account{
		AccountID:      accountID,
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
		Balance:        1000.50,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, accountID).Return(expectedAccount, nil).Once()

	accResponse, err := s.accountUseCase.GetAccountByID(ctx, accountID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), accResponse)
	assert.Equal(s.T(), expectedAccount.AccountID, accResponse.AccountID)
	assert.Equal(s.T(), expectedAccount.DocumentNumber, accResponse.DocumentNumber)
	assert.Equal(s.T(), expectedAccount.FullName, accResponse.FullName)
	assert.Equal(s.T(), expectedAccount.Email, accResponse.Email)
}

func (s *AccountUseCaseTestSuite) TestGetAccountByID_NotFound() {
	ctx := context.Background()
	accountID := int64(999)

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Warn("Use case: Account not found", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, accountID).Return(nil, exception.ErrNotFound).Once()

	accResponse, err := s.accountUseCase.GetAccountByID(ctx, accountID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), accResponse)
	assert.ErrorIs(s.T(), err, exception.ErrNotFound)
}

func (s *AccountUseCaseTestSuite) TestGetAccountByID_DatabaseError() {
	ctx := context.Background()
	accountID := int64(1)

	dbError := errors.New("connection refused")

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Warn("Use case: Account not found", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, accountID).Return(nil, dbError).Once()

	accResponse, err := s.accountUseCase.GetAccountByID(ctx, accountID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), accResponse)
	assert.Equal(s.T(), dbError, err)
}

func TestAccountUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AccountUseCaseTestSuite))
}
