package middleware

import (
	"encoding/json/v2"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/hlog"

	"myapp/pkg/ctxutil"
	e "myapp/pkg/errors"
	vl "myapp/pkg/validator"
)

const maxBytes = 65536 // 64 * 1024 = 64KB

func Validate[T any](v *validator.Validate) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var form T

			if err := json.UnmarshalRead(http.MaxBytesReader(w, r.Body, maxBytes), &form); err != nil {
				hlog.FromRequest(r).Error().Err(err).Msg("")
				e.BadRequest(w, e.RespJSONDecodeFailure)
				return
			}

			if err := v.Struct(form); err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				if err := json.MarshalWrite(w, vl.ToErrResponse(err)); err != nil {
					hlog.FromRequest(r).Error().Err(err).Msg("")
					e.ServerError(w, e.RespJSONEncodeFailure)
					return
				}
				return
			}

			next.ServeHTTP(w, r.WithContext(ctxutil.SetValidatedForm(r.Context(), form)))
		})
	}
}
