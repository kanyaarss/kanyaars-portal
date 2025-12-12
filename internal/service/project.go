package service

import (
	"database/sql"
	"fmt"

	"github.com/kanyaarss/kanyaars-portal/internal/domain"
)

// ProjectService handles project operations
type ProjectService struct {
	db *sql.DB
}

// NewProjectService creates a new project service
func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{db: db}
}

// GetAllProjects retrieves all projects
func (s *ProjectService) GetAllProjects() ([]domain.Project, error) {
	rows, err := s.db.Query(
		"SELECT id, name, slug, description, url, icon_url, status, \"order\", created_at, updated_at FROM projects ORDER BY \"order\" ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.URL, &p.IconURL, &p.Status, &p.Order, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetActiveProjects retrieves only active projects
func (s *ProjectService) GetActiveProjects() ([]domain.Project, error) {
	rows, err := s.db.Query(
		"SELECT id, name, slug, description, url, icon_url, status, \"order\", created_at, updated_at FROM projects WHERE status = 'active' ORDER BY \"order\" ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.URL, &p.IconURL, &p.Status, &p.Order, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetProjectByID retrieves a project by ID
func (s *ProjectService) GetProjectByID(id int) (*domain.Project, error) {
	var p domain.Project
	err := s.db.QueryRow(
		"SELECT id, name, slug, description, url, icon_url, status, \"order\", created_at, updated_at FROM projects WHERE id = $1",
		id,
	).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.URL, &p.IconURL, &p.Status, &p.Order, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &p, nil
}

// GetProjectBySlug retrieves a project by slug
func (s *ProjectService) GetProjectBySlug(slug string) (*domain.Project, error) {
	var p domain.Project
	err := s.db.QueryRow(
		"SELECT id, name, slug, description, url, icon_url, status, \"order\", created_at, updated_at FROM projects WHERE slug = $1",
		slug,
	).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.URL, &p.IconURL, &p.Status, &p.Order, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &p, nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(req *domain.CreateProjectRequest) (int, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO projects (name, slug, description, url, icon_url, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Name, req.Slug, req.Description, req.URL, req.IconURL, req.Status,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(id int, req *domain.UpdateProjectRequest) error {
	_, err := s.db.Exec(
		"UPDATE projects SET name = COALESCE(NULLIF($1, ''), name), slug = COALESCE(NULLIF($2, ''), slug), description = COALESCE(NULLIF($3, ''), description), url = COALESCE(NULLIF($4, ''), url), icon_url = COALESCE(NULLIF($5, ''), icon_url), status = COALESCE(NULLIF($6, ''), status), updated_at = CURRENT_TIMESTAMP WHERE id = $7",
		req.Name, req.Slug, req.Description, req.URL, req.IconURL, req.Status, id,
	)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id int) error {
	result, err := s.db.Exec("DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}
