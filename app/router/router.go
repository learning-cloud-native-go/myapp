package router

import (
	"github.com/go-chi/chi"

	"myapp/app/app"
	"myapp/app/handler"
	"myapp/app/router/middleware"
)

func New(a *app.App) *chi.Mux {
	l := a.Logger()

	r := chi.NewRouter()

	// Routes for healthz
	r.Get("/healthz/liveness", app.HandleLive)
	r.Method("GET", "/healthz/readiness", handler.NewHandler(a.HandleReady, l))

	// Routes for APIs
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		// Routes for books
		r.Method("GET", "/books", handler.NewHandler(a.HandleListBooks, l))
		r.Method("POST", "/books", handler.NewHandler(a.HandleCreateBook, l))
		r.Method("GET", "/books/{id}", handler.NewHandler(a.HandleReadBook, l))
		r.Method("PUT", "/books/{id}", handler.NewHandler(a.HandleUpdateBook, l))
		r.Method("DELETE", "/books/{id}", handler.NewHandler(a.HandleDeleteBook, l))
	})

	r.Method("GET", "/", handler.NewHandler(a.HandleIndex, l))

	return r
}
