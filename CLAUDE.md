# Conduit — AI Development Guide for Claude

## Project Overview

Conduit is an open-source AI cost-optimization gateway. It intercepts LLM API calls and optimizes them for cost through intelligent routing, semantic caching, and privacy filtering.

- **Go Gateway**: Core proxy (single binary, <1ms overhead)
- **Python Engine**: ML sidecar for classification and caching
- **Apache 2.0** license
- **Stage**: Alpha (Phase 1 complete, Phase 2 starting)

## Repository Structure

```
conduit/
├── cmd/
│   ├── gateway/main.go          # Go entry point, CLI, server setup
│   └── engine/main.py           # Python FastAPI entry point
├── internal/
│   ├── config/config.go         # YAML config with env var expansion
│   ├── cost/tracker.go          # Per-request cost calculation
│   ├── provider/
│   │   ├── provider.go          # Provider interface
│   │   ├── openai.go            # OpenAI adapter
│   │   └── anthropic.go         # Anthropic adapter (format conversion)
│   ├── router/router.go         # Rule-based model routing
│   └── server/
│       ├── handler.go           # HTTP handlers, stats endpoint
│       └── stream.go            # SSE streaming support
├── internal/engine/             # Python ML modules
│   ├── classifier/router.py     # Task complexity classification
│   ├── cache/store.py           # Semantic vector cache
│   ├── embeddings/generator.py  # Embedding generation
│   ├── privacy/scanner.py       # PII detection & redaction
│   └── optimizer/prompt.py      # Token optimization
├── configs/                     # Pre-built config profiles
├── deploy/                      # Docker, Helm
├── scripts/                     # Quickstart, benchmarks
└── docs/brand/                  # SVG logo, branding
```

## Key Design Decisions

1. **Go for the proxy** — single binary, sub-ms overhead, excellent concurrency
2. **Python for ML** — best ecosystem for embeddings, NER, classification
3. **OpenAI-compatible API** — drop-in replacement, no code changes needed
4. **Cost-first design** — every feature tied to measurable cost reduction
5. **Apache 2.0** — corporate-friendly, patent protection

## Unique Value Proposition

Unlike LiteLLM (multi-provider proxy) or Helicone (observability), Conduit is the only OSS gateway built specifically for cost optimization:
- ML-based task classification (not just rule-based routing)
- Semantic caching (understands query intent, not just exact text)
- Built-in PII redaction (not enterprise-only)
- First-class local model support (Ollama, vLLM)

## Common Tasks

### Adding a provider adapter
1. Create `internal/provider/<name>.go` implementing `Provider` interface
2. Add provider case in `main.go`
3. Add model pricing to `cost/tracker.go`
4. Add example to `config.example.yaml`

### Adding a new endpoint
1. Add handler in `internal/server/handler.go` or create new file
2. Register route in `main.go`

### Working with config
- Config is YAML in `internal/config/config.go`
- Environment variables auto-expanded: `api_key: "${OPENAI_API_KEY}"`
- Hot reload via SIGHUP (planned)

## Building & Testing

```bash
# Build
go build -o build/conduit ./cmd/gateway/

# Test
go test ./... -v -count=1

# Lint
go vet ./...

# Run
./build/conduit --config config.yaml
./build/conduit --config config.test.yaml  # Test config (mock keys)
./build/conduit --init                      # Generate config
./build/conduit --version                   # Show version
```

## Current Gaps (Great for Contributors)

- **Integration tests** with real API mocking
- **More provider adapters** (Google Gemini, AWS Bedrock, Mistral API, Groq)
- **Python ML engine → Go gateway integration** (wire up classifier)
- **pgvector support** for persistent semantic cache
- **Dashboard** (React/Next.js)
- **Benchmark suite** with published numbers
- **Helm chart** for Kubernetes deployment
- **Rate limiting** (per-key, per-IP)
- **Admin API** for dynamic config changes
