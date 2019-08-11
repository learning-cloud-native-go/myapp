package app

import (
	"github.com/jinzhu/gorm"

	"myapp/util/logger"
)

const (
	appErrDataCreationFailure = "data creation failure"
	appErrDataAccessFailure   = "data access failure"
	appErrDataUpdateFailure   = "data update failure"
	appErrJsonCreationFailure = "json creation failure"
	appErrFormDecodingFailure = "form decoding failure"
)

type App struct {
	logger *logger.Logger
	db     *gorm.DB
}

func New(
	logger *logger.Logger,
	db *gorm.DB,
) *App {
	return &App{
		logger: logger,
		db:     db,
	}
}

func (app *App) Logger() *logger.Logger {
	return app.logger
}
