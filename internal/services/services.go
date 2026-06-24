package services

import (
	"github.com/oojoseph67/ecommerce/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Services struct {
	db     *gorm.DB
	config *config.Config
	logger zerolog.Logger
}

func NewService(db *gorm.DB, config *config.Config, logger zerolog.Logger) *Services {
	return &Services{
		db:     db,
		config: config,
		logger: logger,
	}
}

func (s *Services) internalLogger(service string) *zerolog.Logger {
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

// func coalesce(newVal, oldVal string) string {
// 	if newVal == "" {
// 		return oldVal
// 	}
// 	return newVal
// }
