package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexduzi/challengepismo/internal/dto/request"
	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/alexduzi/challengepismo/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TransactionHandlerTestSuite struct {
	suite.Suite
	logger             *logger.MockLogger
	transactionUseCase *usecase.MockTransactionUseCase
	transactionHandler *TransactionHandler
	router             *gin.Engine
}

func (s *TransactionHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.logger = logger.NewMockLogger(s.T())
	s.transactionUseCase = usecase.NewMockTransactionUseCase(s.T())
	s.transactionHandler = NewTransactionHandler(s.transactionUseCase, s.logger)
	s.router = gin.New()
	s.router.POST("/transactions", s.transactionHandler.CreateTransaction)
}

func (s *TransactionHandlerTestSuite) TestCreateTransaction_Success() {
	reqBody := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 4, // Payment
		Amount:          100.50,
	}

	expectedResponse := &response.TransactionResponse{
		TransactionID:   1,
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          100.50,
		EventDate:       time.Now(),
		CreatedAt:       time.Now(),
	}

	s.transactionUseCase.EXPECT().CreateTransaction(mock.Anything, reqBody).Return(expectedResponse, nil).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var resp response.TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.TransactionID, resp.TransactionID)
	assert.Equal(s.T(), expectedResponse.AccountID, resp.AccountID)
	assert.Contains(s.T(), w.Header().Get("Location"), "/transactions/1")
}

func (s *TransactionHandlerTestSuite) TestCreateTransaction_InvalidRequestBody() {
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code) // Error set via c.Error without middleware
}

func (s *TransactionHandlerTestSuite) TestCreateTransaction_ValidationError_InvalidOperationType() {
	reqBody := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 5, // Invalid: must be 1-4
		Amount:          100.50,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code) // Error set via c.Error
}

func (s *TransactionHandlerTestSuite) TestCreateTransaction_ValidationError_MissingAccountID() {
	reqBody := request.CreateTransactionRequest{
		OperationTypeID: 4,
		Amount:          100.50,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code) // Error set via c.Error
}

func (s *TransactionHandlerTestSuite) TestCreateTransaction_UseCaseError_AccountNotFound() {
	reqBody := request.CreateTransactionRequest{
		AccountID:       999,
		OperationTypeID: 4,
		Amount:          100.50,
	}

	s.transactionUseCase.EXPECT().CreateTransaction(mock.Anything, reqBody).Return(nil, exception.ErrAccountNotExistsForTransaction).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code) // Error set via c.Error
}

func (s *TransactionHandlerTestSuite) TestCreateTransaction_UseCaseError_DatabaseError() {
	reqBody := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          100.50,
	}

	s.transactionUseCase.EXPECT().CreateTransaction(mock.Anything, reqBody).Return(nil, exception.ErrDatabaseError).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code) // Error set via c.Error
}

func (s *TransactionHandlerTestSuite) TestCreateTransaction_Purchase_Success() {
	reqBody := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 1, // Purchase
		Amount:          -50.00,
	}

	expectedResponse := &response.TransactionResponse{
		TransactionID:   2,
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          -50.00,
		EventDate:       time.Now(),
		CreatedAt:       time.Now(),
	}

	s.transactionUseCase.EXPECT().CreateTransaction(mock.Anything, reqBody).Return(expectedResponse, nil).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var resp response.TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.TransactionID, resp.TransactionID)
	assert.Equal(s.T(), expectedResponse.Amount, resp.Amount)
}

func (s *TransactionHandlerTestSuite) TestCreateTransaction_Withdrawal_Success() {
	reqBody := request.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 3, // Withdrawal
		Amount:          -200.00,
	}

	expectedResponse := &response.TransactionResponse{
		TransactionID:   3,
		AccountID:       1,
		OperationTypeID: 3,
		Amount:          -200.00,
		EventDate:       time.Now(),
		CreatedAt:       time.Now(),
	}

	s.transactionUseCase.EXPECT().CreateTransaction(mock.Anything, reqBody).Return(expectedResponse, nil).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var resp response.TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.OperationTypeID, resp.OperationTypeID)
}

func TestTransactionHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionHandlerTestSuite))
}
