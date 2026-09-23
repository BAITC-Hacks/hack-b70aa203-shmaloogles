package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type gemini struct {
	endpoint string
	key      string
	client   *http.Client
}

// FromEnv defaults to labelled mock mode. Misconfigured real mode fails at startup.
func FromEnv() (*Service, error) {
	mode := strings.TrimSpace(os.Getenv("AI_MODE"))
	if mode == "" || mode == "mock" {
		return New(nil, 0, false), nil
	}
	if mode != "gemini" {
		return nil, fmt.Errorf("AI_MODE must be mock or gemini")
	}
	key, model := strings.TrimSpace(os.Getenv("AI_API_KEY")), strings.TrimSpace(os.Getenv("AI_MODEL"))
	if key == "" || !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(model) {
		return nil, fmt.Errorf("gemini requires AI_API_KEY and a valid AI_MODEL")
	}
	base := strings.TrimRight(os.Getenv("AI_BASE_URL"), "/")
	if base == "" {
		base = "https://generativelanguage.googleapis.com/v1beta"
	}
	u, err := url.Parse(base)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("AI_BASE_URL must be an HTTPS URL without credentials, query or fragment")
	}
	timeout := 20 * time.Second
	if value := os.Getenv("AI_TIMEOUT"); value != "" {
		timeout, err = time.ParseDuration(value)
		if err != nil || timeout <= 0 {
			return nil, fmt.Errorf("AI_TIMEOUT must be a positive duration")
		}
	}
	fallback := false
	if value := os.Getenv("AI_FALLBACK"); value != "" {
		fallback, err = strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("AI_FALLBACK must be a boolean")
		}
	}
	client := &http.Client{Timeout: timeout, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	return New(&gemini{endpoint: base + "/models/" + model + ":generateContent", key: key, client: client}, timeout, fallback), nil
}

func (g *gemini) Complete(ctx context.Context, prompt string, input, schema json.RawMessage) ([]byte, error) {
	body, err := json.Marshal(map[string]any{
		"systemInstruction": map[string]any{"parts": []any{map[string]any{"text": prompt}}},
		"contents":          []any{map[string]any{"role": "user", "parts": []any{map[string]any{"text": string(input)}}}},
		"generationConfig":  map[string]any{"responseMimeType": "application/json", "responseJsonSchema": schema},
	})
	if err != nil {
		return nil, ErrProvider
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, ErrProvider
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.key)
	resp, err := g.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, context.DeadlineExceeded
		}
		return nil, ErrProvider
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrProvider
	}
	const maxBody = 1 << 20
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return nil, ErrProvider
	}
	if len(raw) > maxBody {
		return nil, ErrResponse
	}
	var envelope struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text    string `json:"text"`
					Thought bool   `json:"thought"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
	}
	if json.Unmarshal(raw, &envelope) != nil || len(envelope.Candidates) != 1 {
		return nil, ErrResponse
	}
	candidate := envelope.Candidates[0]
	if candidate.FinishReason != "STOP" {
		return nil, ErrResponse
	}
	var output strings.Builder
	for _, part := range candidate.Content.Parts {
		if !part.Thought {
			output.WriteString(part.Text)
		}
	}
	if output.Len() == 0 {
		return nil, ErrResponse
	}
	return []byte(output.String()), nil
}
