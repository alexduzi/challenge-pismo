package response

type TransactionResponse struct {
	AccountID       int   `json:"account_id"`
	OperationTypeID int   `json:"operation"`
	Amount          int64 `json:"amount"`
}
