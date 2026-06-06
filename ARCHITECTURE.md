# Conduit Architecture

Conduit is a two-component system: a **high-performance Go gateway** (the proxy) and an **optional Python ML engine** (the brain).

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     External Clients                             │
│  (OpenAI SDK, LangChain, cURL, Vercel AI SDK, Cursor, etc.)     │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────┐
│                  Conduit Gateway (Go)                            │
│                                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────────┐  │
│  │  HTTP     │  │  Auth/   │  │  Router  │  │  Provider      │  │
│  │  Server   │──│  Rate    │──│  Engine  │──│  Adapters      │  │
│  │  (chi)    │  │  Limit   │  │          │  │                │  │
│  └──────────┘  └──────────┘  └────┬─────┘  │  ┌──────────┐  │  │
│                                    │         │  │ OpenAI   │  │  │
│  ┌──────────┐  ┌──────────┐       │         │  ├──────────┤  │  │
│  │  Cost     │  │  Stats   │       │         │  │Anthropic │  │  │
│  │  Tracker  │──│  Handler │       │         │  ├──────────┤  │  │
│  │          │  │          │       │         │  │ Ollama   │  │  │
│  └──────────┘  └──────────┘       │         │  └──────────┘  │  │
│                                    │         └────────────────┘  │
│                           ┌───────┴───────┐                     │
│                           │  Cache Client  │                     │
│                           │  (HTTP to ML   │                     │
│                           │   Engine)      │                     │
│                           └───────┬───────┘                     │
└───────────────────────────────────┼──────────────────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────────┐
          │                         ▼                             │
          │          ┌──────────────────────────────┐             │
          │          │   Conduit ML Engine (Python)   │             │
          │          │                                │             │
          │          │  ┌──────────┐ ┌─────────────┐ │             │
          │          │  │Classifier│ │   Semantic   │ │             │
          │          │  │(FastAPI) │ │  Cache      │ │             │
          │          │  └──────────┘ │  (pgvector)  │ │             │
          │          │               └─────────────┘ │             │
          │          │  ┌──────────┐ ┌─────────────┐ │             │
          │          │  │  PII     │ │   Embedding  │ │             │
          │          │  │ Scanner  │ │   Generator  │ │             │
          │          │  └──────────┘ └─────────────┘ │             │
          │          └──────────────────────────────┘             │
          │                                                         │
          ▼                                                         ▼
   ┌──────────────┐                                       ┌──────────────┐
   │  LLM APIs    │                                       │  Local Models│
   │  (Cloud)     │                                       │  (Ollama,     │
   │              │                                       │   vLLM, etc.) │
   └──────────────┘                                       └──────────────┘
```

## Component Breakdown

### 1. Gateway (Go)

**Purpose**: High-performance proxy that intercepts all LLM API calls.

**Why Go**: Single binary deployment, sub-millisecond overhead, excellent concurrency for handling thousands of simultaneous requests, minimal resource footprint (~15MB binary, ~10MB RAM idle).

**Key packages**:

| Package | Responsibility |
|---|---|
| `internal/config` | YAML config loading, env var expansion, validation |
| `internal/server` | HTTP handlers, middleware (auth, logging, recovery, request ID) |
| `internal/router` | Provider selection based on model name patterns and ML classification |
| `internal/provider` | Provider interface + adapters for OpenAI, Anthropic, local models |
| `internal/cost` | Real-time cost tracking with pricing database |

**Request flow**:
1. HTTP request arrives at chi router
2. Middleware: logging, recovery, request ID
3. Handler parses JSON body
4. Router selects provider based on model name
5. Provider adapter sends request to upstream API
6. Response is parsed and returned in OpenAI-compatible format
7. Cost tracker records tokens and calculates cost
8. Response headers include cost and provider info

### 2. ML Engine (Python)

**Purpose**: Intelligence layer for routing decisions, semantic caching, and privacy.

**Why Python**: Rich ecosystem for ML (sentence-transformers, scikit-learn, spaCy), easy to extend with new classification models.

**Key modules**:

| Module | Responsibility |
|---|---|
| `classifier/` | Task complexity detection (simple, medium, complex) |
| `embeddings/` | Vector embeddings for semantic cache keys |
| `cache/` | Semantic cache with cosine similarity search |
| `privacy/` | PII detection and redaction |
| `optimizer/` | Prompt compression and token optimization |

### 3. Dashboard (Planned, React/Next.js)

**Purpose**: Visual analytics and management interface.

**Features planned**:
- Real-time usage and cost dashboards
- Cost breakdowns by model, provider, user, team
- Cache hit rate and savings
- Alert configuration
- Config editor

## Data Flow (Complete Request Lifecycle)

```
Client                    Conduit                    ML Engine              Provider
  │                         │                           │                     │
  │── POST /chat/completions│                           │                     │
  │                         │── POST /classify ────────▶│                     │
  │                         │                           │── analyze text      │
  │                         │◀──── complexity ──────────│                     │
  │                         │                           │                     │
  │                         │── POST /cache/lookup ────▶│                     │
  │                         │◀── miss or hit ───────────│                     │
  │                         │                           │                     │
  │                         │ (if miss)                 │                     │
  │                         │── POST /privacy/scan ────▶│                     │
  │                         │◀── redacted text ─────────│                     │
  │                         │                           │                     │
  │                         │───────────────────────────│── POST /chat/ ...  │
  │                         │                           │◀──── response ─────│
  │                         │                           │                     │
  │                         │── POST /cache/store ─────▶│                     │
  │                         │                           │                     │
  │◀─── response + headers──│                           │                     │
  │  X-Conduit-Cost         │                           │                     │
  │  X-Conduit-Provider     │                           │                     │
  │  X-Conduit-Cache-Hit    │                           │                     │
```

## Design Decisions

### Why not Python-only?
Python is 10-50x slower for HTTP proxy workloads and requires a runtime environment. Go gives us a single ~15MB binary with <1ms overhead. The Python ML engine runs as an optional sidecar for intelligence features.

### Why not Rust?
Rust would be faster but has a steeper learning curve that would limit community contributions. Go strikes the best balance of performance, simplicity, and contributor accessibility.

### Why not integrate ML into Go directly?
Go's ML ecosystem is immature. Python has the best tools for embeddings, classification, and NLP. Running Python as a sidecar keeps components decoupled and independently deployable.

### Why semantic caching over exact-match?
Exact-match caching misses most opportunities. Users rarely ask the exact same question twice. Semantic caching captures similar intent, catching 3-5x more cache hits.

## Configuration Philosophy

Configuration is YAML-based with environment variable expansion:
- `api_key: "${OPENAI_API_KEY}"` — auto-expanded from environment
- Hot-reload support planned (SIGHUP)
- Validation on load catches misconfiguration early

## Performance Targets

| Metric | Target |
|---|---|
| P50 latency overhead | <1ms |
| P95 latency overhead | <3ms |
| Throughput (single instance) | 10,000 req/s |
| Memory (idle) | <15MB |
| Binary size | <20MB |

## Security Model

- API keys stored in environment variables, never in config files committed to git
- PII redacted before reaching external providers
- Rate limiting prevents abuse
- Request logging without sensitive content (planned)

## Deployment

### Minimal (single binary)
```
./conduit --config config.yaml
```

### Docker
```
docker run conduit
```

### Docker Compose (with ML engine)
```
docker compose -f deploy/docker-compose.yml up
```

### Kubernetes (planned)
Helm chart with:
- Horizontal pod autoscaling
- Prometheus metrics
- ConfigMaps for configuration
- Ingress support

---

*For implementation details, see individual package documentation and source code.*
