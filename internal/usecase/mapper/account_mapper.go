package mapper

import (
	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/dto/response"
)

func ToAccountResponse(account *domain.Account) *response.AccountResponse {
	return &response.AccountResponse{
		AccountID:      account.AccountID,
		DocumentNumber: account.DocumentNumber,
		FullName:       account.FullName,
		Email:          account.Email,
		Phone:          account.Phone,
		AccountType:    account.AccountType,
		Balance:        account.Balance,
		CreatedAt:      account.CreatedAt,
		UpdatedAt:      account.UpdatedAt,
	}
}
