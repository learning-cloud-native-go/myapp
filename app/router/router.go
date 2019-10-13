package router

import (
	"github.com/go-chi/chi"

	"myapp/app/app"
	"myapp/app/requestlog"
	"myapp/app/router/middleware"
)

func New(a *app.App) *chi.Mux {
	l := a.Logger()

	r := chi.NewRouter()

	// Routes for healthz
	r.Get("/healthz/liveness", app.HandleLive)
	r.Method("GET", "/healthz/readiness", requestlog.NewHandler(a.HandleReady, l))

	// Routes for APIs
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		// Routes for books
		r.Method("GET", "/books", requestlog.NewHandler(a.HandleListBooks, l))
		r.Method("POST", "/books", requestlog.NewHandler(a.HandleCreateBook, l))
		r.Method("GET", "/books/{id}", requestlog.NewHandler(a.HandleReadBook, l))
		r.Method("PUT", "/books/{id}", requestlog.NewHandler(a.HandleUpdateBook, l))
		r.Method("DELETE", "/books/{id}", requestlog.NewHandler(a.HandleDeleteBook, l))
	})

	r.Method("GET", "/", requestlog.NewHandler(a.HandleIndex, l))

	return r
}
