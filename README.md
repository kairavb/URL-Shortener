# Scalable URL Shortening Service

A Go backend that shortens URLs, redirects with caching, and stays safe under concurrent load.

## Run

```bash
go run ./cmd
```

Open `http://localhost:8080` for the UI, or use curl:

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

## Architecture

```
Client (HTML / curl)
        ↓
HTTP API (net/http)
        ↓
Rate Limiter → Cache → DB
        ↓
Short Code Generator (atomic-safe)

Data flow model:
Request
  -> Rate limiter (drop if too fast)
  -> Logger
  -> Handler
      POST /shorten -> validate URL -> generate code -> save DB -> cache it
      GET /{code}   -> check cache -> else DB -> cache it -> redirect
```

## Features

- REST APIs for shorten and redirect
- Collision-safe Base62 short codes (atomic counter)
- SQLite persistence with restart-safe ID recovery
- Token-bucket rate limiting (per client IP)
- In-memory LRU cache for hot redirects
- Health and readiness probes
- Load testing script with `hey`
- Graceful shutdown

Mainly:
No collisions -> atomic counter + Base62, not random strings.
Restart safe -> on startup, max ID is read from DB so codes never repeat.
Cache -> redirects are read-heavy; LRU keeps hot URLs out of SQLite.
Rate limit -> token bucket per IP stops abuse without blocking normal traffic.


## Example Flow
### Flow 1: Shortening a URL (POST /shorten):
- Client sends: {"url": "https://example.com"}
- Rate limiter checks if the IP is sending too many requests.
- Handler parses the JSON.
- Service validates the URL (must be http or https).
- Generator creates a unique short code using an atomic counter + Base62 (a, b, … 9, ba, …).
- Storage saves short_code → long_url in urls.db.
- Cache stores the mapping in memory.
- Response: {"short_code":"b","short_url":"http://localhost:8080/b"}

### Flow 2: Opening a short link (GET /b)
- User visits http://localhost:8080/b
- Handler extracts b from the path.
- Service checks cache first.
- If not cached, reads from urls.db and caches it.
- Handler sends HTTP redirect (307) to the long URL.
- Redirects are read-heavy, so the cache avoids hitting the DB on every click.

### Flow 3: Restart safety
- On startup in main.go:
- Open urls.db
- Find the last short code used
- Resume the counter from there
- So after restart you don’t reuse old short codes.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/shorten` | Create a short URL |
| GET | `/{shortCode}` | Redirect to the long URL |
| GET | `/health` | Liveness check |
| GET | `/ready` | Readiness check (includes DB ping) |
| GET | `/` | Simple web UI |

### Request / Response

```json
POST /shorten
{ "url": "https://example.com" }

Response
{
  "short_code": "b",
  "short_url": "http://localhost:8080/b"
}
```

## Load Testing

Install [hey](https://github.com/rakyll/hey):

```bash
go install github.com/rakyll/hey@latest
```

Run the script (server must be running):

```bash
chmod +x scripts/load_test.sh
./scripts/load_test.sh
```

Typical results on a local machine:

- **Redirect (cached):** ~15k–25k req/s, p99 under 10ms
- **Shorten (write):** limited by SQLite writes (~500–1500 req/s)

The LRU cache moves repeated redirects off the DB read path, which is the main hot-path optimization.

## Project Structure

```
url-shortener/
├── cmd/main.go                 # wires everything together
├── internal/
│   ├── api/                    # handlers, routes, logging middleware, HTTP layer
│   ├── cache/                  # LRU cache
│   ├── health/                 # /health and /ready
│   ├── ratelimit/              # token-bucket limiter
│   ├── shortener/              # ID generator and service logic
│   └── storage/                # Store interface and SQLite
├── scripts/load_test.sh        # benchmarks
├── web/index.html              # UI
└── README.md
```

## Tests

```bash
go test ./...
```

## License

MIT
