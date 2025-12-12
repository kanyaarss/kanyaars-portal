package service

import (
	"database/sql"
	"fmt"

	"github.com/kanyaarss/kanyaars-portal/internal/domain"
)

// PortalService handles portal configuration operations
type PortalService struct {
	db *sql.DB
}

// NewPortalService creates a new portal service
func NewPortalService(db *sql.DB) *PortalService {
	return &PortalService{db: db}
}

// GetPortal retrieves portal configuration
func (s *PortalService) GetPortal() (*domain.Portal, error) {
	var portal domain.Portal

	err := s.db.QueryRow(
		"SELECT id, name, description, logo_url, website, email, phone, address, created_at, updated_at FROM portal_config LIMIT 1",
	).Scan(
		&portal.ID,
		&portal.Name,
		&portal.Description,
		&portal.LogoURL,
		&portal.Website,
		&portal.Email,
		&portal.Phone,
		&portal.Address,
		&portal.CreatedAt,
		&portal.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("portal not configured")
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &portal, nil
}

// UpdatePortal updates portal configuration
func (s *PortalService) UpdatePortal(req *domain.UpdatePortalRequest) error {
	// Check if portal config exists
	var id int
	err := s.db.QueryRow("SELECT id FROM portal_config LIMIT 1").Scan(&id)

	if err == sql.ErrNoRows {
		// Create new portal config
		_, err := s.db.Exec(
			"INSERT INTO portal_config (name, description, logo_url, website, email, phone, address) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			req.Name, req.Description, req.LogoURL, req.Website, req.Email, req.Phone, req.Address,
		)
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}
	} else if err == nil {
		// Update existing portal config
		_, err := s.db.Exec(
			"UPDATE portal_config SET name = COALESCE(NULLIF($1, ''), name), description = COALESCE(NULLIF($2, ''), description), logo_url = COALESCE(NULLIF($3, ''), logo_url), website = COALESCE(NULLIF($4, ''), website), email = COALESCE(NULLIF($5, ''), email), phone = COALESCE(NULLIF($6, ''), phone), address = COALESCE(NULLIF($7, ''), address), updated_at = CURRENT_TIMESTAMP WHERE id = $8",
			req.Name, req.Description, req.LogoURL, req.Website, req.Email, req.Phone, req.Address, id,
		)
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}
	} else {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}
