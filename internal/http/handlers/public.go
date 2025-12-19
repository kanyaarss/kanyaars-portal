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
	rows, err := h.db.Query(`
		SELECT id, name, slug, description, url, icon_url, status
		FROM projects
		WHERE status = 'active'
		ORDER BY "order" ASC
	`)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to fetch projects",
		})
		return
	}
	defer rows.Close()

	projects := []map[string]interface{}{}

	for rows.Next() {
		var (
			id          int
			name        string
			slug        string
			description string
			url         string
			iconURL     string
			status      string
		)

		if err := rows.Scan(
			&id,
			&name,
			&slug,
			&description,
			&url,
			&iconURL,
			&status,
		); err != nil {
			continue
		}

		projects = append(projects, gin.H{
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

	var (
		id          int
		name        string
		slugDB      string
		description string
		url         string
		iconURL     string
		status      string
	)

	err := h.db.QueryRow(`
		SELECT id, name, slug, description, url, icon_url, status
		FROM projects
		WHERE slug = $1
	`, slug).Scan(
		&id,
		&name,
		&slugDB,
		&description,
		&url,
		&iconURL,
		&status,
	)

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

	project := gin.H{
		"id":          id,
		"name":        name,
		"slug":        slugDB,
		"description": description,
		"url":         url,
		"icon_url":    iconURL,
		"status":      status,
	}

	c.HTML(http.StatusOK, "project-detail.html", gin.H{
		"title":   "Project Detail",
		"project": project,
	})
}