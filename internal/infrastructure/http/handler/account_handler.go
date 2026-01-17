package handler

import (
	"net/http"
	"strconv"

	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/dto/response"
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
	var req request.Account
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	if err := h.accountUseCase.CreateAccount(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "account created successfully"})
}

func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "invalid account id",
		})
		return
	}

	account, err := h.accountUseCase.GetAccountByID(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "account not found",
		})
		return
	}

	c.JSON(http.StatusOK, account)
}
