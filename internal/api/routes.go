package api

import (
	"net/http"

	"url-shortener/internal/health"
	"url-shortener/internal/storage"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, store storage.Store) {
	// Register fixed paths before /{shortCode} so they are not treated as codes.
	mux.HandleFunc("GET /health", health.Health)
	mux.HandleFunc("GET /ready", health.Ready(store))
	mux.HandleFunc("POST /shorten", handler.Shorten)
	mux.HandleFunc("GET /", handler.ServeUI)
	mux.HandleFunc("GET /{shortCode}", handler.Resolve)
}
