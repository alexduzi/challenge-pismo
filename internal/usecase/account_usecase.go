package usecase

import (
	"context"
	"log/slog"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/alexduzi/challengepismo/internal/repository"
	"github.com/alexduzi/challengepismo/internal/usecase/mapper"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, request request.CreateAccountRequest) (*response.AccountResponse, error)
	GetAccountByID(ctx context.Context, id int64) (*response.AccountResponse, error)
}

type AccountUseCaseImpl struct {
	accountRepository repository.AccountRepository
	logger            logger.Logger
}

func NewAccountUseCase(accountRepository repository.AccountRepository, logger logger.Logger) *AccountUseCaseImpl {
	return &AccountUseCaseImpl{
		accountRepository: accountRepository,
		logger:            logger,
	}
}

func (a *AccountUseCaseImpl) CreateAccount(ctx context.Context, request request.CreateAccountRequest) (*response.AccountResponse, error) {
	log := a.logger.WithContext(ctx)

	log.Debug("Use case: Creating account",
		slog.String("document_number", request.DocumentNumber),
	)

	account, err := a.accountRepository.Save(ctx, domain.Account{
		DocumentNumber: request.DocumentNumber,
		FullName:       request.FullName,
		Email:          request.Email,
		Phone:          request.Phone,
		AccountType:    request.AccountType,
	})

	if err != nil {
		log.Error("Use case: Failed to save account",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	log.Info("Use case: Account created",
		slog.Int64("account_id", account.AccountID),
	)

	return mapper.ToAccountResponse(account), nil
}

func (a *AccountUseCaseImpl) GetAccountByID(ctx context.Context, id int64) (*response.AccountResponse, error) {
	log := a.logger.WithContext(ctx)

	log.Debug("Use case: Fetching account by ID",
		slog.Int64("account_id", id),
	)

	account, err := a.accountRepository.GetByID(ctx, id)
	if err != nil {
		log.Warn("Use case: Account not found",
			slog.Int64("account_id", id),
		)
		return nil, err
	}

	return mapper.ToAccountResponse(account), nil
}
