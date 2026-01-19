package mapper

import (
	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/response"
)

func ToTransactionResponse(transaction *domain.Transaction) *response.TransactionResponse {
	return &response.TransactionResponse{
		TransactionID:   transaction.TransactionID,
		AccountID:       transaction.AccountID,
		OperationTypeID: transaction.OperationTypeID,
		Amount:          transaction.Amount,
		EventDate:       transaction.EventDate,
		CreatedAt:       transaction.CreatedAt,
	}
}
