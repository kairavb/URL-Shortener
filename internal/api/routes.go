package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("/shorten", handler.Shorten)
	mux.HandleFunc("/", handler.Resolve)
}
