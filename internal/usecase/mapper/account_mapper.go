package mapper

import (
	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/response"
)

func ToAccountResponse(account *domain.Account) *response.AccountResponse {
	return &response.AccountResponse{
		AccountID:      account.AccountID,
		DocumentNumber: account.DocumentNumber,
	}
}
