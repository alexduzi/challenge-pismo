package entity

type Transaction struct {
	TransactionID   int
	AccountID       int
	OperationTypeID int
	Amount          int64
	EventDate       string
}
