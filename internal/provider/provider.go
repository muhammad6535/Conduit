package provider

import (
	"context"
	"encoding/json"
	"time"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Stream     bool      `json:"stream,omitempty"`
	MaxTokens  int       `json:"max_tokens,omitempty"`
	Temperature float64  `json:"temperature,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionResponse struct {
	ID        string    `json:"id"`
	Model     string    `json:"model"`
	Created   int64     `json:"created"`
	Choices   []Choice  `json:"choices"`
	Usage     Usage     `json:"usage"`
	Provider  string    `json:"-"`
}

func (r *CompletionResponse) MarshalJSON() ([]byte, error) {
	type Alias CompletionResponse
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}
	return json.Marshal((*Alias)(r))
}

type Provider interface {
	Name() string
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
}
