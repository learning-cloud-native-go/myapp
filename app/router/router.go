package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"myapp/app/requestlog"
	"myapp/app/router/middleware"
	"myapp/app/service/book"
	"myapp/app/service/health"
	"myapp/util/logger"
)

func New(l *logger.Logger, v *validator.Validate, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	// Routes for healthz
	r.Get("/healthz", health.HandleHealth)

	// Routes for APIs
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		// Routes for books
		srvBook := book.NewApp(l, v, db)
		r.Method("GET", "/books", requestlog.NewHandler(srvBook.HandleListBooks, l))
		r.Method("POST", "/books", requestlog.NewHandler(srvBook.HandleCreateBook, l))
		r.Method("GET", "/books/{id}", requestlog.NewHandler(srvBook.HandleReadBook, l))
		r.Method("PUT", "/books/{id}", requestlog.NewHandler(srvBook.HandleUpdateBook, l))
		r.Method("DELETE", "/books/{id}", requestlog.NewHandler(srvBook.HandleDeleteBook, l))
	})

	return r
}
