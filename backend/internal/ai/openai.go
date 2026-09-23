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

type openAI struct {
	endpoint string
	key      string
	model    string
	client   *http.Client
}

// FromEnv defaults to labelled mock mode. Misconfigured real mode fails at startup.
func FromEnv() (*Service, error) {
	mode := strings.TrimSpace(os.Getenv("AI_MODE"))
	if mode == "" || mode == "mock" {
		return New(nil, 0, false), nil
	}
	if mode != "openai" {
		return nil, fmt.Errorf("AI_MODE must be mock or openai")
	}
	key, model := strings.TrimSpace(os.Getenv("OPENAI_API_KEY")), strings.TrimSpace(os.Getenv("AI_MODEL"))
	if model == "" {
		model = "gpt-4.1-mini"
	}
	if key == "" || !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(model) {
		return nil, fmt.Errorf("openai requires OPENAI_API_KEY and a valid AI_MODEL")
	}
	base := strings.TrimRight(os.Getenv("AI_BASE_URL"), "/")
	if base == "" {
		base = "https://api.openai.com/v1"
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
	return New(&openAI{endpoint: base + "/responses", key: key, model: model, client: client}, timeout, fallback), nil
}

func (g *openAI) Complete(ctx context.Context, prompt string, input, schema json.RawMessage) ([]byte, error) {
	body, err := json.Marshal(map[string]any{
		"model":             g.model,
		"instructions":      prompt,
		"input":             string(input),
		"store":             false,
		"max_output_tokens": 4096,
		"text": map[string]any{"format": map[string]any{
			"type": "json_schema", "name": "task_response", "strict": true, "schema": schema,
		}},
	})
	if err != nil {
		return nil, ErrProvider
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, ErrProvider
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.key)
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
		Status string `json:"status"`
		Output []struct {
			Type    string `json:"type"`
			Status  string `json:"status"`
			Role    string `json:"role"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if json.Unmarshal(raw, &envelope) != nil || envelope.Status != "completed" {
		return nil, ErrResponse
	}
	var output strings.Builder
	messages := 0
	for _, item := range envelope.Output {
		if item.Type == "reasoning" {
			continue
		}
		if item.Type != "message" || item.Status != "completed" || item.Role != "assistant" {
			return nil, ErrResponse
		}
		messages++
		for _, part := range item.Content {
			// Refusals and unexpected content are never accepted as task data.
			if part.Type != "output_text" {
				return nil, ErrResponse
			}
			output.WriteString(part.Text)
		}
	}
	if messages != 1 || output.Len() == 0 {
		return nil, ErrResponse
	}
	return []byte(output.String()), nil
}
