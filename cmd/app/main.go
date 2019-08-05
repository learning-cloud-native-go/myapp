package main

import (
	"fmt"
	"net/http"

	"myapp/app/app"
	"myapp/app/router"
	"myapp/config"
	lr "myapp/util/logger"
)

func main() {
	appConf := config.AppConfig()

	logger := lr.New(appConf.Debug)

	application := app.New(logger)

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

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server startup failed")
	}
}
