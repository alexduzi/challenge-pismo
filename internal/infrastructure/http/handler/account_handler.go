package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	customValidator "github.com/alexduzi/challengepismo/internal/infrastructure/validator"
	"github.com/alexduzi/challengepismo/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountUseCase usecase.AccountUseCase
	logger         logger.Logger
}

func NewAccountHandler(accountUseCase usecase.AccountUseCase, logger logger.Logger) *AccountHandler {
	return &AccountHandler{
		accountUseCase: accountUseCase,
		logger:         logger,
	}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req request.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext(c.Request.Context()).Warn("Invalid request body",
			slog.String("error", err.Error()),
		)
		_ = c.Error(exception.NewValidationError("invalid request body"))
		return
	}

	if err := customValidator.ValidateStruct(&req); err != nil {
		h.logger.WithContext(c.Request.Context()).Warn("Validation failed",
			slog.String("error", err.Error()),
		)
		_ = c.Error(exception.NewValidationError(err.Error()))
		return
	}

	h.logger.WithContext(c.Request.Context()).Info("Creating account",
		slog.String("document_number", req.DocumentNumber),
		slog.String("email", req.Email),
	)

	acc, err := h.accountUseCase.CreateAccount(c.Request.Context(), req)
	if err != nil {
		h.logger.WithContext(c.Request.Context()).Error("Failed to create account",
			slog.String("error", err.Error()),
			slog.String("document_number", req.DocumentNumber),
		)
		_ = c.Error(err)
		return
	}

	h.logger.WithContext(c.Request.Context()).Info("Account created successfully",
		slog.Int64("account_id", acc.AccountID),
		slog.String("document_number", acc.DocumentNumber),
	)

	uri := fmt.Sprintf("http://%s/api/v1/accounts/%d", c.Request.Host, acc.AccountID)
	c.Header("Location", uri)
	c.JSON(http.StatusCreated, acc)
}

func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		h.logger.WithContext(c.Request.Context()).Warn("Invalid account ID",
			slog.String("account_id", c.Param("accountId")),
		)
		_ = c.Error(exception.NewValidationError("account_id must be a valid integer"))
		return
	}

	h.logger.WithContext(c.Request.Context()).Debug("Fetching account",
		slog.Int("account_id", accountID),
	)

	account, err := h.accountUseCase.GetAccountByID(c.Request.Context(), int64(accountID))
	if err != nil {
		h.logger.WithContext(c.Request.Context()).Error("Failed to fetch account",
			slog.String("error", err.Error()),
			slog.Int("account_id", accountID),
		)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}
