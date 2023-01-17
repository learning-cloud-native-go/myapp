package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	gosql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"myapp/api/router"
	"myapp/config"
	"myapp/util/logger"
	"myapp/util/validator"
)

func main() {
	c := config.New()
	l := logger.New(c.Server.Debug)
	v := validator.New()

	cfg := &gosql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", c.DB.Host, c.DB.Port),
		DBName:               c.DB.DBName,
		User:                 c.DB.Username,
		Passwd:               c.DB.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	var logLevel gormlogger.LogLevel
	if c.DB.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		l.Fatal().Err(err).Msg("DB connection start failure")
		return
	}

	r := router.New(l, v, db)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
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
}
