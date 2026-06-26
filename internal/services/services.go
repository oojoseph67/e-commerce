package services

import (
	"github.com/oojoseph67/ecommerce/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// STRUCT COMPOSITION: creating a parent struct and splitting the methods into sub-service struct
// embedding

type BaseService struct {
	db     *gorm.DB
	config *config.Config
	logger zerolog.Logger
}

type AuthService struct {
	BaseService
}

type UserService struct {
	BaseService
}

type ProductService struct {
	BaseService
}

type Services struct {
	AuthService    *AuthService
	UserService    *UserService
	ProductService *ProductService
}

func (s *BaseService) internalLogger(service string) *zerolog.Logger {
	l := s.logger.With().Str("service", service).Logger()
	return &l
}

func coalesce[T comparable](newVal, oldVal T) T {
	var zero T
	if newVal == zero {
		return oldVal
	}
	return newVal
}

// use when checking for bool
func coalescePtr[T any](newVal *T, oldVal T) T {
	if newVal == nil {
		return oldVal
	}
	return *newVal
}

func NewService(db *gorm.DB, cfg *config.Config, logger zerolog.Logger) *Services {

	base := BaseService{
		db:     db,
		config: cfg,
		logger: logger,
	}

	return &Services{
		AuthService:    &AuthService{BaseService: base},
		UserService:    &UserService{BaseService: base},
		ProductService: &ProductService{BaseService: base},
	}
}

// func coalesce(newVal, oldVal string) string {
// 	if newVal == "" {
// 		return oldVal
// 	}
// 	return newVal
// }
