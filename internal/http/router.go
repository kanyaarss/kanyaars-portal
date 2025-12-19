package http

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/kanyaarss/kanyaars-portal/internal/config"
	"github.com/kanyaarss/kanyaars-portal/internal/http/handlers"
	"github.com/kanyaarss/kanyaars-portal/internal/http/middleware"
)

// NewRouter creates and configures the Gin router
func NewRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	// Set Gin mode
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Load HTML Templates
	router.LoadHTMLGlob("web/templates/*.html")
	
	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.CORS))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg.JWT.Secret)
	publicHandler := handlers.NewPublicHandler(db)
	apiHandler := handlers.NewAPIHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	// Public routes
	router.GET("/", publicHandler.Home)
	router.GET("/projects", publicHandler.Projects)
	router.GET("/projects/:slug", publicHandler.ProjectDetail)

	// API routes (public)
	api := router.Group("/api/v1")
	{
		api.GET("/health", apiHandler.HealthCheck)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/portal", apiHandler.GetPortal)
		api.GET("/projects", apiHandler.GetProjects)
		api.GET("/projects/:id", apiHandler.GetProject)
	}

	// Admin routes (protected)
	admin := router.Group("/admin")
	admin.Use(middleware.Auth(cfg.JWT.Secret))
	{
		admin.GET("/", adminHandler.Dashboard)
		admin.GET("/projects", adminHandler.ListProjects)
		admin.POST("/projects", adminHandler.CreateProject)
		admin.GET("/projects/:id", adminHandler.GetProject)
		admin.PUT("/projects/:id", adminHandler.UpdateProject)
		admin.DELETE("/projects/:id", adminHandler.DeleteProject)
		admin.GET("/portal", adminHandler.GetPortal)
		admin.PUT("/portal", adminHandler.UpdatePortal)
	}

	// Static files
	router.Static("/static", "./web/static")

	return router
}
