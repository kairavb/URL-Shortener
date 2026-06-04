package shortener

import (
	"errors"
	"net/url"

	"url-shortener/internal/cache"
	"url-shortener/internal/storage"
)

var ErrInvalidURL = errors.New("invalid url")

// Service holds the core shorten/resolve logic shared by all HTTP handlers.
type Service struct {
	generator *Generator
	store     storage.Store
	cache     *cache.LRU
}

func NewService(gen *Generator, store storage.Store, cache *cache.LRU) *Service {
	return &Service{
		generator: gen,
		store:     store,
		cache:     cache,
	}
}

func (s *Service) CreateShortURL(longURL string) (string, error) {
	if err := validateURL(longURL); err != nil {
		return "", err
	}

	shortCode := s.generator.Generate()
	if err := s.store.Save(shortCode, longURL); err != nil {
		return "", err
	}

	s.cache.Put(shortCode, longURL)
	return shortCode, nil
}

func (s *Service) Resolve(shortCode string) (string, error) {
	// Cache-aside: check memory first, fall back to DB, then warm the cache.
	if longURL, ok := s.cache.Get(shortCode); ok {
		return longURL, nil
	}

	longURL, err := s.store.Get(shortCode)
	if err != nil {
		return "", err
	}

	s.cache.Put(shortCode, longURL)
	return longURL, nil
}

func validateURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ErrInvalidURL
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidURL
	}
	if parsed.Host == "" {
		return ErrInvalidURL
	}
	return nil
}
