package middleware

import (
	"github.com/alexduzi/challengepismo/internal/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(constants.XRequestID)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(string(constants.RequestIDKey), requestID)
		c.Header(constants.XRequestID, requestID)

		c.Next()
	}
}
