package response

import "time"

type TransactionResponse struct {
	TransactionID   int64     `json:"transaction_id" example:"1"`
	AccountID       int64     `json:"account_id" example:"1"`
	OperationTypeID int       `json:"operation" example:"4"`
	Amount          float64   `json:"amount" example:"100.50"`
	EventDate       time.Time `json:"event_date" example:"2024-01-01T10:00:00Z"`
	CreatedAt       time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`
}
