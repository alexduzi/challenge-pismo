package domain

import (
	"math"

	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
)

const (
	NormalPurchase       = 1
	PurchaseInstallments = 2
	Withdrawal           = 3
	CreditVoucher        = 4
)

func NormalizeAmount(operationType int, amount float64) (float64, error) {
	if amount == 0 {
		return 0, exception.ErrAmountGreaterThanZero
	}

	absAmount := math.Abs(amount)

	switch operationType {
	case NormalPurchase, PurchaseInstallments, Withdrawal:
		return -absAmount, nil
	case CreditVoucher:
		return absAmount, nil
	default:
		return 0, exception.ErrInvalidOperationType
	}
}

func ValidateOperationTypeForAmount(op int, amount float64) error {
	if amount == 0 {
		return exception.ErrAmountGreaterThanZero
	}

	switch op {
	case NormalPurchase, PurchaseInstallments, Withdrawal, CreditVoucher:
		return nil
	default:
		return exception.ErrInvalidOperationType
	}
}
