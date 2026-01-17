package usecase

import (
	"context"
	"fmt"

	"github.com/alexduzi/challengepismo/internal/domain"
	dto "github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/repository"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, request dto.Transaction) error
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

func (t *TransactionUseCaseImpl) CreateTransaction(ctx context.Context, request dto.Transaction) error {
	acc, err := t.accountRepository.GetByID(ctx, request.AccountID)
	if err != nil {
		return err
	}

	if acc == nil {
		return fmt.Errorf("account with id %d not found", request.AccountID)
	}

	if _, err := domain.GetOperationType(domain.OperationType(request.OperationTypeID)); err != nil {
		return fmt.Errorf("can't convert operation type %w", err)
	}

	return t.transactionRepository.Save(ctx, domain.Transaction{
		AccountID:       request.AccountID,
		OperationTypeID: request.OperationTypeID,
		Amount:          request.Amount,
		EventDate:       request.EventDate,
	})
}
