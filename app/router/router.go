package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"gorm.io/gorm"

	"myapp/app/book"
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
		r.Route("/books", book.New(v, db).Register)
	})

	return r
}
