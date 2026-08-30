package router

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"gorm.io/gorm"

	"myapp/app/book"
	"myapp/config"
	m "myapp/pkg/middleware"
	mrl "myapp/pkg/middleware/requestlog"
)

func New(c *config.ConfCORS, l *zerolog.Logger, db *gorm.DB, v *validator.Validate) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("."))
	})

	cors := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(c.AllowedOrigins, ","),
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   strings.Split(c.AllowedHeaders, ","),
		AllowCredentials: true,
		ExposedHeaders:   []string{""},
		MaxAge:           c.MaxAge,
	})

	r.Route("/v1", func(r chi.Router) {
		r.Use(cors.Handler)
		r.Use(m.ContentTypeJSON)
		r.Use(hlog.NewHandler(*l))
		r.Use(hlog.RequestIDHandler("request_id", "X-Request-ID"))
		r.Use(mrl.NewHandler)

		r.Route("/books", book.New(db, v).Register)
	})

	return r
}
