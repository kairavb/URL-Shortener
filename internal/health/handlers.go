package health

import (
	"encoding/json"
	"net/http"

	"url-shortener/internal/storage"
)

// Health answers "is the process running?" — used by load balancers for liveness.
func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready answers "can this instance serve traffic?" — checks DB connectivity.
func Ready(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if err := store.Ping(); err != nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
