package repository

import (
	"context"

	"github.com/alexduzi/challengepismo/internal/domain"
)

type TransactionRepository interface {
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	GetByID(ctx context.Context, id int) (*domain.Transaction, error)
	Save(ctx context.Context, transaction domain.Transaction) error
	Update(ctx context.Context, transaction domain.Transaction) error
	Delete(ctx context.Context, id int) error
}
