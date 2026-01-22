package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HealthHandlerTestSuite struct {
	suite.Suite
	cfg           *config.Config
	healthHandler *HealthHandler
	router        *gin.Engine
}

func (s *HealthHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.cfg = &config.Config{
		AppName: "test-api",
		Port:    "8080",
	}
	s.healthHandler = NewHealthHandler(s.cfg)
	s.router = gin.New()
	s.router.GET("/health", s.healthHandler.GetStatus)
	s.router.GET("/readiness", s.healthHandler.GetStatus)
}

func (s *HealthHandlerTestSuite) TestGetStatus_Success() {
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resp response.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "healthy", resp.Status)
	assert.Equal(s.T(), "test-api", resp.Service)
	assert.NotZero(s.T(), resp.Timestamp)
}

func (s *HealthHandlerTestSuite) TestGetStatus_Readiness() {
	req, _ := http.NewRequest(http.MethodGet, "/readiness", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resp response.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "healthy", resp.Status)
	assert.Equal(s.T(), "test-api", resp.Service)
}

func (s *HealthHandlerTestSuite) TestGetStatus_ReturnsCorrectServiceName() {
	s.cfg.AppName = "custom-service-name"
	s.healthHandler = NewHealthHandler(s.cfg)
	s.router = gin.New()
	s.router.GET("/health", s.healthHandler.GetStatus)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resp response.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "custom-service-name", resp.Service)
}

func (s *HealthHandlerTestSuite) TestGetStatus_TimestampIsRecent() {
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resp response.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(s.T(), err)

	// Timestamp should be within the last few seconds
	assert.WithinDuration(s.T(), resp.Timestamp, resp.Timestamp, 5)
}

func TestHealthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HealthHandlerTestSuite))
}
