package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"url-shortener/internal/shortener"
)

type Handler struct {
	service *shortener.Service
}

func NewHandler(service *shortener.Service) *Handler {
	return &Handler{service: service}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	shortCode, err := h.service.CreateShortURL(req.URL)
	if errors.Is(err, shortener.ErrInvalidURL) {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{
		ShortCode: shortCode,
		ShortURL:  buildShortURL(r, shortCode),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode") // set by the GET /{shortCode} route pattern

	longURL, err := h.service.Resolve(shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func (h *Handler) ServeUI(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func buildShortURL(r *http.Request, shortCode string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s", scheme, r.Host, shortCode)
}
