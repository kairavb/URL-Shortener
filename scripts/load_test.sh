#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
REQUESTS="${REQUESTS:-5000}"
CONCURRENCY="${CONCURRENCY:-50}"

if ! command -v hey >/dev/null 2>&1; then
  echo "hey is required: go install github.com/rakyll/hey@latest"
  exit 1
fi

echo "==> Creating a test short URL"
RESPONSE=$(curl -sf -X POST "$BASE_URL/shorten" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com/load-test"}')
SHORT_CODE=$(echo "$RESPONSE" | sed -n 's/.*"short_code":"\([^"]*\)".*/\1/p')

echo "==> Redirect load test (cache hot path)"
hey -n "$REQUESTS" -c "$CONCURRENCY" "$BASE_URL/$SHORT_CODE"

echo
echo "==> Shorten load test (write path)"
hey -n 500 -c 20 -m POST \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com/bench"}' \
  "$BASE_URL/shorten"
