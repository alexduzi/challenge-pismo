package usecase

import (
	"context"
	"fmt"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/infrastructure/db"
	"github.com/alexduzi/challengepismo/internal/model"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, request model.Transaction) error
}

type TransactionUseCaseImpl struct {
	accountRepository     db.AccountRepository
	transactionRepository db.TransactionRepository
}

func NewTransactionUseCase(
	accountRepository db.AccountRepository,
	transactionRepository db.TransactionRepository) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
	}
}

func (t *TransactionUseCaseImpl) CreateTransaction(ctx context.Context, request model.Transaction) error {
	acc, err := t.accountRepository.GetByID(ctx, request.AccountID)
	if err != nil {
		return err
	}

	if acc == nil {
		return fmt.Errorf("account with id %d not found", request.AccountID)
	}

	if _, err := model.GetOperationType(model.OperationType(request.OperationTypeID)); err != nil {
		return fmt.Errorf("can't convert operation type %w", err)
	}

	return t.transactionRepository.Save(ctx, domain.Transaction{
		AccountID:       request.AccountID,
		OperationTypeID: request.OperationTypeID,
		Amount:          request.Amount,
		EventDate:       request.EventDate,
	})
}
