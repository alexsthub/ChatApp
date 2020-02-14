package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var expectedHeaders = map[string]string{
	"Access-Control-Allow-Origin":   "*",
	"Access-Control-Allow-Methods":  "GET, PUT, POST, PATCH, DELETE",
	"Access-Control-Allow-Headers":  "Content-Type, Authorization",
	"Access-Control-Expose-Headers": "Authorization",
	"Access-Control-Max-Age":        "600",
}

func TestCors(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	cors := CorsMW{h}

	req := httptest.NewRequest("GET", "https://api.alexst.me/v1/sessions", nil)
	w := httptest.NewRecorder()
	cors.ServeHTTP(w, req)
	resp := w.Result()
	for key, value := range expectedHeaders {
		seenHeader := resp.Header.Get(key)
		if seenHeader != value {
			t.Errorf("Header: %s value is %s but should be %s", key, seenHeader, value)
		}
	}
}
