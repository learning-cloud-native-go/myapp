package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"myapp/api/router/middleware"
)

const testRespBody = `{"k":"v"}`

func TestContentTypeJSON(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	middleware.ContentTypeJSON(http.HandlerFunc(testHandlerFunc())).ServeHTTP(w, r)
	response := w.Result()

	if respBody := w.Body.String(); respBody != testRespBody {
		t.Errorf("Wrong response body:  got %v want %v ", respBody, testRespBody)
	}

	if status := response.StatusCode; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
	}

	if contentType := response.Header.Get(middleware.HeaderKeyContentType); contentType != middleware.HeaderValueContentTypeJSON {
		t.Errorf("Wrong status code: got %v want %v", contentType, middleware.HeaderValueContentTypeJSON)
	}
}

func testHandlerFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testRespBody)
	}
}
