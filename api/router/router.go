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

	// Routes for healthz
	r.Get("/healthz", health.Read)

	// Routes for APIs
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		// Routes for books
		srvBook := book.NewApp(l, v, db)
		r.Method("GET", "/books", requestlog.NewHandler(srvBook.List, l))
		r.Method("POST", "/books", requestlog.NewHandler(srvBook.Create, l))
		r.Method("GET", "/books/{id}", requestlog.NewHandler(srvBook.Read, l))
		r.Method("PUT", "/books/{id}", requestlog.NewHandler(srvBook.Update, l))
		r.Method("DELETE", "/books/{id}", requestlog.NewHandler(srvBook.Delete, l))
	})

	return r
}
