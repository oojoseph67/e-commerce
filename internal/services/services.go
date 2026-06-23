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
