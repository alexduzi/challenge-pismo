package domain

type Transaction struct {
	TransactionID   int    `db:"transaction_id"`
	AccountID       int    `db:"account_id"`
	OperationTypeID int    `db:"operation"`
	Amount          int64  `db:"amount"`
	EventDate       string `db:"event_date"`
}
