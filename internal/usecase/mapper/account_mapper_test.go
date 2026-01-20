package mapper

import (
	"testing"
	"time"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAccountMapper_CanMapDomainToDto(t *testing.T) {
	acc := &domain.Account{
		AccountID:      1,
		DocumentNumber: "123456789",
		FullName:       "John Doe",
		Email:          "john.doe@test.com",
		Phone:          "111112345",
		AccountType:    "savings",
		Balance:        50.0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	dto := ToAccountResponse(acc)

	assert.NotNil(t, dto)
	assert.Equal(t, acc.AccountID, dto.AccountID)
	assert.Equal(t, acc.DocumentNumber, dto.DocumentNumber)
	assert.Equal(t, acc.FullName, dto.FullName)
	assert.Equal(t, acc.Email, dto.Email)
	assert.Equal(t, acc.Phone, dto.Phone)
	assert.Equal(t, acc.AccountType, dto.AccountType)
	assert.Equal(t, acc.Balance, dto.Balance)
	assert.Equal(t, acc.CreatedAt, dto.CreatedAt)
	assert.Equal(t, acc.UpdatedAt, dto.UpdatedAt)
}
