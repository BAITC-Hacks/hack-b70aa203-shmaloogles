package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/shmaloogles/business-task-platform/backend/internal/ai"
	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
)

// RegisterAIRoutes adds stateless AI routes without reading or writing tasks.
func RegisterAIRoutes(mux *http.ServeMux, service *ai.Service) {
	mux.HandleFunc("POST /api/tasks/clarify", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Description string `json:"description"`
		}
		if !decodeAIRequest(w, r, &input) {
			return
		}
		result, err := service.Clarify(r.Context(), input.Description)
		if err != nil {
			writeAIError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("POST /api/tasks/generate", func(w http.ResponseWriter, r *http.Request) {
		var input ai.GenerateInput
		if !decodeAIRequest(w, r, &input) {
			return
		}
		result, err := service.Generate(r.Context(), input)
		if err != nil {
			writeAIError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, struct {
			ai.Generation
			Readiness scoring.Result `json:"readiness"`
		}{Generation: result, Readiness: scoring.Calculate(result.Card)})
	})
}

func decodeAIRequest(w http.ResponseWriter, r *http.Request, value any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := decodeJSON(r, value); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds 1 MiB")
		} else {
			writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON with known fields")
		}
		return false
	}
	return true
}

func writeAIError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ai.ErrInput):
		writeError(w, http.StatusBadRequest, "invalid_input", err.Error())
	case errors.Is(err, context.Canceled):
		// The caller is gone; do not replace cancellation with a mock success.
		return
	case errors.Is(err, context.DeadlineExceeded):
		writeError(w, http.StatusGatewayTimeout, "ai_timeout", "AI request timed out")
	case errors.Is(err, ai.ErrResponse):
		writeError(w, http.StatusBadGateway, "invalid_ai_response", "AI returned invalid task data")
	default:
		writeError(w, http.StatusBadGateway, "ai_unavailable", "AI provider unavailable")
	}
}
