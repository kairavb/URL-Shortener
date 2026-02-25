package api

import (
	"encoding/json"
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
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	shortCode, err := h.service.CreateShortURL(req.URL)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{ShortCode: shortCode}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:] // remove leading /

	longURL, err := h.service.Resolve(shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}
