package domain

import "github.com/alexduzi/challengepismo/internal/infrastructure/exception"

const (
	NormalPurchase       = 1
	PurchaseInstallments = 2
	Withdrawal           = 3
	CreditVoucher        = 4
)

func ValidateOperationTypeForAmount(op int, amount float64) error {
	if amount == 0 {
		return exception.ErrAmountGreaterThanZero
	}

	switch op {
	case NormalPurchase, PurchaseInstallments, Withdrawal:
		if amount > 0 {
			return exception.ErrInvalidAmountForOperationType
		}
	case CreditVoucher:
		if amount < 0 {
			return exception.ErrInvalidAmountForOperationType
		}
	}
	return nil
}
