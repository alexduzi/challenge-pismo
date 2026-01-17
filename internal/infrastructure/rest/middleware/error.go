package middleware

import (
	"net/http"

	"github.com/alexduzi/challengepismo/internal/model"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 && !c.Writer.Written() {
			_ = c.Errors.Last().Err

			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Message: "internal server error",
			})
		}
	}
}
