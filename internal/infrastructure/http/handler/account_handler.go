package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexduzi/challengepismo/internal/dto/request"
	_ "github.com/alexduzi/challengepismo/internal/dto/response" // swagger docs
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

// CreateAccount creates a new account
// @Summary Create a new account
// @Description Creates a new account with the provided details
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body request.CreateAccountRequest true "Account creation request"
// @Success 201 {object} response.AccountResponse "Account created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 409 {object} response.ErrorResponse "Document number or email already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/accounts [post]
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

// GetAccountByID retrieves an account by ID
// @Summary Get account by ID
// @Description Retrieves account details by account ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param accountId path int true "Account ID"
// @Success 200 {object} response.AccountResponse "Account found"
// @Failure 400 {object} response.ErrorResponse "Invalid account ID"
// @Failure 404 {object} response.ErrorResponse "Account not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/accounts/{accountId} [get]
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
