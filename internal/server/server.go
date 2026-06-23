package server

import (
	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/config"
	"github.com/oojoseph67/ecommerce/internal/handlers"
	"github.com/oojoseph67/ecommerce/internal/middleware"
	"github.com/oojoseph67/ecommerce/internal/services"
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
	router.Use(middleware.CORS())

	authMiddleware := middleware.Auth(s.config.JWT.Secret)

	// services constructor
	service := services.NewService(s.db, s.config, s.logger)

	// handlers constructor
	handler := handlers.NewHandler(service)

	// routes
	// API v1
	v1 := router.Group("/api/v1")
	{
		// health
		v1.GET("/health", handlers.HealthCheck)

		// auth
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", handler.Signup)
			auth.POST("/login", handler.Login)
			auth.POST("/refresh", handler.RefreshToken)
			auth.POST("/logout", handler.Logout)
		}

		// user
		user := v1.Group("/user")
		{
			user.GET("/me", authMiddleware, handler.Me)
			user.PATCH("/update", authMiddleware, handler.UpdateProfile)
		}

		// Example protected routes
		// protected := v1.Group("/user")
		// protected.Use(middleware.Auth(s.config.JWT.Secret))
		// {
		//     protected.GET("/me", profileHandler.GetMe)
		// }

		// Example admin routes
		// admin := v1.Group("/admin")
		// admin.Use(middleware.Auth(s.config.JWT.Secret), middleware.Admin())
		// {
		//     admin.GET("/users", adminHandler.ListUsers)
		// }
	}

	return router
}
