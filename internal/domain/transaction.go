package domain

import "time"

type Transaction struct {
	TransactionID   int64     `db:"transaction_id"`
	AccountID       int64     `db:"account_id"`
	OperationTypeID int       `db:"operation_type_id"`
	Amount          float64   `db:"amount"`
	Balance         float64   `db:"balance"`
	EventDate       time.Time `db:"event_date"`
	CreatedAt       time.Time `db:"created_at"`
}
