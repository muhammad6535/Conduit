# Conduit AI Assistant Guide

Conduit is an open-source AI cost-optimization gateway written primarily in Go, with a Python ML engine sidecar.

## Quick Facts

- **Gateway**: Go, single binary, <1ms overhead, Apache 2.0
- **ML Engine**: Python FastAPI sidecar for classification, caching, privacy
- **Config**: YAML-based with `${ENV_VAR}` expansion
- **API**: OpenAI-compatible (`/v1/chat/completions`)
- **Target**: Cost reduction (50-80% for most orgs)
- **License**: Apache 2.0

## Architecture

see [ARCHITECTURE.md](ARCHITECTURE.md)

Key directories:
- `cmd/gateway/main.go` — Entry point
- `internal/server/` — HTTP handlers, middleware
- `internal/provider/` — Provider adapters (OpenAI, Anthropic)
- `internal/cost/` — Cost tracking
- `internal/router/` — Request routing
- `internal/config/` — Config loading
- `cmd/engine/main.py` — ML engine entry
- `internal/engine/` — Classifier, cache, privacy modules

## Commands

```bash
go build -o build/conduit ./cmd/gateway/
go test ./...
go vet ./...
./build/conduit --config config.yaml
./build/conduit --init      # Generate config
./build/conduit --version   # Show version
```

## Coding Standards

- Follow standard Go conventions (`gofmt`, `go vet`)
- Python: follow PEP 8, use `ruff` for linting
- No comments in code unless explaining complex logic
- Error handling: return errors, don't panic
- Logging: use `log.Printf` for server-side, `fmt.Printf` for CLI

## Adding a New Provider

1. Create `internal/provider/<name>.go`
2. Implement `provider.Provider` interface
3. Add to `cmd/gateway/main.go` switch statement
4. Add pricing to `internal/cost/tracker.go`
5. Update config example and README

## Testing

- Go tests live next to source (`*_test.go`)
- Python tests use pytest
- Integration tests start the server and make HTTP requests
