package handler

import (
	"fmt"
	"net/http"

	"github.com/alexduzi/challengepismo/internal/dto/request"
	_ "github.com/alexduzi/challengepismo/internal/dto/response" // swagger docs
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	customValidator "github.com/alexduzi/challengepismo/internal/infrastructure/validator"
	"github.com/alexduzi/challengepismo/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUseCase usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

// CreateTransaction creates a new transaction
// @Summary Create a new transaction
// @Description Creates a new transaction for an account. Operation types: 1=Purchase, 2=Installment Purchase, 3=Withdrawal, 4=Payment
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body request.CreateTransactionRequest true "Transaction creation request"
// @Success 201 {object} response.TransactionResponse "Transaction created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 404 {object} response.ErrorResponse "Account not found"
// @Failure 422 {object} response.ErrorResponse "Invalid amount for operation type"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req request.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(exception.NewValidationError("invalid request body"))
		return
	}

	if err := customValidator.ValidateStruct(&req); err != nil {
		_ = c.Error(exception.NewValidationError(err.Error()))
		return
	}

	tran, err := h.transactionUseCase.CreateTransaction(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	uri := fmt.Sprintf("http://%s/api/v1/transactions/%d", c.Request.Host, tran.TransactionID)
	c.Header("Location", uri)
	c.JSON(http.StatusCreated, tran)
}
