package repository

import (
	"context"

	"github.com/alexduzi/challengepismo/internal/domain"
)

type TransactionRepository interface {
	UpdateForBalance(ctx context.Context, transaction domain.Transaction) error
	GetAllByAccountID(ctx context.Context, accountId int64) ([]domain.Transaction, error)
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	GetByID(ctx context.Context, id int64) (*domain.Transaction, error)
	Save(ctx context.Context, transaction domain.Transaction) (*domain.Transaction, error)
	Update(ctx context.Context, transaction domain.Transaction) error
	Delete(ctx context.Context, id int64) error
}
