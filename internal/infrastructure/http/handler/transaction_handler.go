package handler

import (
	"fmt"
	"net/http"

	"github.com/alexduzi/challengepismo/internal/dto/request"
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
