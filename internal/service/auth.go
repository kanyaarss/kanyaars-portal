package service

import (
	"database/sql"
	"fmt"

	"github.com/kanyaarss/kanyaars-portal/internal/domain"
	"github.com/kanyaarss/kanyaars-portal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	db        *sql.DB
	jwtSecret string
	jwtExpiry int64
}

// NewAuthService creates a new auth service
func NewAuthService(db *sql.DB, jwtSecret string, jwtExpiry int64) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Login authenticates user and returns JWT token
func (s *AuthService) Login(email, password string) (*domain.User, string, error) {
	// Find user by email
	var user domain.User
	err := s.db.QueryRow(
		"SELECT id, email, name, password, role FROM users WHERE email = $1 AND is_active = true",
		email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	if err != nil {
		return nil, "", fmt.Errorf("database error: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, "", fmt.Errorf("token generation failed: %w", err)
	}

	return &user, token, nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(id int) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRow(
		"SELECT id, email, name, role, is_active FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.IsActive)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(user *domain.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}

	_, err = s.db.Exec(
		"INSERT INTO users (email, name, password, role, is_active) VALUES ($1, $2, $3, $4, $5)",
		user.Email, user.Name, string(hashedPassword), user.Role, user.IsActive,
	)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

// UpdatePassword updates user password
func (s *AuthService) UpdatePassword(userID int, newPassword string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}

	_, err = s.db.Exec(
		"UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		string(hashedPassword), userID,
	)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}
