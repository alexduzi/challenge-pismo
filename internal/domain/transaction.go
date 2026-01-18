package domain

import "time"

type Transaction struct {
	TransactionID   int       `db:"transaction_id"`
	AccountID       int       `db:"account_id"`
	OperationTypeID int       `db:"operation"`
	Amount          int64     `db:"amount"`
	EventDate       time.Time `db:"event_date"`
}
