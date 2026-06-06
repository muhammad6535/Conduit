package provider

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type MockProvider struct {
	name string
}

func NewMock(name string) *MockProvider {
	return &MockProvider{name: name}
}

func (p *MockProvider) Name() string {
	return p.name
}

func (p *MockProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(time.Duration(100+rand.Intn(200)) * time.Millisecond):
	}

	content := ""
	switch {
		case len(req.Messages) == 0:
			content = "Hello! I'm Conduit's demo mode. Send me a message!"
		default:
			lastMsg := req.Messages[len(req.Messages)-1].Content
			content = fmt.Sprintf(
				`**Conduit Demo Response** (mock mode)

You said: "%s"

**What happened behind the scenes:**
1. **Classified** your request as "simple chat"
2. **Routed** to cheapest capable model (%s → %s)
3. **Checked cache** → miss (first time seeing this query)
4. **Cost tracked** → $0.0001 (saved $0.0099 vs GPT-4o)

> This is a mock response. Connect real API keys to see actual completions.`,
				truncate(lastMsg, 80), req.Model, req.Model)

			if req.Stream {
				content = fmt.Sprintf(
					`**Streaming Demo Response**

Processing "%s"...

- Task complexity: simple
- Recommended model: gpt-4o-mini
- Cache: miss ✓
- Cost: $0.0000 (simulated)

Connect your API keys to go live!`,
					truncate(lastMsg, 40))
			}
	}

	randTokens := 10 + rand.Intn(50)
	return &CompletionResponse{
		ID:      fmt.Sprintf("mock-%d", time.Now().UnixNano()),
		Model:   req.Model,
		Created: time.Now().Unix(),
		Choices: []Choice{
			{
				Index:        0,
				Message:      Message{Role: RoleAssistant, Content: content},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10 + rand.Intn(30),
			CompletionTokens: randTokens,
			TotalTokens:      10 + rand.Intn(30) + randTokens,
		},
		Provider: p.name,
	}, nil
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}
