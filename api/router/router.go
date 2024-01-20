package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"myapp/api/resource/book"
	"myapp/api/resource/health"
	"myapp/api/router/middleware"
	"myapp/api/router/middleware/requestlog"
)

func New(l *zerolog.Logger, v *validator.Validate, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.ContentTypeJSON)

		bookAPI := book.New(l, v, db)
		r.Method(http.MethodGet, "/books", requestlog.NewHandler(bookAPI.List, l))
		r.Method(http.MethodPost, "/books", requestlog.NewHandler(bookAPI.Create, l))
		r.Method(http.MethodGet, "/books/{id}", requestlog.NewHandler(bookAPI.Read, l))
		r.Method(http.MethodPut, "/books/{id}", requestlog.NewHandler(bookAPI.Update, l))
		r.Method(http.MethodDelete, "/books/{id}", requestlog.NewHandler(bookAPI.Delete, l))
	})

	return r
}
