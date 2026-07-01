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
	adminMiddleware := middleware.Admin()

	// services constructor
	service := services.NewService(s.db, s.config, s.logger)

	// handlers constructor
	handler := handlers.NewHandler(
		service.AuthService,
		service.UserService,
		service.ProductService,
		service.CategoryService,
		service.UploadService,
		service.CartService,
		service.OrderService,
	)

	// ROUTES

	// upload
	router.Static("/uploads", "./uploads")

	// API v1
	v1 := router.Group("/api/v1")
	{
		// health
		v1.GET("/health", handlers.HealthCheck)

		// auth
		auth := v1.Group("/auth")
		{
			auth.POST("/admin/signup", handler.AdminSignup)
			auth.POST("/admin/login", handler.AdminLogin)
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

		// categories
		categories := v1.Group("/categories")
		{
			categories.GET("/", authMiddleware, handler.GetCategories)
			categories.POST("/", authMiddleware, adminMiddleware, handler.CreateCategory)
			categories.PUT("/:id", authMiddleware, adminMiddleware, handler.UpdateCategory)
			categories.DELETE("/:id", authMiddleware, adminMiddleware, handler.DeleteCategory)
		}

		// product
		products := v1.Group("/products")
		{
			products.GET("/", authMiddleware, handler.GetProducts)
			products.GET("/:id", authMiddleware, handler.GetProduct)
			products.POST("/", authMiddleware, adminMiddleware, handler.CreateProduct)
			products.PUT("/:id", authMiddleware, adminMiddleware, handler.UpdateProduct)
			products.DELETE("/:id", authMiddleware, adminMiddleware, handler.DeleteProduct)
			products.POST("/:id/upload", authMiddleware, adminMiddleware, handler.UploadProductImage)
			products.DELETE("/:id/upload", authMiddleware, adminMiddleware, handler.DeleteProductImage)
		}

		// cart
		cart := v1.Group("/cart")
		{
			cart.GET("/", authMiddleware, handler.GetCart)
			cart.POST("/items", authMiddleware, handler.AddItemToCart)
			cart.PATCH("/items/:id", authMiddleware, handler.UpdateCartItem)
			cart.DELETE("/items/:id", authMiddleware, handler.RemoveCartItem)
		}

		// order
		order := v1.Group("/order")
		{
			order.GET("/", authMiddleware, handler.GetOrders)
			order.GET("/:id", authMiddleware, handler.GetOrder)
			order.POST("/", authMiddleware, handler.CreateOrder)
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
