package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	dbConn "myapp/adapter/gorm"
	"myapp/app/app"
	"myapp/app/router"
	"myapp/config"
	lr "myapp/util/logger"
	vr "myapp/util/validator"
)

func main() {
	appConf := config.AppConfig()

	logger := lr.New(appConf.Debug)

	db, err := dbConn.New(appConf)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
		return
	}
	if appConf.Debug {
		db.LogMode(true)
	}

	validator := vr.New()

	application := app.New(logger, db, validator)

	appRouter := router.New(application)

	address := fmt.Sprintf(":%d", appConf.Server.Port)

	logger.Info().Msgf("Starting server %v", address)

	s := &http.Server{
		Addr:         address,
		Handler:      appRouter,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	// Make a channel to listen for errors coming from the listener. Use a buffered
	// channel so the goroutine exits if we dont get this error
	serverErrors := make(chan error, 1)

	// Start server in a goroutine so that its does not block
	go func() {
		logger.Info().Msgf("main :  API listening on %s.", s.Addr)
		serverErrors <- s.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, os.Kill)

	// Block main until signal is received
	select {
	case err := <-serverErrors:
		logger.Info().Msgf("error : listening and serving %s", err)

	case <-shutdown:
		logger.Info().Msg("main : Start shutdown")

		// Give outstanding requests a deadline for completion
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking for listener to shutdown and load shed
		if err := s.Shutdown(ctx); err != nil {
			logger.Info().Msgf("main : Graceful shutdown did not complete in %v : %v", timeout, err)

			// Immediately closes all active net.Listeners
			if err = s.Close(); err != nil {
				logger.Fatal().Err(err).Msgf("main : Could not stop server gracefully %v", err)
			}
		}

	}

}
