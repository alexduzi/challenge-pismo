package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/alexduzi/challengepismo/internal/repository"
	"github.com/alexduzi/challengepismo/internal/usecase/mapper"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, request request.CreateTransactionRequest) (*response.TransactionResponse, error)
}

type TransactionUseCaseImpl struct {
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
	logger                logger.Logger
}

func NewTransactionUseCase(
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
	logger logger.Logger) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		logger:                logger,
	}
}

func (t *TransactionUseCaseImpl) CreateTransaction(ctx context.Context, request request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	log := t.logger.WithContext(ctx)

	normalizedAmount, err := domain.NormalizeAmount(request.OperationTypeID, request.Amount)
	if err != nil {
		log.Warn("Use case: Failed to normalize amount",
			slog.Int("operation_type_id", request.OperationTypeID),
		)
		return nil, err
	}

	log.Debug("Use case: Creating transaction",
		slog.Int64("account_id", request.AccountID),
	)

	log.Debug("Use case: Fetching account by ID",
		slog.Int64("account_id", request.AccountID),
	)

	_, err = t.accountRepository.GetByID(ctx, request.AccountID)
	if err != nil {
		log.Warn("Use case: Account not found",
			slog.Int64("account_id", request.AccountID),
		)
		return nil, err
	}

	tran, err := t.transactionRepository.Save(ctx, domain.Transaction{
		AccountID:       request.AccountID,
		OperationTypeID: request.OperationTypeID,
		Amount:          normalizedAmount,
		EventDate:       time.Now(),
	})
	if err != nil {
		log.Error("Use case: Failed to create transaction",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	log.Info("Use case: Transaction created",
		slog.Int64("transaction_id", tran.TransactionID),
	)

	return mapper.ToTransactionResponse(tran), nil
}
