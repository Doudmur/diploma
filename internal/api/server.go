package api

import (
	"diploma/internal/api/routes"
	"diploma/internal/config"
	"github.com/gin-gonic/gin"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	config *config.Config
}

// NewServer creates a new Server instance
func NewServer(cfg *config.Config) *Server {
	router := gin.Default()
	return &Server{
		router: router,
		config: cfg,
	}
}

// Start runs the API server
func (s *Server) Start() error {
	// Set up routes
	routes.SetupRoutes(s.router)

	// Start the server
	return s.router.Run(":" + s.config.Port)
}
