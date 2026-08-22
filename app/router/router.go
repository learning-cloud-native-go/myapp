package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"gorm.io/gorm"

	"myapp/app/book"
	"myapp/form"
	m "myapp/pkg/middleware"
	mrl "myapp/pkg/middleware/requestlog"
)

func New(l *zerolog.Logger, v *validator.Validate, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("."))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Use(hlog.NewHandler(*l))
		r.Use(hlog.RequestIDHandler("request_id", "X-Request-ID"))
		r.Use(mrl.NewHandler)
		r.Use(m.ContentTypeJSON)

		bookAPI := book.New(v, db)
		r.Get("/books", bookAPI.List)
		r.With(m.Validate[form.BookForm](v)).Post("/books", bookAPI.Create)
		r.Get("/books/{id}", bookAPI.Read)
		r.With(m.Validate[form.BookForm](v)).Put("/books/{id}", bookAPI.Update)
		r.Delete("/books/{id}", bookAPI.Delete)
	})

	return r
}
