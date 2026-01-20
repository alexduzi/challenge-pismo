package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoggerMiddleware_InfoOnSuccess(t *testing.T) {
	mockLog := logger.NewMockLogger(t)
	mockLog.EXPECT().Info("HTTP Request", mock.Anything).Once()

	router := setupTestRouter()
	router.Use(LoggerMiddleware(mockLog))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggerMiddleware_WarnOn4xx(t *testing.T) {
	mockLog := logger.NewMockLogger(t)
	mockLog.EXPECT().Warn("HTTP Request", mock.Anything).Once()

	router := setupTestRouter()
	router.Use(LoggerMiddleware(mockLog))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLoggerMiddleware_ErrorOn5xx(t *testing.T) {
	mockLog := logger.NewMockLogger(t)
	mockLog.EXPECT().Error("HTTP Request", mock.Anything).Once()

	router := setupTestRouter()
	router.Use(LoggerMiddleware(mockLog))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLoggerMiddleware_WithRequestID(t *testing.T) {
	mockLog := logger.NewMockLogger(t)
	mockLog.EXPECT().Info("HTTP Request", mock.Anything).Once()

	router := setupTestRouter()
	router.Use(RequestIDMiddleware())
	router.Use(LoggerMiddleware(mockLog))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggerMiddleware_WithQueryString(t *testing.T) {
	mockLog := logger.NewMockLogger(t)
	mockLog.EXPECT().Info("HTTP Request", mock.Anything).Once()

	router := setupTestRouter()
	router.Use(LoggerMiddleware(mockLog))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?foo=bar", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggerMiddleware_ErrorAppended(t *testing.T) {
	mockLog := logger.NewMockLogger(t)
	mockLog.EXPECT().Warn("HTTP Request", mock.Anything).Once()

	router := setupTestRouter()
	router.Use(LoggerMiddleware(mockLog))
	router.Use(ErrorHandlerMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.Error(exception.NewValidationError("invalid request body"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?foo=bar", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", response["message"])
}
