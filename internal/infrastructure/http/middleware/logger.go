package middleware

import (
	"time"

	"log/slog"

	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		requestID, _ := c.Get("request_id")

		c.Next()

		latency := time.Since(start)

		status := c.Writer.Status()

		fields := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		if requestID != nil {
			fields = append(fields, slog.String("request_id", requestID.(string)))
		}

		if query != "" {
			fields = append(fields, slog.String("query", query))
		}

		if len(c.Errors) > 0 {
			fields = append(fields, slog.String("error", c.Errors.String()))
		}

		msg := "HTTP Request"
		if status >= 500 {
			log.Error(msg, fields...)
		} else if status >= 400 {
			log.Warn(msg, fields...)
		} else {
			log.Info(msg, fields...)
		}
	}
}
