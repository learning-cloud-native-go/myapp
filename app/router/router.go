package router

import (
	"github.com/go-chi/chi"

	"myapp/app/app"
	"myapp/app/handler"
)

func New(a *app.App) *chi.Mux {
	l := a.Logger()

	r := chi.NewRouter()

	// Routes for healthz
	r.Get("/healthz/liveness", app.HandleLive)
	r.Method("GET", "/healthz/readiness", handler.NewHandler(a.HandleReady, l))

	r.Method("GET", "/", handler.NewHandler(a.HandleIndex, l))

	return r
}
