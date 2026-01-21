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

type TransactionUseCaseTestSuite struct {
	suite.Suite
	logger                *logger.MockLogger
	accountRepository     *repository.MockAccountRepository
	transactionRepository *repository.MockTransactionRepository
	transactionUseCase    *TransactionUseCaseImpl
}

func (s *TransactionUseCaseTestSuite) SetupTest() {
	s.logger = logger.NewMockLogger(s.T())
	s.accountRepository = repository.NewMockAccountRepository(s.T())
	s.transactionRepository = repository.NewMockTransactionRepository(s.T())
	s.transactionUseCase = NewTransactionUseCase(s.accountRepository, s.transactionRepository, s.logger)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_Success() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.CreditVoucher,
		Amount:          100.50,
	}

	expectedAccount := &domain.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
	}

	expectedTransaction := &domain.Transaction{
		TransactionID:   1,
		AccountID:       1,
		OperationTypeID: domain.CreditVoucher,
		Amount:          100.50,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating transaction", mock.Anything).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Info("Use case: Transaction created", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, transactionRequest.AccountID).Return(expectedAccount, nil).Once()

	s.transactionRepository.EXPECT().Save(ctx, mock.MatchedBy(func(t domain.Transaction) bool {
		return t.AccountID == transactionRequest.AccountID &&
			t.OperationTypeID == transactionRequest.OperationTypeID &&
			t.Amount == transactionRequest.Amount
	})).Return(expectedTransaction, nil).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), transactionResponse)
	assert.Equal(s.T(), expectedTransaction.TransactionID, transactionResponse.TransactionID)
	assert.Equal(s.T(), expectedTransaction.AccountID, transactionResponse.AccountID)
	assert.Equal(s.T(), expectedTransaction.OperationTypeID, transactionResponse.OperationTypeID)
	assert.Equal(s.T(), expectedTransaction.Amount, transactionResponse.Amount)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_NormalPurchase_Success() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.NormalPurchase,
		Amount:          -50.00,
	}

	expectedAccount := &domain.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
	}

	expectedTransaction := &domain.Transaction{
		TransactionID:   2,
		AccountID:       1,
		OperationTypeID: domain.NormalPurchase,
		Amount:          -50.00,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating transaction", mock.Anything).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Info("Use case: Transaction created", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, transactionRequest.AccountID).Return(expectedAccount, nil).Once()

	s.transactionRepository.EXPECT().Save(ctx, mock.MatchedBy(func(t domain.Transaction) bool {
		return t.AccountID == transactionRequest.AccountID &&
			t.OperationTypeID == transactionRequest.OperationTypeID &&
			t.Amount == transactionRequest.Amount
	})).Return(expectedTransaction, nil).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), transactionResponse)
	assert.Equal(s.T(), expectedTransaction.TransactionID, transactionResponse.TransactionID)
	assert.Equal(s.T(), expectedTransaction.Amount, transactionResponse.Amount)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_InvalidAmountForOperationType() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.NormalPurchase,
		Amount:          100.50,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Warn("Use case: Operation type is invalid for this amount", mock.Anything).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), transactionResponse)
	assert.ErrorIs(s.T(), err, exception.ErrInvalidAmountForOperationType)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_CreditVoucherWithNegativeAmount() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.CreditVoucher,
		Amount:          -100.50,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Warn("Use case: Operation type is invalid for this amount", mock.Anything).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), transactionResponse)
	assert.ErrorIs(s.T(), err, exception.ErrInvalidAmountForOperationType)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_CreditVoucherWithZeroAmount() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.CreditVoucher,
		Amount:          0,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Warn("Use case: Operation type is invalid for this amount", mock.Anything).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), transactionResponse)
	assert.ErrorIs(s.T(), err, exception.ErrAmountGreaterThanZero)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_AccountNotFound() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       999,
		OperationTypeID: domain.CreditVoucher,
		Amount:          100.50,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating transaction", mock.Anything).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Warn("Use case: Account not found", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, transactionRequest.AccountID).Return(nil, exception.ErrNotFound).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), transactionResponse)
	assert.ErrorIs(s.T(), err, exception.ErrNotFound)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_RepositoryError() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.CreditVoucher,
		Amount:          100.50,
	}

	expectedAccount := &domain.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating transaction", mock.Anything).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Error("Use case: Failed to create transaction", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, transactionRequest.AccountID).Return(expectedAccount, nil).Once()
	s.transactionRepository.EXPECT().Save(ctx, mock.Anything).Return(nil, exception.ErrDatabaseError).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), transactionResponse)
	assert.ErrorIs(s.T(), err, exception.ErrDatabaseError)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_AccountRepositoryDatabaseError() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.CreditVoucher,
		Amount:          100.50,
	}

	dbError := errors.New("connection refused")

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating transaction", mock.Anything).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Warn("Use case: Account not found", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, transactionRequest.AccountID).Return(nil, dbError).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), transactionResponse)
	assert.Equal(s.T(), dbError, err)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_Withdrawal_Success() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.Withdrawal,
		Amount:          -200.00,
	}

	expectedAccount := &domain.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
	}

	expectedTransaction := &domain.Transaction{
		TransactionID:   3,
		AccountID:       1,
		OperationTypeID: domain.Withdrawal,
		Amount:          -200.00,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating transaction", mock.Anything).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Info("Use case: Transaction created", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, transactionRequest.AccountID).Return(expectedAccount, nil).Once()

	s.transactionRepository.EXPECT().Save(ctx, mock.MatchedBy(func(t domain.Transaction) bool {
		return t.AccountID == transactionRequest.AccountID &&
			t.OperationTypeID == transactionRequest.OperationTypeID &&
			t.Amount == transactionRequest.Amount
	})).Return(expectedTransaction, nil).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), transactionResponse)
	assert.Equal(s.T(), expectedTransaction.TransactionID, transactionResponse.TransactionID)
	assert.Equal(s.T(), domain.Withdrawal, transactionResponse.OperationTypeID)
}

func (s *TransactionUseCaseTestSuite) TestCreateTransaction_PurchaseInstallments_Success() {
	ctx := context.Background()

	transactionRequest := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: domain.PurchaseInstallments,
		Amount:          -300.00,
	}

	expectedAccount := &domain.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
	}

	expectedTransaction := &domain.Transaction{
		TransactionID:   4,
		AccountID:       1,
		OperationTypeID: domain.PurchaseInstallments,
		Amount:          -300.00,
	}

	s.logger.EXPECT().WithContext(ctx).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Use case: Creating transaction", mock.Anything).Once()
	s.logger.EXPECT().Debug("Use case: Fetching account by ID", mock.Anything).Once()
	s.logger.EXPECT().Info("Use case: Transaction created", mock.Anything).Once()

	s.accountRepository.EXPECT().GetByID(ctx, transactionRequest.AccountID).Return(expectedAccount, nil).Once()

	s.transactionRepository.EXPECT().Save(ctx, mock.MatchedBy(func(t domain.Transaction) bool {
		return t.AccountID == transactionRequest.AccountID &&
			t.OperationTypeID == transactionRequest.OperationTypeID &&
			t.Amount == transactionRequest.Amount
	})).Return(expectedTransaction, nil).Once()

	transactionResponse, err := s.transactionUseCase.CreateTransaction(ctx, transactionRequest)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), transactionResponse)
	assert.Equal(s.T(), expectedTransaction.TransactionID, transactionResponse.TransactionID)
	assert.Equal(s.T(), domain.PurchaseInstallments, transactionResponse.OperationTypeID)
}

func TestTransactionUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionUseCaseTestSuite))
}
