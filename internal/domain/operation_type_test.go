package domain

import (
	"testing"

	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/stretchr/testify/assert"
)

func TestValidateOperationTypeForAmount_NoErr(t *testing.T) {
	var tests = []struct {
		name   string
		op     int
		amount float64
	}{
		{"NormalPurchase", 1, -10.0},
		{"PurchaseInstallments", 2, -20.0},
		{"Withdrawal", 3, -30.0},
		{"CreditVoucher", 4, 50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOperationTypeForAmount(tt.op, tt.amount)
			if err != nil {
				t.Errorf("ValidateOperationTypeForAmount(%d, %f) = no error; got %v", tt.op, tt.amount, err)
			}
		})
	}
}

func TestValidateOperationTypeForAmount_WithErr(t *testing.T) {
	var tests = []struct {
		name   string
		op     int
		amount float64
	}{
		{"NormalPurchase", 1, 10.0},
		{"PurchaseInstallments", 2, 20.0},
		{"Withdrawal", 3, 30.0},
		{"CreditVoucher", 4, -50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOperationTypeForAmount(tt.op, tt.amount)
			assert.ErrorIs(t, err, exception.ErrInvalidAmountForOperationType, "ErrInvalidAmountForOperationType")
		})
	}
}

func TestValidateOperationTypeForAmount_AmountZero(t *testing.T) {
	err := ValidateOperationTypeForAmount(1, 0)
	assert.ErrorIs(t, err, exception.ErrAmountGreaterThanZero, "ErrAmountGreaterThanZero")
}
