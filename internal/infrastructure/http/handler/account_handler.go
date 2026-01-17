package handler

import (
	"net/http"

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
	c.JSON(http.StatusCreated, gin.H{"message": "account created successfully"})
}

func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
