package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kanyaarss/kanyaars-portal/internal/domain"
	"github.com/kanyaarss/kanyaars-portal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *sql.DB
	jwtSecret string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *sql.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Login handles admin login
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Invalid request",
			err.Error(),
		))
		return
	}

	// Find user by email
	var user domain.User
	err := h.db.QueryRow(
		"SELECT id, email, name, password, role FROM users WHERE email = $1 AND is_active = true",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
			"Login failed",
			"Invalid email or password",
		))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
			"Login failed",
			"Invalid email or password",
		))
		return
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, h.jwtSecret, 86400)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Token generation failed",
			err.Error(),
		))
		return
	}

	// Return response
	c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Login successful", domain.UserLoginResponse{
		Token:     token,
		ExpiresIn: 86400,
		User: &domain.UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}))
}
