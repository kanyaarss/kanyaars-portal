package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	db *sql.DB
}

// NewPublicHandler creates a new public handler
func NewPublicHandler(db *sql.DB) *PublicHandler {
	return &PublicHandler{db: db}
}

// Home renders the home page
func (h *PublicHandler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Kanyaars Portal",
	})
}

// Projects renders the projects page
func (h *PublicHandler) Projects(c *gin.Context) {
	// Fetch projects from database
	rows, err := h.db.Query(
		"SELECT id, name, slug, description, url, icon_url, status FROM projects WHERE status = 'active' ORDER BY \"order\" ASC",
	)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to fetch projects",
		})
		return
	}
	defer rows.Close()

	var projects []map[string]interface{}
	for rows.Next() {
		var id int
		var name, slug, description, url, iconURL, status string
		if err := rows.Scan(&id, &name, &slug, &description, &url, &iconURL, &status); err != nil {
			continue
		}

		projects = append(projects, map[string]interface{}{
			"id":          id,
			"name":        name,
			"slug":        slug,
			"description": description,
			"url":         url,
			"icon_url":    iconURL,
			"status":      status,
		})
	}

	c.HTML(http.StatusOK, "projects.html", gin.H{
		"title":    "Projects",
		"projects": projects,
	})
}

// ProjectDetail renders the project detail page
func (h *PublicHandler) ProjectDetail(c *gin.Context) {
	slug := c.Param("slug")

	var id int
	var name, slug2, description, url, iconURL, status string
	err := h.db.QueryRow(
		"SELECT id, name, slug, description, url, icon_url, status FROM projects WHERE slug = $1",
		slug,
	).Scan(
		&id,
		&name,
		&slug2,
		&description,
		&url,
		&iconURL,
		&status,
	)

	project := map[string]interface{}{
		"id":          id,
		"name":        name,
		"slug":        slug2,
		"description": description,
		"url":         url,
		"icon_url":    iconURL,
		"status":      status,
	}

	if err == sql.ErrNoRows {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Project not found",
		})
		return
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to fetch project",
		})
		return
	}

	c.HTML(http.StatusOK, "project-detail.html", gin.H{
		"title":   "Project Detail",
		"project": project,
	})
}
