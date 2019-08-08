package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"myapp/repository"
)

func (app *App) HandleListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repository.ListBooks(app.db)
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	if books == nil {
		fmt.Fprint(w, "[]")
		return
	}

	dtos := books.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
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
