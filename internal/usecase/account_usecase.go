package usecase

import (
	"context"

	"github.com/alexduzi/challengepismo/internal/domain"
	dto "github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/repository"
	"github.com/alexduzi/challengepismo/internal/usecase/mapper"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, model dto.Account) error
	GetAccountByID(ctx context.Context, id int) (*dto.Account, error)
}

type AccountUseCaseImpl struct {
	accountRepository repository.AccountRepository
}

func NewAccountUseCase(accountRepository repository.AccountRepository) *AccountUseCaseImpl {
	return &AccountUseCaseImpl{
		accountRepository: accountRepository,
	}
}

func (a *AccountUseCaseImpl) CreateAccount(ctx context.Context, model dto.Account) error {
	err := a.accountRepository.Save(ctx, domain.Account{
		DocumentNumber: model.DocumentNumber,
	})

	if err != nil {
		// log error
		return err
	}

	return nil
}

func (a *AccountUseCaseImpl) GetAccountByID(ctx context.Context, id int) (*dto.Account, error) {
	account, err := a.accountRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.ToAccountResponse(account), nil
}
