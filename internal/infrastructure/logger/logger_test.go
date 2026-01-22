package logger

import (
	"context"
	"testing"

	"github.com/alexduzi/challengepismo/internal/constants"
	"github.com/alexduzi/challengepismo/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestLogger_CanCreateInstance(t *testing.T) {
	cfg, _ := config.LoadConfig()

	logger := NewLogger(cfg)

	assert.NotNil(t, logger)
}

func TestLogger_WithDifferentLogLevels(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
	}{
		{"debug level", "debug"},
		{"info level", "info"},
		{"warn level", "warn"},
		{"error level", "error"},
		{"default level", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				LogLevel:  tt.logLevel,
				LogFormat: "text",
				AppName:   "test-app",
				GinMode:   "test",
			}

			logger := NewLogger(cfg)

			assert.NotNil(t, logger)
		})
	}
}

func TestLogger_WithJSONFormat(t *testing.T) {
	cfg := &config.Config{
		LogLevel:  "info",
		LogFormat: "json",
		AppName:   "test-app",
		GinMode:   "test",
	}

	logger := NewLogger(cfg)

	assert.NotNil(t, logger)
}

func TestLogger_LogMethods(t *testing.T) {
	cfg := &config.Config{
		LogLevel:  "debug",
		LogFormat: "text",
		AppName:   "test-app",
		GinMode:   "test",
	}

	logger := NewLogger(cfg)

	// These should not panic
	assert.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")
	})
}

func TestLogger_LogMethodsWithArgs(t *testing.T) {
	cfg := &config.Config{
		LogLevel:  "debug",
		LogFormat: "text",
		AppName:   "test-app",
		GinMode:   "test",
	}

	logger := NewLogger(cfg)

	assert.NotPanics(t, func() {
		logger.Debug("debug message", "key", "value")
		logger.Info("info message", "count", 42)
		logger.Warn("warn message", "flag", true)
		logger.Error("error message", "error", "something went wrong")
	})
}

func TestLogger_With(t *testing.T) {
	cfg := &config.Config{
		LogLevel:  "info",
		LogFormat: "text",
		AppName:   "test-app",
		GinMode:   "test",
	}

	logger := NewLogger(cfg)
	newLogger := logger.With("user_id", "123")

	assert.NotNil(t, newLogger)
	assert.NotSame(t, logger, newLogger)
}

func TestLogger_WithContext_WithRequestID(t *testing.T) {
	cfg := &config.Config{
		LogLevel:  "info",
		LogFormat: "text",
		AppName:   "test-app",
		GinMode:   "test",
	}
	requestID := "req-123"

	logger := NewLogger(cfg)
	ctx := context.WithValue(context.Background(), constants.RequestIDKey, requestID)
	newLogger := logger.WithContext(ctx)

	assert.NotNil(t, newLogger)
	assert.NotSame(t, logger, newLogger)
}

func TestLogger_WithContext_WithoutRequestID(t *testing.T) {
	cfg := &config.Config{
		LogLevel:  "info",
		LogFormat: "text",
		AppName:   "test-app",
		GinMode:   "test",
	}

	logger := NewLogger(cfg)
	ctx := context.Background()
	newLogger := logger.WithContext(ctx)

	assert.NotNil(t, newLogger)
}

func TestSlogLogger_GetLogger(t *testing.T) {
	cfg := &config.Config{
		LogLevel:  "info",
		LogFormat: "text",
		AppName:   "test-app",
		GinMode:   "test",
	}

	logger := NewLogger(cfg)
	slogLogger, ok := logger.(*SlogLogger)

	assert.True(t, ok)
	assert.NotNil(t, slogLogger.GetLogger())
}
