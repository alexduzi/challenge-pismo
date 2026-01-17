package mapper

import (
	"github.com/alexduzi/challengepismo/internal/domain"
	dto "github.com/alexduzi/challengepismo/internal/dto/request"
)

func ToAccountResponse(account *domain.Account) *dto.Account {
	return &dto.Account{
		DocumentNumber: account.DocumentNumber,
	}
}
