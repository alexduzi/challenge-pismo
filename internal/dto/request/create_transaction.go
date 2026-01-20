package request

type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id" validate:"required,numeric" example:"1"`
	OperationTypeID int     `json:"operation_type_id" validate:"required,numeric,gte=1,lte=4" example:"4"`
	Amount          float64 `json:"amount" validate:"required,numeric" example:"100.50"`
}
