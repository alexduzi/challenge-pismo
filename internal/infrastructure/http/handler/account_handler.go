package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	customValidator "github.com/alexduzi/challengepismo/internal/infrastructure/validator"
	"github.com/alexduzi/challengepismo/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountUseCase usecase.AccountUseCase
}

func NewAccountHandler(accountUseCase usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{
		accountUseCase: accountUseCase,
	}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req request.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(exception.NewValidationError("invalid request body"))
		return
	}

	if err := customValidator.ValidateStruct(&req); err != nil {
		_ = c.Error(exception.NewValidationError(err.Error()))
		return
	}

	acc, err := h.accountUseCase.CreateAccount(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	uri := fmt.Sprintf("http://%s/api/v1/accounts/%d", c.Request.Host, acc.AccountID)
	c.Header("Location", uri)

	c.JSON(http.StatusCreated, acc)
}

func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		_ = c.Error(exception.NewValidationError("account_id must be a valid integer"))
		return
	}

	account, err := h.accountUseCase.GetAccountByID(c.Request.Context(), accountID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}
