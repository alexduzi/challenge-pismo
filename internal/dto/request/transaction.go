package dto

type Transaction struct {
	AccountID       int    `json:"account_id"`
	OperationTypeID int    `json:"operation"`
	Amount          int64  `json:"amount"`
	EventDate       string `json:"event_date"`
}
