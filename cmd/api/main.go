package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/config"
	"github.com/oojoseph67/ecommerce/internal/database"
	"github.com/oojoseph67/ecommerce/internal/logger"
	"github.com/oojoseph67/ecommerce/internal/server"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

func main() {
	log := logger.New()
	configuration, err := config.Load()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if err := validators.RegisterValidators(); err != nil {
		log.Fatal().Err(err).Msg("failed to register validators")
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

	srv := server.New(configuration, db, log)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", configuration.Server.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	go func() {
		log.Info().Str("port", configuration.Server.Port).Msg("starting http server")
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Panic().Err(err).Msg("failed to shutdown http server")
		return
	}

	log.Info().Msg("shutting down database")

}
