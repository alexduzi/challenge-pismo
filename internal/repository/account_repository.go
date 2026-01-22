package repository

import (
	"context"

	"github.com/alexduzi/challengepismo/internal/domain"
)

type AccountRepository interface {
	GetAll(ctx context.Context) ([]domain.Account, error)
	GetByID(ctx context.Context, id int64) (*domain.Account, error)
	Save(ctx context.Context, account domain.Account) (*domain.Account, error)
	Update(ctx context.Context, account domain.Account) (*domain.Account, error)
	Delete(ctx context.Context, id int64) error
}
