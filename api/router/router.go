package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"myapp/api/requestlog"
	"myapp/api/resource/book"
	"myapp/api/resource/health"
	"myapp/api/router/middleware"
	"myapp/util/logger"
)

func New(l *logger.Logger, v *validator.Validate, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		bookAPI := book.New(l, v, db)
		r.Method("GET", "/books", requestlog.NewHandler(bookAPI.List, l))
		r.Method("POST", "/books", requestlog.NewHandler(bookAPI.Create, l))
		r.Method("GET", "/books/{id}", requestlog.NewHandler(bookAPI.Read, l))
		r.Method("PUT", "/books/{id}", requestlog.NewHandler(bookAPI.Update, l))
		r.Method("DELETE", "/books/{id}", requestlog.NewHandler(bookAPI.Delete, l))
	})

	return r
}
