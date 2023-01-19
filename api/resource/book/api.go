package book

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"myapp/util/logger"
)

type API struct {
	logger     *logger.Logger
	validator  *validator.Validate
	repository *Repository
}

func New(logger *logger.Logger, validator *validator.Validate, db *gorm.DB) *API {
	return &API{
		logger:     logger,
		validator:  validator,
		repository: NewRepository(db),
	}
}
