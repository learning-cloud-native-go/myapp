package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	dbConn "myapp/adapter/gorm"
	"myapp/app/router"
	"myapp/config"
	lr "myapp/util/logger"
	vr "myapp/util/validator"
)

func main() {
	appConf := config.AppConfig()

	logger := lr.New(appConf.Debug)

	validator := vr.New()

	db, err := dbConn.New(&appConf.Db)
	if err != nil {
		logger.Fatal().Err(err).Msg("Db connection start failure")
		return
	}

	appRouter := router.New(logger, validator, db)
	address := fmt.Sprintf(":%d", appConf.Server.Port)
	srv := &http.Server{
		Addr:         address,
		Handler:      appRouter,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		logger.Info().Msgf("Shutting down server %v", address)

		ctx, cancel := context.WithTimeout(context.Background(), appConf.Server.TimeoutIdle)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("Server shutdown failure")
		}

		sqlDB, err := db.DB()
		if err == nil {
			if err = sqlDB.Close(); err != nil {
				logger.Error().Err(err).Msg("Db connection closing failure")
			}
		}

		close(closed)
	}()

	logger.Info().Msgf("Starting server %v", address)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
}
