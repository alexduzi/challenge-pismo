package usecase

import (
	"context"
	"time"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/repository"
	"github.com/alexduzi/challengepismo/internal/usecase/mapper"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, request request.CreateTransactionRequest) (*response.TransactionResponse, error)
}

type TransactionUseCaseImpl struct {
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
}

func NewTransactionUseCase(
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
	}
}

func (t *TransactionUseCaseImpl) CreateTransaction(ctx context.Context, request request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	_, err := t.accountRepository.GetByID(ctx, request.AccountID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateOperationTypeForAmount(request.OperationTypeID, request.Amount); err != nil {
		return nil, err
	}

	tran, err := t.transactionRepository.Save(ctx, domain.Transaction{
		AccountID:       request.AccountID,
		OperationTypeID: request.OperationTypeID,
		Amount:          request.Amount,
		EventDate:       time.Now(),
	})
	if err != nil {
		// log error
		return nil, err
	}

	return mapper.ToTransactionResponse(tran), nil
}
