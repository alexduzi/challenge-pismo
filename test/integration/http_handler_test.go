package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	conn "github.com/alexduzi/challengepismo/internal/infrastructure/persistence/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type HttpHandlerTestSuite struct {
	suite.Suite
	ctx               context.Context
	logger            logger.Logger
	postgresContainer *postgres.PostgresContainer
	db                *sqlx.DB
	engine            *gin.Engine
}

func (suite *HttpHandlerTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	conf, _ := config.LoadConfig()

	testServer, container := setupHttpHandler(suite.ctx, conf)

	suite.postgresContainer = container
	suite.db = testServer.DB
	suite.engine = testServer.Engine
	suite.logger = testServer.Logger
}

func (suite *HttpHandlerTestSuite) TearDownSuite() {
	if err := suite.postgresContainer.Terminate(suite.ctx); err != nil {
		suite.logger.Error("Error terminating postgres container: ", err)
	}
	if err := conn.CloseDB(suite.db); err != nil {
		suite.logger.Error("Error terminating sqlx connection: ", err)
	}
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_HealthCheck() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "healthy")
	assert.Contains(suite.T(), w.Body.String(), "api-accounts")
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_ReadinessCheck() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/readiness", nil)
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "healthy")
	assert.Contains(suite.T(), w.Body.String(), "api-accounts")
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CanCreateAccount() {
	reqBody := request.CreateAccountRequest{
		DocumentNumber: "12345678901",
		FullName:       "Test User",
		Email:          "testuser@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var resp response.AccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), reqBody.DocumentNumber, resp.DocumentNumber)
	assert.Equal(suite.T(), reqBody.FullName, resp.FullName)
	assert.Equal(suite.T(), reqBody.Email, resp.Email)
	assert.Equal(suite.T(), reqBody.Phone, resp.Phone)
	assert.Equal(suite.T(), reqBody.AccountType, resp.AccountType)
	assert.NotZero(suite.T(), resp.AccountID)
	assert.Contains(suite.T(), w.Header().Get("Location"), fmt.Sprintf("/api/v1/accounts/%d", resp.AccountID))
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CanGetAccountByID() {
	createReq := request.CreateAccountRequest{
		DocumentNumber: "98765432100",
		FullName:       "Get Test User",
		Email:          "gettest@example.com",
		Phone:          "11999998888",
		AccountType:    "savings",
	}

	body, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createResp response.AccountResponse
	_ = json.Unmarshal(w.Body.Bytes(), &createResp)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/accounts/%d", createResp.AccountID), nil)
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var getResp response.AccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &getResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createResp.AccountID, getResp.AccountID)
	assert.Equal(suite.T(), createReq.DocumentNumber, getResp.DocumentNumber)
	assert.Equal(suite.T(), createReq.FullName, getResp.FullName)
	assert.Equal(suite.T(), createReq.Email, getResp.Email)
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_GetAccountByID_NotFound() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accounts/999999", nil)
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var errResp response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "resource not found", errResp.Message)
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_GetAccountByID_InvalidID() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accounts/invalid", nil)
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var errResp response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "account_id must be a valid integer", errResp.Message)
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CreateAccount_ValidationError() {
	reqBody := request.CreateAccountRequest{
		DocumentNumber: "123",
		FullName:       "Te",
		Email:          "invalid-email",
		Phone:          "123",
		AccountType:    "invalid",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var errResp response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "documentnumber must be exactly 11 characters", errResp.Message)
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CreateAccount_DuplicateDocument() {
	reqBody := request.CreateAccountRequest{
		DocumentNumber: "11111111111",
		FullName:       "Duplicate Test",
		Email:          "duplicate1@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	reqBody.Email = "duplicate2@example.com" // Different email
	body, _ = json.Marshal(reqBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusConflict, w.Code)

	var errResp response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "document number already exists", errResp.Message)
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CanCreateTransaction() {
	createAccReq := request.CreateAccountRequest{
		DocumentNumber: "22222222222",
		FullName:       "Transaction Test User",
		Email:          "transactiontest@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	body, _ := json.Marshal(createAccReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var accResp response.AccountResponse
	_ = json.Unmarshal(w.Body.Bytes(), &accResp)

	tranReq := request.CreateTransactionRequest{
		AccountID:       accResp.AccountID,
		OperationTypeID: 4, // Payment
		Amount:          100.50,
	}

	body, _ = json.Marshal(tranReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var tranResp response.TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &tranResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), accResp.AccountID, tranResp.AccountID)
	assert.Equal(suite.T(), tranReq.OperationTypeID, tranResp.OperationTypeID)
	assert.NotZero(suite.T(), tranResp.TransactionID)
	assert.Contains(suite.T(), w.Header().Get("Location"), fmt.Sprintf("/api/v1/transactions/%d", tranResp.TransactionID))
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CreateTransaction_AccountNotFound() {
	tranReq := request.CreateTransactionRequest{
		AccountID:       999999, // Non-existent account
		OperationTypeID: 4,
		Amount:          100.50,
	}

	body, _ := json.Marshal(tranReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CreateTransaction_ValidationError() {
	tranReq := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 10, // Invalid - must be 1-4
		Amount:          -100,
	}

	body, _ := json.Marshal(tranReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CreateTransaction_Purchase() {
	createAccReq := request.CreateAccountRequest{
		DocumentNumber: "33333333333",
		FullName:       "Purchase Test User",
		Email:          "purchasetest@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	body, _ := json.Marshal(createAccReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var accResp response.AccountResponse
	_ = json.Unmarshal(w.Body.Bytes(), &accResp)

	tranReq := request.CreateTransactionRequest{
		AccountID:       accResp.AccountID,
		OperationTypeID: 1, // Purchase
		Amount:          50.00, // Positive input, should be normalized to negative
	}

	body, _ = json.Marshal(tranReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var tranResp response.TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &tranResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, tranResp.OperationTypeID)
	// Purchase should have negative amount (normalized from positive input)
	assert.Equal(suite.T(), -50.00, tranResp.Amount)
}

func (suite *HttpHandlerTestSuite) TestHttpHandler_CreateTransaction_Withdrawal() {
	createAccReq := request.CreateAccountRequest{
		DocumentNumber: "44444444444",
		FullName:       "Withdrawal Test User",
		Email:          "withdrawaltest@example.com",
		Phone:          "11987654321",
		AccountType:    "savings",
	}

	body, _ := json.Marshal(createAccReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var accResp response.AccountResponse
	_ = json.Unmarshal(w.Body.Bytes(), &accResp)

	tranReq := request.CreateTransactionRequest{
		AccountID:       accResp.AccountID,
		OperationTypeID: 3, // Withdrawal
		Amount:          75.00, // Positive input, should be normalized to negative
	}

	body, _ = json.Marshal(tranReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	suite.engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var tranResp response.TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &tranResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, tranResp.OperationTypeID)
	// Withdrawal should have negative amount (normalized from positive input)
	assert.Equal(suite.T(), -75.00, tranResp.Amount)
}

func TestHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HttpHandlerTestSuite))
}
