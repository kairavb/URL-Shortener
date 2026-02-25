package shortener

import (
	"url-shortener/internal/storage"
)

type Service struct {
	generator *Generator
	store     storage.Store
}

func NewService(gen *Generator, store storage.Store) *Service {
	return &Service{
		generator: gen,
		store:     store,
	}
}

func (s *Service) CreateShortURL(longURL string) (string, error) {
	shortCode := s.generator.Generate()

	err := s.store.Save(shortCode, longURL)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func (s *Service) Resolve(shortCode string) (string, error) {
	return s.store.Get(shortCode)
}
