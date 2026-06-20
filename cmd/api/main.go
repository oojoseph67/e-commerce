package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/config"
	"github.com/oojoseph67/ecommerce/internal/database"
	"github.com/oojoseph67/ecommerce/internal/logger"
)

func main() {
	log := logger.New()
	configuration, err := config.Load()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(&configuration.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}

	defer func() {
		if err := mainDB.Close(); err != nil {
			log.Fatal().Err(err).Msg("failed to close database connection")
		}
	}()

	gin.SetMode(configuration.Server.GinMode)

	log.Info().Msg("starting server")

}
