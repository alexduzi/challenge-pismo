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

type AccountHandlerTestSuite struct {
	suite.Suite
	logger         *logger.MockLogger
	accountUseCase *usecase.MockAccountUseCase
	accountHandler *AccountHandler
	router         *gin.Engine
}

func (s *AccountHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.logger = logger.NewMockLogger(s.T())
	s.accountUseCase = usecase.NewMockAccountUseCase(s.T())
	s.accountHandler = NewAccountHandler(s.accountUseCase, s.logger)
	s.router = gin.New()
	s.router.POST("/accounts", s.accountHandler.CreateAccount)
	s.router.GET("/accounts/:accountId", s.accountHandler.GetAccountByID)
}

func (s *AccountHandlerTestSuite) TestCreateAccount_Success() {
	reqBody := request.CreateAccountRequest{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	expectedResponse := &response.AccountResponse{
		AccountID:      1,
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
		Balance:        0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Times(2)
	s.logger.EXPECT().Info("Creating account", mock.Anything).Once()
	s.logger.EXPECT().Info("Account created successfully", mock.Anything).Once()

	s.accountUseCase.EXPECT().CreateAccount(mock.Anything, reqBody).Return(expectedResponse, nil).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var resp response.AccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.AccountID, resp.AccountID)
	assert.Equal(s.T(), expectedResponse.DocumentNumber, resp.DocumentNumber)
	assert.Contains(s.T(), w.Header().Get("Location"), "/accounts/1")
}

func (s *AccountHandlerTestSuite) TestCreateAccount_InvalidRequestBody() {
	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Once()
	s.logger.EXPECT().Warn("Invalid request body", mock.Anything).Once()

	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *AccountHandlerTestSuite) TestCreateAccount_ValidationError() {
	reqBody := request.CreateAccountRequest{
		DocumentNumber: "123",
		FullName:       "Jo",
		Email:          "invalid-email",
		Phone:          "123",
		AccountType:    "invalid",
	}

	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Once()
	s.logger.EXPECT().Warn("Validation failed", mock.Anything).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *AccountHandlerTestSuite) TestCreateAccount_UseCaseError() {
	reqBody := request.CreateAccountRequest{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Times(2)
	s.logger.EXPECT().Info("Creating account", mock.Anything).Once()
	s.logger.EXPECT().Error("Failed to create account", mock.Anything).Once()

	s.accountUseCase.EXPECT().CreateAccount(mock.Anything, reqBody).Return(nil, exception.ErrDatabaseError).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *AccountHandlerTestSuite) TestCreateAccount_DuplicateDocumentError() {
	reqBody := request.CreateAccountRequest{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Times(2)
	s.logger.EXPECT().Info("Creating account", mock.Anything).Once()
	s.logger.EXPECT().Error("Failed to create account", mock.Anything).Once()

	s.accountUseCase.EXPECT().CreateAccount(mock.Anything, reqBody).Return(nil, exception.ErrDuplicateDocument).Once()

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *AccountHandlerTestSuite) TestGetAccountByID_Success() {
	accountID := int64(1)

	expectedResponse := &response.AccountResponse{
		AccountID:      accountID,
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
		Balance:        1000.50,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Once()
	s.logger.EXPECT().Debug("Fetching account", mock.Anything).Once()

	s.accountUseCase.EXPECT().GetAccountByID(mock.Anything, accountID).Return(expectedResponse, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resp response.AccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.AccountID, resp.AccountID)
	assert.Equal(s.T(), expectedResponse.DocumentNumber, resp.DocumentNumber)
}

func (s *AccountHandlerTestSuite) TestGetAccountByID_InvalidID() {
	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Once()
	s.logger.EXPECT().Warn("Invalid account ID", mock.Anything).Once()

	req, _ := http.NewRequest(http.MethodGet, "/accounts/invalid", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *AccountHandlerTestSuite) TestGetAccountByID_NotFound() {
	accountID := int64(999)

	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Times(2)
	s.logger.EXPECT().Debug("Fetching account", mock.Anything).Once()
	s.logger.EXPECT().Error("Failed to fetch account", mock.Anything).Once()

	s.accountUseCase.EXPECT().GetAccountByID(mock.Anything, accountID).Return(nil, exception.ErrNotFound).Once()

	req, _ := http.NewRequest(http.MethodGet, "/accounts/999", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *AccountHandlerTestSuite) TestGetAccountByID_DatabaseError() {
	accountID := int64(1)

	s.logger.EXPECT().WithContext(mock.Anything).Return(s.logger).Times(2)
	s.logger.EXPECT().Debug("Fetching account", mock.Anything).Once()
	s.logger.EXPECT().Error("Failed to fetch account", mock.Anything).Once()

	s.accountUseCase.EXPECT().GetAccountByID(mock.Anything, accountID).Return(nil, exception.ErrDatabaseError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func TestAccountHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AccountHandlerTestSuite))
}
