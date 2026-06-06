package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/conduit-ai/gateway/internal/cost"
	"github.com/conduit-ai/gateway/internal/provider"
	"github.com/conduit-ai/gateway/internal/router"
)

type Handler struct {
	router  *router.Router
	tracker *cost.Tracker
}

func NewHandler(r *router.Router, t *cost.Tracker) *Handler {
	return &Handler{router: r, tracker: t}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	jsonResp(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) ListModels(w http.ResponseWriter, r *http.Request) {
	models := h.router.AvailableModels()
	jsonResp(w, http.StatusOK, map[string]interface{}{
		"object": "list",
		"data":   models,
	})
}

type chatRequest struct {
	Model       string             `json:"model"`
	Messages    []provider.Message `json:"messages"`
	Stream      bool               `json:"stream,omitempty"`
	MaxTokens   int                `json:"max_tokens,omitempty"`
	Temperature *float64           `json:"temperature,omitempty"`
}

func (h *Handler) ChatCompletions(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	if len(req.Messages) == 0 {
		jsonError(w, http.StatusBadRequest, "messages is required")
		return
	}

	route, err := h.router.Route(req.Model)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Sprintf("routing error: %v", err))
		return
	}

	log.Printf("[route] model=%q → provider=%q actual_model=%q",
		req.Model, route.Provider.Name(), route.Model)

	if req.Stream {
		h.streamChatCompletion(w, r, &req, route)
		return
	}

	compReq := &provider.CompletionRequest{
		Model:    route.Model,
		Messages: req.Messages,
		Stream:   false,
	}

	if req.MaxTokens > 0 {
		compReq.MaxTokens = req.MaxTokens
	}
	if req.Temperature != nil {
		compReq.Temperature = *req.Temperature
	}

	resp, err := route.Provider.Complete(r.Context(), compReq)
	if err != nil {
		jsonError(w, http.StatusBadGateway, fmt.Sprintf("provider error: %v", err))
		return
	}

	cost := h.tracker.Record(resp.Model, resp.Provider, resp.Usage.PromptTokens, resp.Usage.CompletionTokens)

	log.Printf("[cost] model=%s tokens=%d+%d cost=$%.6f monthly=$%.4f",
		resp.Model, resp.Usage.PromptTokens, resp.Usage.CompletionTokens,
		cost, h.tracker.MonthlyCost())

	w.Header().Set("X-Conduit-Provider", resp.Provider)
	w.Header().Set("X-Conduit-Cost", fmt.Sprintf("%.6f", cost))
	w.Header().Set("X-Conduit-Monthly-Cost", fmt.Sprintf("%.4f", h.tracker.MonthlyCost()))

	jsonResp(w, http.StatusOK, resp)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	records := h.tracker.Records()
	byModel := make(map[string]struct {
		Requests int
		Tokens   int
		Cost     float64
	})
	var totalTokens int
	var totalCost float64

	for _, rec := range records {
		m := byModel[rec.Model]
		m.Requests++
		m.Tokens += rec.TotalTokens
		m.Cost += rec.Cost
		byModel[rec.Model] = m
		totalTokens += rec.TotalTokens
		totalCost += rec.Cost
	}

	jsonResp(w, http.StatusOK, map[string]interface{}{
		"total_requests": len(records),
		"total_tokens":   totalTokens,
		"total_cost":     fmt.Sprintf("$%.4f", totalCost),
		"by_model":       byModel,
	})
}

func jsonResp(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResp(w, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"type":    strings.ReplaceAll(strings.ToLower(http.StatusText(status)), " ", "_"),
			"code":    status,
		},
	})
}
