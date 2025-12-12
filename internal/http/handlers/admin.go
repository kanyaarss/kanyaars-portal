package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kanyaarss/kanyaars-portal/internal/domain"
)

type AdminHandler struct {
	db *sql.DB
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// Dashboard renders the admin dashboard
func (h *AdminHandler) Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/dashboard.html", gin.H{
		"title": "Admin Dashboard",
	})
}

// ListProjects returns all projects
func (h *AdminHandler) ListProjects(c *gin.Context) {
	rows, err := h.db.Query(
		"SELECT id, name, slug, description, url, icon_url, status, \"order\", created_at FROM projects ORDER BY \"order\" ASC",
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
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.URL, &p.IconURL, &p.Status, &p.Order, &p.CreatedAt); err != nil {
			continue
		}
		projects = append(projects, p)
	}

	c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Projects retrieved", projects))
}

// CreateProject creates a new project
func (h *AdminHandler) CreateProject(c *gin.Context) {
	var req domain.CreateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Invalid request",
			err.Error(),
		))
		return
	}

	var id int
	err := h.db.QueryRow(
		"INSERT INTO projects (name, slug, description, url, icon_url, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Name, req.Slug, req.Description, req.URL, req.IconURL, req.Status,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, domain.NewAPIResponse(true, "Project created", gin.H{"id": id}))
}

// GetProject returns a single project
func (h *AdminHandler) GetProject(c *gin.Context) {
	id := c.Param("id")

	var project domain.Project
	err := h.db.QueryRow(
		"SELECT id, name, slug, description, url, icon_url, status, \"order\", created_at FROM projects WHERE id = $1",
		id,
	).Scan(&project.ID, &project.Name, &project.Slug, &project.Description, &project.URL, &project.IconURL, &project.Status, &project.Order, &project.CreatedAt)

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

// UpdateProject updates a project
func (h *AdminHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Invalid request",
			err.Error(),
		))
		return
	}

	_, err := h.db.Exec(
		"UPDATE projects SET name = COALESCE(NULLIF($1, ''), name), slug = COALESCE(NULLIF($2, ''), slug), description = COALESCE(NULLIF($3, ''), description), url = COALESCE(NULLIF($4, ''), url), icon_url = COALESCE(NULLIF($5, ''), icon_url), status = COALESCE(NULLIF($6, ''), status), updated_at = CURRENT_TIMESTAMP WHERE id = $7",
		req.Name, req.Slug, req.Description, req.URL, req.IconURL, req.Status, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Project updated", nil))
}

// DeleteProject deletes a project
func (h *AdminHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")

	result, err := h.db.Exec("DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(
			"Not found",
			"Project not found",
		))
		return
	}

	c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Project deleted", nil))
}

// GetPortal returns portal configuration
func (h *AdminHandler) GetPortal(c *gin.Context) {
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

// UpdatePortal updates portal configuration
func (h *AdminHandler) UpdatePortal(c *gin.Context) {
	var req domain.UpdatePortalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Invalid request",
			err.Error(),
		))
		return
	}

	// Check if portal config exists
	var id int
	err := h.db.QueryRow("SELECT id FROM portal_config LIMIT 1").Scan(&id)

	if err == sql.ErrNoRows {
		// Create new portal config
		_, err := h.db.Exec(
			"INSERT INTO portal_config (name, description, logo_url, website, email, phone, address) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			req.Name, req.Description, req.LogoURL, req.Website, req.Email, req.Phone, req.Address,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
				"Database error",
				err.Error(),
			))
			return
		}
	} else {
		// Update existing portal config
		_, err := h.db.Exec(
			"UPDATE portal_config SET name = COALESCE(NULLIF($1, ''), name), description = COALESCE(NULLIF($2, ''), description), logo_url = COALESCE(NULLIF($3, ''), logo_url), website = COALESCE(NULLIF($4, ''), website), email = COALESCE(NULLIF($5, ''), email), phone = COALESCE(NULLIF($6, ''), phone), address = COALESCE(NULLIF($7, ''), address), updated_at = CURRENT_TIMESTAMP WHERE id = $8",
			req.Name, req.Description, req.LogoURL, req.Website, req.Email, req.Phone, req.Address, id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
				"Database error",
				err.Error(),
			))
			return
		}
	}

	c.JSON(http.StatusOK, domain.NewAPIResponse(true, "Portal updated", nil))
}
