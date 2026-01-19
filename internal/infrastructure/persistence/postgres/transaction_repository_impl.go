package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/jmoiron/sqlx"
)

type TransactionRepositoryImpl struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewTransactionRepository(db *sqlx.DB, logger logger.Logger) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (t *TransactionRepositoryImpl) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	query := `
		SELECT * FROM transactions
	`
	err := t.db.SelectContext(ctx, &transactions, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return transactions, nil
}

func (t *TransactionRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	var transaction domain.Transaction
	query := `
		SELECT * FROM transactions WHERE transaction_id = $1
	`
	err := t.db.GetContext(ctx, &transaction, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return &transaction, nil
}

func (t *TransactionRepositoryImpl) Save(ctx context.Context, transaction domain.Transaction) (*domain.Transaction, error) {
	t.logger.Debug("Executing SELECT query", slog.String("table", "transactions"))

	query := `
		INSERT INTO transactions (account_id, operation_type_id, amount)
		VALUES (:account_id, :operation_type_id, :amount)
		RETURNING transaction_id, event_date, created_at
	`
	rows, err := t.db.NamedQueryContext(ctx, query, transaction)
	if err != nil {
		t.logger.Error("INSERT query failed", slog.String("error", err.Error()))
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

func (t *TransactionRepositoryImpl) Update(ctx context.Context, transaction domain.Transaction) error {
	query := `
		UPDATE transactions
		SET account_id = $2,
			operation = $3,
			amount = $4
		WHERE transaction_id = $1
	`
	_, err := t.db.ExecContext(ctx, query, transaction.TransactionID, transaction.AccountID, transaction.OperationTypeID, transaction.Amount)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}

func (t *TransactionRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM transactions WHERE transaction_id = $1
	`
	_, err := t.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}
