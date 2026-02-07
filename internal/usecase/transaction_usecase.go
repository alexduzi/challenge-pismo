package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/alexduzi/challengepismo/internal/repository"
	"github.com/alexduzi/challengepismo/internal/usecase/mapper"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, request request.CreateTransactionRequest) (*response.TransactionResponse, error)
}

type TransactionUseCaseImpl struct {
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
	logger                logger.Logger
}

func NewTransactionUseCase(
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
	logger logger.Logger) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		logger:                logger,
	}
}

func (t *TransactionUseCaseImpl) CreateTransaction(ctx context.Context, request request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	log := t.logger.WithContext(ctx)

	normalizedAmount, err := domain.NormalizeAmount(request.OperationTypeID, request.Amount)
	if err != nil {
		log.Warn("Use case: Failed to normalize amount",
			slog.Int("operation_type_id", request.OperationTypeID),
		)
		return nil, err
	}

	log.Debug("Use case: Fetching account by ID",
		slog.Int64("account_id", request.AccountID),
	)

	_, err = t.accountRepository.GetByID(ctx, request.AccountID)
	if err != nil {
		log.Warn("Use case: Account not found",
			slog.Int64("account_id", request.AccountID),
		)
		return nil, err
	}

	trans, err := t.transactionRepository.GetAllByAccountIDTx(ctx, request.AccountID)
	if err != nil {
		log.Warn("Use case: Transactions not found",
			slog.Int64("account_id", request.AccountID),
		)
		return nil, err
	}

	newBalance := request.Amount

	if len(trans) > 0 && request.OperationTypeID == domain.CreditVoucher {
		for _, tran := range trans {

			discharge := newBalance
			if discharge > -tran.Balance {
				discharge = -tran.Balance
			}

			newBal := tran.Balance + discharge
			tran.Balance = newBal

			log.Debug("Use case: Updating transaction",
				slog.Int64("account_id", request.AccountID),
			)

			err = t.transactionRepository.UpdateForBalanceTx(ctx, tran)
			if err != nil {
				// log.Warn("Use case: Transactions not found",
				// 	slog.Int64("account_id", request.AccountID),
				// )
				return nil, err
			}

			newBalance -= discharge
			if newBalance == 0 {
				break
			}
		}
	}

	log.Debug("Use case: Creating transaction",
		slog.Int64("account_id", request.AccountID),
	)

	tran, err := t.transactionRepository.SaveTx(ctx, domain.Transaction{
		AccountID:       request.AccountID,
		OperationTypeID: request.OperationTypeID,
		Amount:          normalizedAmount,
		Balance:         newBalance, // balance corrigido
		EventDate:       time.Now(),
	})

	if err != nil {
		log.Error("Use case: Failed to create transaction",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	log.Info("Use case: Transaction created",
		slog.Int64("transaction_id", tran.TransactionID),
	)

	return mapper.ToTransactionResponse(tran), nil
}
