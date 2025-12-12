package domain

import "time"

// Project represents a project in the portal
type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	IconURL     string    `json:"icon_url"`
	Status      string    `json:"status"` // active, inactive, maintenance
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProjectRequest represents create project request
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Slug        string `json:"slug" binding:"required,min=3"`
	Description string `json:"description" binding:"required"`
	URL         string `json:"url" binding:"required,url"`
	IconURL     string `json:"icon_url" binding:"url"`
	Status      string `json:"status" binding:"required,oneof=active inactive maintenance"`
}

// UpdateProjectRequest represents update project request
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"min=3"`
	Slug        string `json:"slug" binding:"min=3"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"url"`
	IconURL     string `json:"icon_url" binding:"url"`
	Status      string `json:"status" binding:"oneof=active inactive maintenance"`
}
