package storage

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

func (s *SQLiteStore) Init() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		short_code TEXT PRIMARY KEY,
		long_url TEXT NOT NULL
	);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) Save(shortCode string, longURL string) error {
	_, err := s.db.Exec(
		"INSERT INTO urls(short_code, long_url) VALUES(?, ?)",
		shortCode,
		longURL,
	)
	return err
}

func (s *SQLiteStore) Get(shortCode string) (string, error) {
	var longURL string
	err := s.db.QueryRow(
		"SELECT long_url FROM urls WHERE short_code = ?",
		shortCode,
	).Scan(&longURL)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}

	return longURL, err
}

func (s *SQLiteStore) GetMaxID() (uint64, error) {
	var maxShort string

	err := s.db.QueryRow(`
		SELECT short_code FROM urls 
		ORDER BY rowid DESC 
		LIMIT 1
	`).Scan(&maxShort)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	// Decode Base62 back to number
	var id uint64
	for i := 0; i < len(maxShort); i++ {
		id = id*62 + uint64(indexOf(maxShort[i]))
	}

	return id, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func indexOf(b byte) int {
	for i := 0; i < len(charset); i++ {
		if charset[i] == b {
			return i
		}
	}
	return 0
}
