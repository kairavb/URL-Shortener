package storage

import "errors"

var ErrNotFound = errors.New("url not found")

// Store defines the interface for a URL shortening service's storage mechanism.
// It provides methods to save a mapping between a short code and a long URL, and to retrieve the long URL based on the short code.
type Store interface {
	Save(shortCode string, longURL string) error
	Get(shortCode string) (string, error)
	GetMaxID() (uint64, error)
}
