package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"myapp/app/book"
	"myapp/form"
	"myapp/pkg/logger"
	"myapp/pkg/middleware"
	rl "myapp/pkg/middleware/requestlog"
)

func New(l *logger.Logger, v *validator.Validate, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("."))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.ContentTypeJSON)

		bookAPI := book.New(l, v, db)
		r.Get("/books", rl.NewHandler(bookAPI.List, l).ServeHTTP)
		r.With(middleware.Validate[form.BookForm](l, v)).Post("/books", rl.NewHandler(bookAPI.Create, l).ServeHTTP)
		r.Get("/books/{id}", rl.NewHandler(bookAPI.Read, l).ServeHTTP)
		r.With(middleware.Validate[form.BookForm](l, v)).Put("/books/{id}", rl.NewHandler(bookAPI.Update, l).ServeHTTP)
		r.Delete("/books/{id}", rl.NewHandler(bookAPI.Delete, l).ServeHTTP)
	})

	return r
}
