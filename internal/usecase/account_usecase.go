package usecase

import (
	"context"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/infrastructure/db"
	"github.com/alexduzi/challengepismo/internal/model"
	"github.com/alexduzi/challengepismo/internal/model/mapper"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, model model.Account) error
	GetAccountByID(ctx context.Context, id int) (*model.Account, error)
}

type AccountUseCaseImpl struct {
	accountRepository db.AccountRepository
}

func NewAccountUseCase(accountRepository db.AccountRepository) *AccountUseCaseImpl {
	return &AccountUseCaseImpl{
		accountRepository: accountRepository,
	}
}

func (a *AccountUseCaseImpl) CreateAccount(ctx context.Context, model model.Account) error {
	err := a.accountRepository.Save(ctx, domain.Account{
		DocumentNumber: model.DocumentNumber,
	})

	if err != nil {
		// log error
		return err
	}

	return nil
}

func (a *AccountUseCaseImpl) GetAccountByID(ctx context.Context, id int) (*model.Account, error) {
	account, err := a.accountRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.ToAccountResponse(account), nil
}
