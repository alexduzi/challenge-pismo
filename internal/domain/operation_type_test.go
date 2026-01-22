package domain

import (
	"testing"

	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeAmount(t *testing.T) {
	tests := []struct {
		name           string
		operationType  int
		inputAmount    float64
		expectedAmount float64
		expectedErr    error
	}{
		{"NormalPurchase positive input", NormalPurchase, 100.0, -100.0, nil},
		{"NormalPurchase negative input", NormalPurchase, -100.0, -100.0, nil},
		{"PurchaseInstallments positive input", PurchaseInstallments, 200.0, -200.0, nil},
		{"PurchaseInstallments negative input", PurchaseInstallments, -200.0, -200.0, nil},
		{"Withdrawal positive input", Withdrawal, 50.0, -50.0, nil},
		{"Withdrawal negative input", Withdrawal, -50.0, -50.0, nil},
		{"CreditVoucher positive input", CreditVoucher, 300.0, 300.0, nil},
		{"CreditVoucher negative input", CreditVoucher, -300.0, 300.0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeAmount(tt.operationType, tt.inputAmount)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedAmount, result)
		})
	}
}

func TestNormalizeAmount_ZeroAmount(t *testing.T) {
	_, err := NormalizeAmount(NormalPurchase, 0)
	assert.ErrorIs(t, err, exception.ErrAmountGreaterThanZero)
}

func TestNormalizeAmount_InvalidOperationType(t *testing.T) {
	_, err := NormalizeAmount(99, 100.0)
	assert.ErrorIs(t, err, exception.ErrInvalidOperationType)
}

func TestValidateOperationTypeForAmount_NoErr(t *testing.T) {
	tests := []struct {
		name   string
		op     int
		amount float64
	}{
		{"NormalPurchase", NormalPurchase, -10.0},
		{"PurchaseInstallments", PurchaseInstallments, -20.0},
		{"Withdrawal", Withdrawal, -30.0},
		{"CreditVoucher", CreditVoucher, 50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOperationTypeForAmount(tt.op, tt.amount)
			assert.NoError(t, err)
		})
	}
}

func TestValidateOperationTypeForAmount_AmountZero(t *testing.T) {
	err := ValidateOperationTypeForAmount(NormalPurchase, 0)
	assert.ErrorIs(t, err, exception.ErrAmountGreaterThanZero)
}

func TestValidateOperationTypeForAmount_InvalidOperationType(t *testing.T) {
	err := ValidateOperationTypeForAmount(99, 100.0)
	assert.ErrorIs(t, err, exception.ErrInvalidOperationType)
}
