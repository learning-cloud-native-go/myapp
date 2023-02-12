package err

import (
	"fmt"
	"net/http"
)

const (
	DataCreationFailure = "data creation failure"
	DataAccessFailure   = "data access failure"
	DataUpdateFailure   = "data update failure"
	DataDeletionFailure = "data deletion failure"

	JsonEncodingFailure = "json encoding failure"
	JsonDecodingFailure = "json decoding failure"

	FormErrResponseFailure = "form error response failure"

	InvalidIdInUrlParam = "invalid id in url param"
)

type Error struct {
	Error string `json:"error"`
}

type Errors struct {
	Errors []string `json:"errors"`
}

func ServerError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error": "%v"}`, msg)
}

func BadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `{"error": "%v"}`, msg)
}

func ValidationErrors(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(reps)
}
