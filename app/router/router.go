package router

import (
	"github.com/go-chi/chi"

	"myapp/app/app"
	"myapp/app/handler"
)

func New(a *app.App) *chi.Mux {
	l := a.Logger()

	r := chi.NewRouter()
	r.Method("GET", "/", handler.NewHandler(a.HandleIndex, l))

	return r
}
