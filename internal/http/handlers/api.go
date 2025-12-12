package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanyaarss/kanyaars-portal/internal/domain"
)

type APIHandler struct {
	db *sql.DB
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(db *sql.DB) *APIHandler {
	return &APIHandler{db: db}
}

// HealthCheck returns the health status of the API
func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, domain.HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
		Uptime:    int64(time.Since(time.Now()).Seconds()),
	})
}

// GetPortal returns portal information
func (h *APIHandler) GetPortal(c *gin.Context) {
	var portal domain.Portal

	err := h.db.QueryRow(
		"SELECT id, name, description, logo_url, website, email, phone, address FROM portal_config LIMIT 1",
	).Scan(
		&portal.ID,
		&portal.Name,
		&portal.Description,
		&portal.LogoURL,
		&portal.Website,
		&portal.Email,
		&portal.Phone,
		&portal.Address,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Portal not configured", nil))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Portal retrieved", portal))
}

// GetProjects returns all active projects
func (h *APIHandler) GetProjects(c *gin.Context) {
	rows, err := h.db.Query(
		"SELECT id, name, slug, description, url, icon_url, status, created_at FROM projects WHERE status = 'active' ORDER BY \"order\" ASC",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.URL, &p.IconURL, &p.Status, &p.CreatedAt); err != nil {
			continue
		}
		projects = append(projects, p)
	}

	c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Projects retrieved", gin.H{
		"data":  projects,
		"total": len(projects),
	}))
}

// GetProject returns a single project by ID
func (h *APIHandler) GetProject(c *gin.Context) {
	id := c.Param("id")

	var project domain.Project
	err := h.db.QueryRow(
		"SELECT id, name, slug, description, url, icon_url, status, created_at FROM projects WHERE id = $1",
		id,
	).Scan(&project.ID, &project.Name, &project.Slug, &project.Description, &project.URL, &project.IconURL, &project.Status, &project.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(
			"Not found",
			"Project not found",
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

	c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Project retrieved", project))
}
