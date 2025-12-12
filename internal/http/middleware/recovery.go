package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kanyaarss/kanyaars-portal/internal/domain"
)

// Recovery returns a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
					"Internal Server Error",
					"An unexpected error occurred",
				))
				c.Abort()
			}
		}()

		c.Next()
	}
}
