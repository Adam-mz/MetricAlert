package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/Adam-mz/MetricAlert/internal/handler"
	"github.com/Adam-mz/MetricAlert/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	ServerAddr string
}

func main() {
	cfg := Config{}
	
	storage := storage.NewMemStorage()
	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "HTTP server address")
	flag.Parse()

	if addr := os.Getenv("ADDRESS"); addr != "" {
		cfg.ServerAddr = addr
	}

	mux := newMux(storage)
	if err := http.ListenAndServe(cfg.ServerAddr, mux); err != nil {
		panic(err)
	}
}

func newMux(storage *storage.MemStorage) *chi.Mux {
	h := handler.Handler{Storage: storage}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Post("/update/{type}/{name}/{value}", h.Update)
	r.Get("/value/{type}/{name}", h.GetValue)
	r.Get("/", h.GetAllMetrics)

	return r
}
