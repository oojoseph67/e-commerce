package services

import (
	"time"

	"github.com/oojoseph67/ecommerce/internal/config"
	"github.com/oojoseph67/ecommerce/internal/events"
	"github.com/oojoseph67/ecommerce/internal/providers"
	"github.com/oojoseph67/ecommerce/internal/utils/interfaces"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const (
	defaultDateFormat = time.RFC1123
)

// STRUCT COMPOSITION: creating a parent struct and splitting the methods into sub-service struct
// embedding

type BaseService struct {
	db             *gorm.DB
	config         *config.Config
	logger         zerolog.Logger
	eventPublisher events.Publisher
}

type AuthService struct {
	BaseService
}

type UserService struct {
	BaseService
}

type ProductService struct {
	BaseService
	categoryService CategoryService
}

type CategoryService struct {
	BaseService
}

type UploadService struct {
	BaseService
	provider interfaces.UploadProvider
}

type CartService struct {
	BaseService
}

type OrderService struct {
	BaseService
}

type Services struct {
	AuthService     *AuthService
	UserService     *UserService
	ProductService  *ProductService
	CategoryService *CategoryService
	UploadService   *UploadService
	CartService     *CartService
	OrderService    *OrderService
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

func NewService(db *gorm.DB, cfg *config.Config, logger zerolog.Logger, eventPublisher events.Publisher) *Services {

	base := BaseService{
		db:             db,
		config:         cfg,
		logger:         logger,
		eventPublisher: eventPublisher,
	}

	categoryService := &CategoryService{BaseService: base}

	var uploadProvider interfaces.UploadProvider

	if cfg.AWS.UploadStorage == "aws" {
		uploadProvider = providers.NewS3Provider(&cfg.AWS, logger)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path, *base.internalLogger("upload"))
	}

	return &Services{
		AuthService: &AuthService{BaseService: base},
		UserService: &UserService{BaseService: base},
		ProductService: &ProductService{
			BaseService:     base,
			categoryService: *categoryService,
		},
		CategoryService: categoryService,
		UploadService:   NewUploadService(uploadProvider, base),
		CartService:     &CartService{BaseService: base},
		OrderService:    &OrderService{BaseService: base},
	}
}

func NewUploadService(provider interfaces.UploadProvider, base BaseService) *UploadService {
	return &UploadService{
		provider:    provider,
		BaseService: base,
	}
}

// func coalesce(newVal, oldVal string) string {
// 	if newVal == "" {
// 		return oldVal
// 	}
// 	return newVal
// }
