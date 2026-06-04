package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/internal/api"
	"url-shortener/internal/cache"
	"url-shortener/internal/ratelimit"
	"url-shortener/internal/shortener"
	"url-shortener/internal/storage"

	_ "modernc.org/sqlite"
)

const (
	addr         = ":8080"
	cacheSize    = 1000
	ratePerSec   = 20
	rateBurst    = 40
	shutdownWait = 10 * time.Second
)

func main() {
	// Open local SQLite file — created automatically on first run.
	db, err := sql.Open("sqlite", "file:urls.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1) // SQLite works best with a single writer

	store := storage.NewSQLiteStore(db)
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Resume the ID counter so short codes stay unique after restarts.
	maxID, err := store.GetMaxID()
	if err != nil {
		log.Fatal(err)
	}

	gen := shortener.NewGenerator(maxID)
	urlCache := cache.New(cacheSize)
	service := shortener.NewService(gen, store, urlCache)
	handler := api.NewHandler(service)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, handler, store)

	// Outermost middleware runs first: rate limit, then request logging.
	limiter := ratelimit.New(ratePerSec, rateBurst)
	server := &http.Server{
		Addr:    addr,
		Handler: limiter.Middleware(api.Logging(mux)),
	}

	go func() {
		log.Printf("server running on http://localhost%s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for Ctrl+C or SIGTERM, then drain in-flight requests.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), shutdownWait)
	defer cancel()

	log.Println("shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
