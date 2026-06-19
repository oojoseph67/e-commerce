package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/config"
	"github.com/oojoseph67/ecommerce/internal/database"
	"github.com/oojoseph67/ecommerce/internal/logger"
)

func main() {
	log := logger.New()
	config, err := config.Load()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(config.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}

	defer mainDB.Close()
	gin.SetMode(config.Server.GinMode)

	log.Info().Msg("starting server")

}
