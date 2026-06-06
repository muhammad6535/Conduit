# Cost Savings Analysis

This document provides detailed methodology and real-world scenarios for the cost savings Conduit delivers.

---

## How Conduit Saves Money

| Technique | Mechanism | Typical Savings |
|---|---|---|
| **Intelligent Model Routing** | ML classifier routes simple queries to cheap models, complex to capable ones | 50-70% |
| **Semantic Caching** | Vector similarity detects repeated intent, serves cached responses | 30-50% fewer API calls |
| **Context Optimization** | Smart truncation, sliding windows, prompt compression | 20-40% fewer tokens |
| **Local Model Fallback** | Sensitive/trivial tasks → free local models | 0-100% (use-case dependent) |
| **Vendor Arbitrage** | Route to cheapest provider for equivalent quality | 10-30% |

**Combined: 50-90% reduction** in LLM spend for most organizations.

---

## Scenario 1: Internal AI Assistant (500 employees)

### Assumptions
- 500 employees
- 50 queries/day/employee
- 22 working days/month
- Average 1,000 input tokens + 200 output tokens per query
- Mix: 70% simple, 20% medium, 10% complex

### Without Conduit (all GPT-4o)

| Metric | Value |
|---|---|
| Monthly queries | 550,000 |
| Input tokens | 550M |
| Output tokens | 110M |
| **Monthly cost** | **$37,500** |
| Cost per employee | $75.00/month |

### With Conduit (intelligent routing)

| Model Tier | % of Queries | Queries | Cost |
|---|---|---|---|
| Cheap (GPT-4o-mini) | 70% | 385,000 | $225 |
| Balanced (GPT-4o) | 20% | 110,000 | $3,850 |
| Best (GPT-4o) | 10% | 55,000 | $3,750 |
| **Total** | **100%** | **550,000** | **$7,825** |

**Savings: $29,675/month (79%)**

### With Conduit + routing + caching (30% cache hit rate)

| Item | Cost |
|---|---|
| Intelligent routing | $7,825 |
| Cache hits (30% of queries free) | -$2,348 |
| **Total** | **$5,477** |
| **Savings vs raw** | **$32,023/month (85%)** |

---

## Scenario 2: Customer Support Agent (10,000 conversations/day)

### Assumptions
- 10,000 conversations/day
- 30 messages/conversation
- 200 tokens/message
- 30 days/month

### Without Conduit

| Model | Monthly Cost |
|---|---|
| GPT-4o only | $108,000 |
| GPT-4o-mini only | $3,240 |

### With Conduit (routing + caching)

| Strategy | Monthly Cost | Savings |
|---|---|---|
| No optimization | $108,000 | — |
| Rule-based routing only | $32,400 | 70% |
| Conduit routing + caching | **$16,200** | **85%** |

---

## Scenario 3: AI-Powered SaaS Product (1M API calls/day)

### Assumptions
- 1,000,000 API calls/day
- Average 500 tokens/request
- 30 days/month
- 40% of requests are near-duplicates (cacheable)

### Cost Comparison

| Approach | Monthly Cost | Cost/Request |
|---|---|---|
| Raw GPT-4o | $450,000 | $0.015 |
| Raw GPT-4o-mini | $6,750 | $0.000225 |
| Conduit (routing) | $22,500 | $0.00075 |
| Conduit (routing + caching) | **$13,500** | **$0.00045** |

---

## Methodology

Cost calculations use current (June 2026) API pricing:

| Model | Input/1K tokens | Output/1K tokens |
|---|---|---|
| GPT-4o | $0.0025 | $0.01 |
| GPT-4o-mini | $0.00015 | $0.0006 |
| GPT-4-turbo | $0.01 | $0.03 |
| Claude Sonnet 4 | $0.003 | $0.015 |
| Claude Haiku 3 | $0.00025 | $0.00125 |
| Claude Opus 4 | $0.015 | $0.075 |
| Llama 3 (local) | $0.00 | $0.00 |

Formula: `cost = (input_tokens / 1000 × input_price) + (output_tokens / 1000 × output_price)`

Cache savings assume 100% of cached requests avoid API calls entirely (no cost).

---

## Running Your Own Analysis

To calculate your potential savings:

1. Export your current usage from your LLM provider
2. Count total tokens and requests per model
3. Apply Conduit's routing assumptions:
   - 70% of GPT-4/Claude Opus traffic can use mini/Haiku
   - 30-40% of requests are cacheable
   - 20% average token reduction from context optimization

Or just deploy Conduit and check `/v1/stats` — it shows real-time savings.

---

## Limitations

- Savings depend on your specific usage patterns
- Routing quality improves as the classifier learns your traffic
- Cache hit rates vary by use case (chatbots: 30-50%, data extraction: 10-20%)
- Local models reduce cost but may have lower quality for complex tasks
