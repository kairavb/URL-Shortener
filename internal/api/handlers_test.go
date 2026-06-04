package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/cache"
	"url-shortener/internal/shortener"
	"url-shortener/internal/storage"
)

type memoryStore struct {
	data map[string]string
}

func (m *memoryStore) Save(shortCode, longURL string) error {
	m.data[shortCode] = longURL
	return nil
}

func (m *memoryStore) Get(shortCode string) (string, error) {
	url, ok := m.data[shortCode]
	if !ok {
		return "", storage.ErrNotFound
	}
	return url, nil
}

func (m *memoryStore) GetMaxID() (uint64, error) { return 0, nil }
func (m *memoryStore) Ping() error               { return nil }

func TestShortenAndResolve(t *testing.T) {
	store := &memoryStore{data: make(map[string]string)}
	service := shortener.NewService(shortener.NewGenerator(0), store, cache.New(10))
	handler := NewHandler(service)

	mux := http.NewServeMux()
	RegisterRoutes(mux, handler, store)

	shortenReq := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(`{"url":"https://example.com"}`))
	shortenReq.Header.Set("Content-Type", "application/json")
	shortenRec := httptest.NewRecorder()
	mux.ServeHTTP(shortenRec, shortenReq)

	if shortenRec.Code != http.StatusOK {
		t.Fatalf("shorten status = %d", shortenRec.Code)
	}

	var resp shortenResponse
	if err := json.NewDecoder(shortenRec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ShortCode == "" || resp.ShortURL == "" {
		t.Fatalf("unexpected response: %+v", resp)
	}

	resolveReq := httptest.NewRequest(http.MethodGet, "/"+resp.ShortCode, nil)
	resolveRec := httptest.NewRecorder()
	mux.ServeHTTP(resolveRec, resolveReq)

	if resolveRec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("resolve status = %d", resolveRec.Code)
	}
	if got := resolveRec.Header().Get("Location"); got != "https://example.com" {
		t.Fatalf("location = %q", got)
	}
}

func TestHealthAndReady(t *testing.T) {
	store := &memoryStore{data: make(map[string]string)}
	handler := NewHandler(shortener.NewService(shortener.NewGenerator(0), store, cache.New(10)))

	mux := http.NewServeMux()
	RegisterRoutes(mux, handler, store)

	for _, path := range []string{"/health", "/ready"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d", path, rec.Code)
		}
	}
}
