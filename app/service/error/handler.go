package error

import (
	"fmt"
	"net/http"
)

const (
	ErrDataCreationFailure = "data creation failure"
	ErrDataAccessFailure   = "data access failure"
	ErrDataUpdateFailure   = "data update failure"
	ErrDataDeletionFailure = "data deletion failure"

	ErrJsonCreationFailure    = "json creation failure"
	ErrFormDecodingFailure    = "form decoding failure"
	ErrFormErrResponseFailure = "form error response failure"

	ErrInvalidIdInUrlParam = "invalid id in url param"
)

func AppError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error": "%v"}`, msg)
}

func ValError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	fmt.Fprintf(w, `{"error": "%v"}`, msg)
}

func FormValError(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(reps)
}
