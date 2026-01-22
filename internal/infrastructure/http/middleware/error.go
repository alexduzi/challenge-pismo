package middleware

import (
	"errors"
	"net/http"

	"github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		status, message := mapError(err)

		c.JSON(status, response.ErrorResponse{
			Message: message,
		})
	}
}

func mapError(err error) (int, string) {
	var validationErr *exception.ValidationError
	if errors.As(err, &validationErr) {
		return http.StatusBadRequest, validationErr.Message
	}

	switch {
	case errors.Is(err, exception.ErrAccountNotExistsForTransaction):
		return http.StatusNotFound, "invalid account for transaction"
	case errors.Is(err, exception.ErrAmountGreaterThanZero):
		return http.StatusBadRequest, "amount must be greater than 0"
	case errors.Is(err, exception.ErrInvalidOperationType):
		return http.StatusBadRequest, "invalid operation type"
	case errors.Is(err, exception.ErrDuplicateDocument):
		return http.StatusConflict, "document number already exists"
	case errors.Is(err, exception.ErrDuplicateEmail):
		return http.StatusConflict, "email already exists"
	case errors.Is(err, exception.ErrNotFound):
		return http.StatusNotFound, "resource not found"
	case errors.Is(err, exception.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"
	case errors.Is(err, exception.ErrDatabaseError):
		return http.StatusInternalServerError, "internal server error"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
