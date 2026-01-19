package request

type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id" validate:"required,numeric"`
	OperationTypeID int     `json:"operation" validate:"required,numeric,gte=1,lte=4"`
	Amount          float64 `json:"amount" validate:"required,numeric"`
}
