package usecase

import (
	"context"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/repository"
	"github.com/alexduzi/challengepismo/internal/usecase/mapper"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, model request.CreateAccountRequest) (*response.AccountResponse, error)
	GetAccountByID(ctx context.Context, id int) (*response.AccountResponse, error)
}

type AccountUseCaseImpl struct {
	accountRepository repository.AccountRepository
}

func NewAccountUseCase(accountRepository repository.AccountRepository) *AccountUseCaseImpl {
	return &AccountUseCaseImpl{
		accountRepository: accountRepository,
	}
}

func (a *AccountUseCaseImpl) CreateAccount(ctx context.Context, model request.CreateAccountRequest) (*response.AccountResponse, error) {
	account, err := a.accountRepository.Save(ctx, domain.Account{
		DocumentNumber: model.DocumentNumber,
		FullName:       model.FullName,
		Email:          model.Email,
		Phone:          model.Phone,
		AccountType:    model.AccountType,
	})

	if err != nil {
		// log error
		return nil, err
	}

	return mapper.ToAccountResponse(account), nil
}

func (a *AccountUseCaseImpl) GetAccountByID(ctx context.Context, id int) (*response.AccountResponse, error) {
	account, err := a.accountRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.ToAccountResponse(account), nil
}
