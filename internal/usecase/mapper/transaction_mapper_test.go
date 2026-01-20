package mapper

import (
	"testing"
	"time"

	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTransactionMapper_CanMapDomainToDto(t *testing.T) {
	tran := &domain.Transaction{
		TransactionID:   1,
		AccountID:       1,
		OperationTypeID: domain.NormalPurchase,
		Amount:          50.0,
		EventDate:       time.Now(),
		CreatedAt:       time.Now(),
	}

	dto := ToTransactionResponse(tran)

	assert.NotNil(t, dto)
	assert.Equal(t, tran.TransactionID, dto.TransactionID)
	assert.Equal(t, tran.AccountID, dto.AccountID)
	assert.Equal(t, tran.OperationTypeID, dto.OperationTypeID)
	assert.Equal(t, tran.Amount, dto.Amount)
	assert.Equal(t, tran.EventDate, dto.EventDate)
	assert.Equal(t, tran.CreatedAt, dto.CreatedAt)
}
