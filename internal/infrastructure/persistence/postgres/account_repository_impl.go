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

type AccountRepositoryImpl struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewAccountRepository(db *sqlx.DB, logger logger.Logger) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (a *AccountRepositoryImpl) GetAll(ctx context.Context) ([]domain.Account, error) {
	var accounts []domain.Account
	query := `
		SELECT * FROM accounts
	`
	err := a.db.SelectContext(ctx, &accounts, query)
	if err != nil {
		return nil, handlePgError(err)
	}

	return accounts, nil
}

func (a *AccountRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	a.logger.Debug("Executing SELECT query", slog.String("table", "accounts"))

	var account domain.Account
	query := `
		SELECT * FROM accounts WHERE account_id = $1
	`
	err := a.db.GetContext(ctx, &account, query, id)
	if err != nil {
		a.logger.Error("SELECT query failed", slog.String("error", err.Error()))

		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) Save(ctx context.Context, account domain.Account) (*domain.Account, error) {
	a.logger.Debug("Executing INSERT query", slog.String("table", "accounts"))

	query := `
		INSERT INTO accounts (document_number, full_name, email, phone, account_type)
		VALUES (:document_number, :full_name, :email, :phone, :account_type)
		RETURNING account_id, created_at, updated_at
	`

	rows, err := a.db.NamedQueryContext(ctx, query, account)
	if err != nil {
		a.logger.Error("INSERT query failed", slog.String("error", err.Error()))
		return nil, handlePgError(err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&account.AccountID, &account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, handlePgError(err)
		}
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) Update(ctx context.Context, account domain.Account) (*domain.Account, error) {
	query := `
		UPDATE accounts
		SET full_name = :full_name,
			email = :email,
			phone = :phone,
			account_type = :account_type,
			updated_at = NOW()
		WHERE account_id = :account_id
		RETURNING updated_at
	`
	rows, err := a.db.NamedQueryContext(ctx, query, account)
	if err != nil {
		return nil, handlePgError(err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&account.UpdatedAt); err != nil {
			return nil, handlePgError(err)
		}
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM accounts WHERE account_id = $1
	`
	_, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return handlePgError(err)
	}

	return nil
}
