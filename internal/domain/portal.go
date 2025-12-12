package domain

import "time"

// Portal represents portal configuration
type Portal struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	Website     string    `json:"website"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdatePortalRequest represents update portal request
type UpdatePortalRequest struct {
	Name        string `json:"name" binding:"min=3"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url" binding:"url"`
	Website     string `json:"website" binding:"url"`
	Email       string `json:"email" binding:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
}
