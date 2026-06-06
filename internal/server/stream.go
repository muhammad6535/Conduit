package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/conduit-ai/gateway/internal/provider"
	"github.com/conduit-ai/gateway/internal/router"
)

func (h *Handler) streamChatCompletion(w http.ResponseWriter, r *http.Request, req *chatRequest, route *router.Route) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		jsonError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	compReq := &provider.CompletionRequest{
		Model:    route.Model,
		Messages: req.Messages,
		Stream:   true,
	}
	if req.MaxTokens > 0 {
		compReq.MaxTokens = req.MaxTokens
	}

	switch p := route.Provider.(type) {
	case *provider.OpenAIProvider:
		h.streamOpenAI(w, flusher, r.Context(), p, compReq)
	case *provider.AnthropicProvider:
		h.streamAnthropic(w, flusher, r.Context(), p, compReq)
	default:
		jsonError(w, http.StatusInternalServerError, "provider does not support streaming")
	}
}

func (h *Handler) streamOpenAI(w http.ResponseWriter, f http.Flusher, ctx interface{ Done() <-chan struct{} }, p *provider.OpenAIProvider, req *provider.CompletionRequest) {
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\":\"marshal error\"}\n\n")
		f.Flush()
		return
	}

	httpReq, err := http.NewRequestWithContext(nil, "POST", p.BaseURL()+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\":\"request error\"}\n\n")
		f.Flush()
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey())

	resp, err := p.HTTPClient().Do(httpReq)
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\":\"%s\"}\n\n", err.Error())
		f.Flush()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(w, "data: %s\n\n", string(respBody))
		f.Flush()
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if _, err := fmt.Fprintf(w, "%s\n\n", line); err != nil {
			break
		}
		f.Flush()
	}
	fmt.Fprintf(w, "\n")
	f.Flush()
}

func (h *Handler) streamAnthropic(w http.ResponseWriter, f http.Flusher, ctx interface{ Done() <-chan struct{} }, p *provider.AnthropicProvider, req *provider.CompletionRequest) {
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\":\"marshal error\"}\n\n")
		f.Flush()
		return
	}

	httpReq, err := http.NewRequestWithContext(nil, "POST", p.BaseURL()+"/messages", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\":\"request error\"}\n\n")
		f.Flush()
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.APIKey())
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.HTTPClient().Do(httpReq)
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\":\"%s\"}\n\n", err.Error())
		f.Flush()
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			oaiLine := convertAnthropicSSEToOpenAI(line)
			if oaiLine != "" {
				if _, err := fmt.Fprintf(w, "%s\n\n", oaiLine); err != nil {
					break
				}
				f.Flush()
			}
		}
	}
	fmt.Fprintf(w, "\n")
	f.Flush()
}

func convertAnthropicSSEToOpenAI(anthropicLine string) string {
	data := strings.TrimPrefix(anthropicLine, "data: ")
	data = strings.TrimSpace(data)

	if data == "[DONE]" {
		return "data: [DONE]"
	}

	var event struct {
		Type  string `json:"type"`
		Delta struct {
			Text string `json:"text"`
		} `json:"delta"`
		Index int `json:"index"`
	}

	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return ""
	}

	switch event.Type {
	case "content_block_delta":
		if event.Delta.Text == "" {
			return ""
		}
		oaiChunk := map[string]interface{}{
			"id":      "chunk",
			"object":  "chat.completion.chunk",
			"created": 0,
			"model":   "",
			"choices": []map[string]interface{}{
				{
					"index": event.Index,
					"delta": map[string]interface{}{
						"content": event.Delta.Text,
						"role":    "assistant",
					},
				},
			},
		}
		oaiBytes, _ := json.Marshal(oaiChunk)
		return "data: " + string(oaiBytes)

	case "message_stop":
		return "data: [DONE]"

	default:
		return ""
	}
}
