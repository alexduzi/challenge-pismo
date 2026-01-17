package mapper

import (
	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/model"
)

func ToAccountResponse(account *domain.Account) *model.Account {
	return &model.Account{
		DocumentNumber: account.DocumentNumber,
	}
}
