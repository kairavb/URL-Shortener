package main

import (
	"database/sql"
	"log"
	"net/http"

	"url-shortener/internal/api"
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
	service := shortener.NewService(gen, store)

	handler := api.NewHandler(service)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, handler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
