package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware_GeneratesID(t *testing.T) {
	router := setupTestRouter()
	router.Use(RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		assert.NotEmpty(t, requestID)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}

func TestRequestIDMiddleware_UsesExistingID(t *testing.T) {
	router := setupTestRouter()
	router.Use(RequestIDMiddleware())

	expectedID := "custom-request-id-123"

	router.GET("/test", func(c *gin.Context) {
		requestID, _ := c.Get("request_id")
		assert.Equal(t, expectedID, requestID)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", expectedID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedID, w.Header().Get("X-Request-ID"))
}

func TestRequestIDMiddleware_PropagatesInResponse(t *testing.T) {
	router := setupTestRouter()
	router.Use(RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	responseID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, responseID)
	assert.Len(t, responseID, 36)
}
