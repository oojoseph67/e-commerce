package server

import (
	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	logger zerolog.Logger
}

func New(config *config.Config, db *gorm.DB, logger zerolog.Logger) *Server {
	return &Server{
		config: config,
		db:     db,
		logger: logger,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {

	router := gin.New()

	// middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// routes
	router.GET("/health", s.healthCheck)

	return router
}
