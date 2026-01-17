package middleware

import (
	"net/http"

	dto "github.com/alexduzi/challengepismo/internal/dto/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 && !c.Writer.Written() {
			_ = c.Errors.Last().Err

			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "internal server error",
			})
		}
	}
}
