package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"myapp/pkg/ctxutil"
	"myapp/pkg/middleware"
)

func TestRequestID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		headerValue string
	}{
		{
			name:        "with header value",
			headerValue: "9m4e2mr0ui3e8a215n4g",
		},
		{
			name:        "without header value",
			headerValue: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			if tt.headerValue != "" {
				r.Header.Set("X-Request-ID", tt.headerValue)
			}

			w := httptest.NewRecorder()
			middleware.RequestID(http.HandlerFunc(testHandlerFuncRequestID())).ServeHTTP(w, r)

			if w.Result().StatusCode != http.StatusOK {
				t.Fatal("context requestID should not be empty")
			}
		})
	}
}

func testHandlerFuncRequestID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := ctxutil.RequestID(r.Context())

		if requestId == "" {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
