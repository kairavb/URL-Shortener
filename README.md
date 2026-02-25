# Scalable URL Shortening Service

## Project Explanation

SOON

## Running the Project

SOON

## Flow

Client (HTML / curl)
↓
HTTP API (net/http)
↓
Rate Limiter → Cache → DB
↓
ShortCode Generator (atomic-safe)

## features

- URL shortening
- REST APIs
- Persistent storage
- Collision-safe short codes
- Concurrency handling
- Token-bucket rate limiting
- In-memory Caching
- Health / readiness checks
- Load / stress testing
- Performance optimization
- Metrics / Analytics

## Endpoints

- POST /shorten
- GET /{shortCode}
- GET /health
- GET /ready

## Request Response Model

```json
POST /shorten
{
  "url": "https://example.com"
}

Response
{
  "short_code": "aZ3f",
  "short_url": "http://localhost:8080/aZ3f"
}
```

## Project Structure

```bash
url-shortener/
├── cmd/
│   └── main.go              # App entrypoint
│
├── internal/
│   ├── api/
│   │   ├── handlers.go          # HTTP handlers
│   │   ├── middleware.go        # Rate limiting, logging
│   │   └── routes.go            # Route registration
│   │
│   ├── shortener/
│   │   ├── generator.go         # Base62, ID counter
│   │   └── generator_test.go
│   │
│   ├── storage/
│   │   ├── store.go             # Interface
│   │   ├── sqlite.go            # SQLite implementation
│   │   └── migrations.sql
│   │
│   ├── cache/
│   │   └── lru.go               # In-memory LRU cache
│   │
│   ├── ratelimit/
│   │   └── token_bucket.go
│   │
│   └── health/
│       └── handlers.go
│
├── web/
│   └── index.html               # Minimal UI
│
├── scripts/
│   └── load_test.sh             # wrk / hey
│
├── go.mod
├── go.sum
└── README.md
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License.

````md
// old readme file

# URL-Shortner

A GO project that shortens long URL's into small ones and redirects them

A part of 8 day challenge completed in just under 4 days, that is,

- learn GO from scratch
- learn a web framework of GO (Fiber)
- build a small Full-Stack App in GO

## Usage

```bash
git clone https://github.com/kairavb/URL-Shortner.git
go mod tidy
go run main.go
```

## Fiber Installation

```bash
git clone repo
go mod init project_name
go go get github.com/gofiber/fiber/v2
```
````
