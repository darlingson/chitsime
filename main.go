package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/darlingson/chitsime/internal/api"
	"github.com/darlingson/chitsime/internal/store"
)

//go:embed frontend/dist
var frontend embed.FS

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbPath := env("DB_PATH", "data/chitsime.db")
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		log.Error("database not found, run setup first", "path", dbPath, "command", "go run ./cmd/setup")
		os.Exit(1)
	}

	st, err := store.Open(dbPath)
	if err != nil {
		log.Error("open database", "err", err)
		os.Exit(1)
	}
	defer st.Close()

	dist, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Error("load frontend", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", api.New(st, log))
	mux.Handle("/", spa(dist))

	addr := env("ADDR", ":8080")
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Info("api listening", "addr", addr, "db", dbPath)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server", "err", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Info("stopped")
}

func spa(dist fs.FS) http.Handler {
	files := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" || name == "." {
			name = "index.html"
		}

		if info, err := fs.Stat(dist, name); err == nil && !info.IsDir() {
			files.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			http.ServeFileFS(w, r, dist, "index.html")
			return
		}

		http.NotFound(w, r)
	})
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
