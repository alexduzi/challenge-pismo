package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/jmoiron/sqlx"
)

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

func (a *TransactionRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	var transaction domain.Transaction
	query := `
		SELECT * FROM transactions WHERE transaction_id = $1
	`
	err := a.db.GetContext(ctx, &transaction, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return &transaction, nil
}

func (a *TransactionRepositoryImpl) Save(ctx context.Context, transaction domain.Transaction) (*domain.Transaction, error) {
	query := `
		INSERT INTO transactions (account_id, operation_type_id, amount)
		VALUES (:account_id, :operation_type_id, :amount)
		RETURNING transaction_id, event_date, created_at
	`
	rows, err := a.db.NamedQueryContext(ctx, query, transaction)
	if err != nil {
		return nil, handlePgError(err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&transaction.TransactionID, &transaction.EventDate, &transaction.CreatedAt); err != nil {
			return nil, handlePgError(err)
		}
	}

	return &transaction, nil
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

func (a *TransactionRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM transactions WHERE transaction_id = $1
	`
	_, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}
