package db

import (
	"context"
	"fmt"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	GetByID(ctx context.Context, id int) (*domain.Transaction, error)
	Save(ctx context.Context, transaction domain.Transaction) error
	Update(ctx context.Context, transaction domain.Transaction) error
	Delete(ctx context.Context, id int) error
}

type TransactionRepositoryImpl struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db}
}

func (a *TransactionRepositoryImpl) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	query := `
		SELECT * FROM transactions
	`
	err := a.db.SelectContext(ctx, &transactions, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return transactions, nil
}

func (a *TransactionRepositoryImpl) GetByID(ctx context.Context, id int) (*domain.Transaction, error) {
	var transaction domain.Transaction
	query := `
		SELECT * FROM transactions WHERE transaction_id = $1
	`
	err := a.db.QueryRowContext(ctx, query, id).Scan(&transaction)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return &transaction, nil
}

func (a *TransactionRepositoryImpl) Save(ctx context.Context, transaction domain.Transaction) error {
	query := `
		INSERT INTO transactions (account_id, operation, amount, event_date)
		VALUES ($1, $2)
	`
	_, err := a.db.ExecContext(ctx, query, transaction.AccountID, transaction.OperationTypeID, transaction.Amount, transaction.EventDate)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}

func (a *TransactionRepositoryImpl) Update(ctx context.Context, transaction domain.Transaction) error {
	query := `
		UPDATE transactions
		SET account_id = $2,
			operation = $3,
			amount = $4
		WHERE transaction_id = $1
	`
	_, err := a.db.ExecContext(ctx, query, transaction.TransactionID, transaction.AccountID, transaction.OperationTypeID, transaction.Amount)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}

func (a *TransactionRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM transactions WHERE transaction_id = $1
	`
	_, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}
