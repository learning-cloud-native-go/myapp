package book

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"myapp/util/logger"
)

type App struct {
	logger     *logger.Logger
	validator  *validator.Validate
	repository *Repository
}

func NewApp(logger *logger.Logger, validator *validator.Validate, db *gorm.DB) *App {
	return &App{
		logger:     logger,
		validator:  validator,
		repository: NewRepository(db),
	}
}
