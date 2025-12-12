package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kanyaarss/kanyaars-portal/internal/domain"
	"github.com/kanyaarss/kanyaars-portal/pkg/jwt"
)

// Auth returns a middleware that validates JWT tokens
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
				"Unauthorized",
				"Missing authorization header",
			))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
				"Unauthorized",
				"Invalid authorization header format",
			))
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := jwt.ValidateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
				"Unauthorized",
				"Invalid or expired token",
			))
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
