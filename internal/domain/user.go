package domain

import "time"

// User represents an admin user
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"` // Never expose password in JSON
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserLoginRequest represents login request
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserLoginResponse represents login response
type UserLoginResponse struct {
	Token     string    `json:"token"`
	ExpiresIn int64     `json:"expires_in"`
	User      *UserInfo `json:"user"`
}

// UserInfo represents user info (safe to expose)
type UserInfo struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
