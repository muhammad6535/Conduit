---
layout: default
title: Conduit — AI Cost-Optimization Gateway
description: Open-source AI gateway that cuts LLM costs 50-80% through intelligent routing, semantic caching, and privacy filtering.
---

# Conduit

**The AI Gateway That Pays for Itself**

Open-source. Self-hosted. Cuts LLM costs 50-80%. Intelligently routes, caches, and optimizes every API call.

## Quick Start

```bash
docker run -p 8080:8080 \
  -e OPENAI_API_KEY="sk-..." \
  ghcr.io/muhammad6535/conduit:latest
```

Connect your app by changing one line:

```python
client = OpenAI(base_url="http://localhost:8080/v1")
```

## Features

| Feature | Description |
|---|---|
| 🎯 **ML-based routing** | Routes to cheapest capable model based on task complexity |
| 💾 **Semantic caching** | Vector similarity cache — understands intent, not just exact text |
| 🛡️ **PII redaction** | Built-in privacy filter strips emails, phones, SSNs, API keys |
| 📊 **Real-time cost tracking** | Per-request cost, monthly totals, savings reports |
| ⚡ **Streaming** | Full SSE support with format conversion |

## Architecture

```
Your App → Conduit Gateway (Go) → ML Classifier → Router → Provider
                                        ↓
                                  Semantic Cache (Python)
```

## Why Conduit?

- **Cost-first design**: Every feature directly reduces LLM spend
- **Go core**: ~15MB binary, <1ms overhead, 10k+ req/s
- **OpenAI-compatible**: Zero code changes — works with existing SDKs
- **Self-hosted**: Full control of your data, no vendor lock-in
- **Apache 2.0**: Corporate-friendly, patent protection

## Links

- [GitHub Repository](https://github.com/muhammad6535/conduit)
- [ARCHITECTURE.md](ARCHITECTURE.md)
- [COST_SAVINGS.md](COST_SAVINGS.md)
- [CONTRIBUTING.md](CONTRIBUTING.md)
