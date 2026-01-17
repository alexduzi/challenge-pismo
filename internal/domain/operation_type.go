package domain

import "fmt"

type OperationType int

const (
	NormalPurchase OperationType = iota
	PurchaseInstallments
	Withdrawal
	CreditVoucher
)

func GetOperationType(op OperationType) (OperationType, error) {
	switch op {
	case NormalPurchase:
		return NormalPurchase, nil
	case PurchaseInstallments:
		return PurchaseInstallments, nil
	case Withdrawal:
		return Withdrawal, nil
	case CreditVoucher:
		return CreditVoucher, nil
	default:
		return -1, fmt.Errorf("unknown operation type: %v", op)
	}
}
