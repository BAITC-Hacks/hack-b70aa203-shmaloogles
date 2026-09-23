package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflight(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/tasks", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()

	CORS(http.NotFoundHandler(), "http://localhost:3000").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("unexpected allowed origin: %q", got)
	}
}

func TestCORSRejectsUnknownPreflightOrigin(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/tasks", nil)
	request.Header.Set("Origin", "https://example.test")
	recorder := httptest.NewRecorder()

	CORS(http.NotFoundHandler(), "http://localhost:3000").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}
