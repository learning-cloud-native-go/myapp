package app

import (
	"net/http"
)

// HandleLive is an http.HandlerFunc that handles liveness checks by
// immediately responding with an HTTP 200 status.
func HandleLive(w http.ResponseWriter, _ *http.Request) {
	writeHealthy(w)
}

// HandleReady is an http.HandlerFunc that handles readiness checks by
// responding with an HTTP 200 status if it is healthy, 500 otherwise.
func (app *App) HandleReady(w http.ResponseWriter, r *http.Request) {
	if err := app.db.DB().Ping(); err != nil {
		app.Logger().Fatal().Err(err).Msg("")
		writeUnhealthy(w)
		return
	}

	writeHealthy(w)
}

func writeHealthy(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("."))
}

func writeUnhealthy(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("."))
}
