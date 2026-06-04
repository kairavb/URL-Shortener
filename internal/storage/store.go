package storage

import "errors"

var ErrNotFound = errors.New("url not found")

// Store is the persistence layer — SQLite today, swappable later.
type Store interface {
	Save(shortCode, longURL string) error
	Get(shortCode string) (string, error)
	GetMaxID() (uint64, error)
	Ping() error
}
