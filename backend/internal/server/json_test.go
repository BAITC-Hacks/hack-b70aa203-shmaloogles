package server

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONRequiresBoundedObject(t *testing.T) {
	for _, body := range []string{"null", "[]", "42", "", `{} {}`, `{"title":"` + strings.Repeat("x", 1<<20) + `"}`} {
		r := httptest.NewRequest("PUT", "/", strings.NewReader(body))
		var input struct {
			Title *string `json:"title"`
		}
		if err := decodeJSON(r, &input); err == nil {
			t.Fatalf("accepted invalid body of length %d", len(body))
		}
	}
	r := httptest.NewRequest("PUT", "/", strings.NewReader(`{"title":null}`))
	var input struct {
		Title *string `json:"title"`
	}
	if err := decodeJSON(r, &input); err != nil {
		t.Fatalf("nullable field must remain supported: %v", err)
	}
}
