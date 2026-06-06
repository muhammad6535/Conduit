#!/bin/bash
# Conduit Benchmark — measures latency and throughput
set -e

GATEWAY_URL="${GATEWAY_URL:-http://localhost:8080}"
CONCURRENCY="${CONCURRENCY:-10}"
REQUESTS="${REQUESTS:-50}"

echo "═══ Conduit Benchmark ═══"
echo "Target:    $GATEWAY_URL"
echo "Requests:  $REQUESTS"
echo "Concurrency: $CONCURRENCY"
echo ""

# Health check
echo -n "• Health: "
curl -s -o /dev/null -w "%{http_code}" "$GATEWAY_URL/health"
echo ""

# Sequential latency test
echo "• Latency (sequential, $REQUESTS requests):"
total=0
for i in $(seq 1 $REQUESTS); do
    start=$(date +%s%N)
    curl -s -X POST "$GATEWAY_URL/v1/chat/completions" \
        -H "Content-Type: application/json" \
        -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"Hi"}]}' \
        -o /dev/null -w "%{http_code}" > /dev/null 2>&1 || true
    end=$(date +%s%N)
    elapsed=$(( (end - start) / 1000000 ))
    total=$((total + elapsed))
done
avg=$((total / REQUESTS))
echo "  Average: ${avg}ms"

# Throughput test (requires parallel)
if command -v parallel &> /dev/null; then
    echo "• Throughput (parallel, $CONCURRENCY concurrent):"
    seq $REQUESTS | parallel -j $CONCURRENCY --progress \
        "curl -s -X POST '$GATEWAY_URL/v1/chat/completions' \
            -H 'Content-Type: application/json' \
            -d '{\"model\":\"gpt-4o-mini\",\"messages\":[{\"role\":\"user\",\"content\":\"Hi\"}]}' \
            -o /dev/null -w '%{http_code}\n'" > /dev/null 2>&1
    echo "  Completed $REQUESTS requests with $CONCURRENCY concurrency"
fi

# Stats
echo ""
echo "• Gateway Stats:"
curl -s "$GATEWAY_URL/v1/stats" | python -m json.tool 2>/dev/null || curl -s "$GATEWAY_URL/v1/stats"
