//go:generate sh -c "go tool swag init -d ../../ -g cmd/app/main.go -o ../../ -ot yaml --parseDependency && mv ../../swagger.yaml ../../openapi.yaml"

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"myapp/app/router"
	"myapp/config"
	"myapp/pkg/validator"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

//	@title			MYAPP API
//	@version		1.0
//	@description	This is a sample RESTful API with a CRUD

//	@contact.name	Dumindu Madunuwan
//	@contact.url	https://github.com/dumindu

//	@license.name	MIT License
//	@license.url	https://github.com/learning-cloud-native-go/myapp/blob/master/LICENSE

// @servers.url	localhost:8080/v1
func main() {
	c := config.New()

	loc, _ := time.LoadLocation(c.TZ)
	time.Local = loc

	logLevel := zerolog.InfoLevel
	if c.Server.Debug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()

	v := validator.New()

	logLevelDB := gormlogger.Error
	if c.DB.Debug {
		logLevelDB = gormlogger.Info
	}

	dbString := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.Username, c.DB.Password, c.DB.DBName, c.DB.Port)
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevelDB)})
	if err != nil {
		l.Fatal().Err(err).Msg("DB connection start failure")
		return
	}

	mux := router.New(&c.CORS, &l, v, db)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      mux,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		l.Info().Msgf("Shutting down server %v", s.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), c.Server.TimeoutIdle)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			l.Error().Err(err).Msg("Server shutdown failure")
		}

		sqlDB, err := db.DB()
		if err == nil {
			if err = sqlDB.Close(); err != nil {
				l.Error().Err(err).Msg("DB connection closing failure")
			}
		}

		close(closed)
	}()

	l.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	l.Info().Msgf("Server shutdown successfully")
}
