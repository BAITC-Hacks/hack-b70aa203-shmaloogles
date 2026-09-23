package ai

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// In-memory HTTP transport: no listening sockets, credentials or network required.
func TestOpenAIHTTP(t *testing.T) {
	text, _ := json.Marshal(validQuestions)
	valid := `{"status":"completed","output":[{"type":"reasoning","summary":[]},{"type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":` + string(text) + `}]}]}`
	for _, tc := range []struct {
		name, body string
		status     int
		want       error
	}{
		{"valid", valid, 200, nil},
		{"provider error", "secret", 500, ErrProvider},
		{"rate limit", "secret", 429, ErrProvider},
		{"auth", "secret", 401, ErrProvider},
		{"redirect", "", 302, ErrProvider},
		{"broken envelope", "{", 200, ErrResponse},
		{"empty output", `{"status":"completed","output":[]}`, 200, ErrResponse},
		{"refusal", `{"status":"completed","output":[{"type":"message","status":"completed","role":"assistant","content":[{"type":"refusal","refusal":"Unable to help"}]}]}`, 200, ErrResponse},
		{"truncated", strings.Replace(valid, `"completed"`, `"incomplete"`, 1), 200, ErrResponse},
		{"failed", strings.Replace(valid, `"completed"`, `"failed"`, 1), 200, ErrResponse},
		{"message incomplete", strings.ReplaceAll(valid, `"status":"completed","content"`, `"status":"incomplete","content"`), 200, ErrResponse},
		{"wrong role", strings.Replace(valid, `"assistant"`, `"user"`, 1), 200, ErrResponse},
		{"unexpected tool", strings.Replace(valid, `"message"`, `"function_call"`, 1), 200, ErrResponse},
		{"missing text", `{"status":"completed","output":[{"type":"message","status":"completed","role":"assistant","content":[]}]}`, 200, ErrResponse},
		{"too large", strings.Repeat(" ", (1<<20)+1), 200, ErrResponse},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := &openAI{endpoint: "https://example.invalid/v1/responses", model: "test-model", key: "test-key", client: &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer test-key" || r.Header.Get("Content-Type") != "application/json" || r.URL.Path != "/v1/responses" || r.URL.RawQuery != "" {
					t.Fatalf("bad request: %v", r)
				}
				var body struct {
					Model           string `json:"model"`
					Instructions    string `json:"instructions"`
					Input           string `json:"input"`
					Store           *bool  `json:"store"`
					MaxOutputTokens int    `json:"max_output_tokens"`
					Text            struct {
						Format struct {
							Type   string          `json:"type"`
							Name   string          `json:"name"`
							Strict bool            `json:"strict"`
							Schema json.RawMessage `json:"schema"`
						} `json:"format"`
					} `json:"text"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body.Model != "test-model" || body.Instructions != clarifyPrompt || !json.Valid([]byte(body.Input)) || body.Store == nil || *body.Store || body.MaxOutputTokens != 4096 || body.Text.Format.Type != "json_schema" || body.Text.Format.Name == "" || !body.Text.Format.Strict || !json.Valid(body.Text.Format.Schema) {
					t.Fatalf("bad body: %+v", body)
				}
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
			})}}
			got, err := New(g, time.Second, false).Clarify(context.Background(), "Описание")
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v want %v", err, tc.want)
			}
			if err == nil && got.Mode != "provider" {
				t.Fatal(got)
			}
		})
	}
}

func TestOpenAITransportErrorAndTimeout(t *testing.T) {
	for _, timeout := range []bool{false, true} {
		g := &openAI{endpoint: "https://example.invalid", client: &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
			if timeout {
				<-r.Context().Done()
				return nil, r.Context().Err()
			}
			return nil, errors.New("sensitive error")
		})}}
		_, err := New(g, time.Millisecond, false).Clarify(context.Background(), "x")
		want := ErrProvider
		if timeout {
			want = context.DeadlineExceeded
		}
		if !errors.Is(err, want) || strings.Contains(err.Error(), "sensitive") {
			t.Fatal(err)
		}
	}
}

func TestProviderOwnTimeout(t *testing.T) {
	g := &openAI{endpoint: "https://example.invalid", client: &http.Client{Transport: transportFunc(func(*http.Request) (*http.Response, error) {
		return nil, context.DeadlineExceeded
	})}}
	got, err := New(g, time.Second, true).Clarify(context.Background(), "x")
	if err != nil || got.FallbackReason != "timeout" {
		t.Fatalf("%+v %v", got, err)
	}
}

func TestFromEnv(t *testing.T) {
	for _, tc := range []struct {
		name  string
		env   map[string]string
		valid bool
		mock  bool
	}{
		{"default", nil, true, true},
		{"mock", map[string]string{"AI_MODE": "mock"}, true, true},
		{"unknown", map[string]string{"AI_MODE": "other"}, false, false},
		{"missing credentials", map[string]string{"AI_MODE": "openai"}, false, false},
		{"openai", map[string]string{"AI_MODE": "openai", "OPENAI_API_KEY": "test", "AI_MODEL": "test-model", "AI_FALLBACK": "true", "AI_TIMEOUT": "2s"}, true, false},
		{"default model", map[string]string{"AI_MODE": "openai", "OPENAI_API_KEY": "test", "AI_FALLBACK": "true", "AI_TIMEOUT": "2s"}, true, false},
		{"bad model", map[string]string{"AI_MODE": "openai", "OPENAI_API_KEY": "test", "AI_MODEL": "../oops"}, false, false},
		{"bad timeout", map[string]string{"AI_MODE": "openai", "OPENAI_API_KEY": "test", "AI_MODEL": "test", "AI_TIMEOUT": "-1s"}, false, false},
		{"bad bool", map[string]string{"AI_MODE": "openai", "OPENAI_API_KEY": "test", "AI_MODEL": "test", "AI_FALLBACK": "maybe"}, false, false},
		{"insecure URL", map[string]string{"AI_MODE": "openai", "OPENAI_API_KEY": "test", "AI_MODEL": "test", "AI_BASE_URL": "http://example.invalid"}, false, false},
		{"query in URL", map[string]string{"AI_MODE": "openai", "OPENAI_API_KEY": "test", "AI_MODEL": "test", "AI_BASE_URL": "https://example.invalid?key=secret"}, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, key := range []string{"AI_MODE", "OPENAI_API_KEY", "AI_MODEL", "AI_TIMEOUT", "AI_BASE_URL", "AI_FALLBACK"} {
				t.Setenv(key, "")
			}
			for key, value := range tc.env {
				t.Setenv(key, value)
			}
			s, err := FromEnv()
			if (err == nil) != tc.valid {
				t.Fatalf("%v", err)
			}
			if err == nil && (s.provider == nil) != tc.mock {
				t.Fatalf("wrong mode: %+v", s)
			}
			if err == nil && !tc.mock {
				g := s.provider.(*openAI)
				model := tc.env["AI_MODEL"]
				if model == "" {
					model = "gpt-4.1-mini"
				}
				if g.model != model || g.endpoint != "https://api.openai.com/v1/responses" {
					t.Fatal("wrong model/endpoint")
				}
				if !s.fallback || s.timeout != 2*time.Second || g.client.CheckRedirect(nil, nil) != http.ErrUseLastResponse {
					t.Fatalf("wrong config: %+v", s)
				}
			}
		})
	}
}
