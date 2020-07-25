package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

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

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info().Msgf("main :  API listening on %s.", s.Addr)
		serverErrors <- s.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, os.Kill)

	select {
	case err := <-serverErrors:
		logger.Info().Msgf("error : listening and serving %s", err)

	case <-shutdown:
		logger.Info().Msg("main : Start shutdown")

		timeout := appConf.Server.TimeoutIdle
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			logger.Info().Msgf("main : Graceful shutdown did not complete in %v : %v", timeout, err)

			if err = s.Close(); err != nil {
				logger.Fatal().Err(err).Msgf("main : Could not stop server gracefully %v", err)
			}
		}

	}

}
