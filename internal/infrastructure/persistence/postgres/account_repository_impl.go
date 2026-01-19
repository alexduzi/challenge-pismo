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

type AccountRepositoryImpl struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{db}
}

func (a *AccountRepositoryImpl) GetAll(ctx context.Context) ([]domain.Account, error) {
	var accounts []domain.Account
	query := `
		SELECT * FROM accounts
	`
	err := a.db.SelectContext(ctx, &accounts, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return accounts, nil
}

func (a *AccountRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	var account domain.Account
	query := `
		SELECT * FROM accounts WHERE account_id = $1
	`
	err := a.db.GetContext(ctx, &account, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) Save(ctx context.Context, account domain.Account) (*domain.Account, error) {
	query := `
		INSERT INTO accounts (document_number, full_name, email, phone, account_type)
		VALUES (:document_number, :full_name, :email, :phone, :account_type)
		RETURNING account_id, created_at, updated_at
	`

	rows, err := a.db.NamedQueryContext(ctx, query, account)
	if err != nil {
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

func (a *AccountRepositoryImpl) Update(ctx context.Context, account domain.Account) error {
	query := `
		UPDATE accounts
		SET document_number = $2
		WHERE account_id = $1
	`
	_, err := a.db.ExecContext(ctx, query, account.AccountID, account.DocumentNumber)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}

func (a *AccountRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM accounts WHERE account_id = $1
	`
	_, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}
