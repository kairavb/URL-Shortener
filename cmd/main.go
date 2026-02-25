package main

import (
	"database/sql"
	"fmt"
	"log"

	"url-shortener/internal/shortener"
	"url-shortener/internal/storage"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:urls.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := storage.NewSQLiteStore(db)

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	maxID, err := store.GetMaxID()
	if err != nil {
		log.Fatal(err)
	}

	gen := shortener.NewGenerator(maxID)

	shortCode := gen.Generate()

	err = store.Save(shortCode, "https://example.com")
	if err != nil {
		log.Fatal(err)
	}

	longURL, err := store.Get(shortCode)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Short:", shortCode)
	fmt.Println("Long :", longURL)
}
