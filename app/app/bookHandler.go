package app

import (
	"net/http"
)

func (app *App) HandleListBooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[]"))
}

func (app *App) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleReadBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (app *App) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func (app *App) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}
