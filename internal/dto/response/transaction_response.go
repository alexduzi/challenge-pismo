package response

import "time"

type TransactionResponse struct {
	TransactionID   int64     `json:"transaction_id"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int       `json:"operation"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
	CreatedAt       time.Time `json:"created_at"`
}
