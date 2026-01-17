package db

import (
	"context"
	"fmt"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	GetAll(ctx context.Context) ([]domain.Account, error)
	GetByID(ctx context.Context, id int) (*domain.Account, error)
	Save(ctx context.Context, account domain.Account) error
	Update(ctx context.Context, account domain.Account) error
	Delete(ctx context.Context, id int) error
}

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

func (a *AccountRepositoryImpl) GetByID(ctx context.Context, id int) (*domain.Account, error) {
	var account domain.Account
	query := `
		SELECT * FROM accounts WHERE account_id = $1
	`
	err := a.db.QueryRowContext(ctx, query, id).Scan(&account)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) Save(ctx context.Context, account domain.Account) error {
	query := `
		INSERT INTO accounts (account_id, document_number)
		VALUES ($1, $2)
	`
	_, err := a.db.ExecContext(ctx, query, account.AccountID, account.DocumentNumber)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
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

func (a *AccountRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM accounts WHERE account_id = $1
	`
	_, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
	}

	return nil
}
