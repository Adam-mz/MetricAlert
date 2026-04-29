package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/Adam-mz/MetricAlert/internal/handler"
	appmiddleware "github.com/Adam-mz/MetricAlert/internal/middleware"
	"github.com/Adam-mz/MetricAlert/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
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

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	mux := newMux(storage, logger)
	if err := http.ListenAndServe(cfg.ServerAddr, mux); err != nil {
		panic(err)
	}
}

func newMux(storage *storage.MemStorage, logger *zap.Logger) *chi.Mux {
	h := handler.Handler{Storage: storage}

	r := chi.NewRouter()

	// Middleware
	r.Use(appmiddleware.RequestLogger(logger))
	r.Use(middleware.Recoverer)

	r.Post("/update/{type}/{name}/{value}", h.Update)
	r.Get("/value/{type}/{name}", h.GetValue)
	r.Get("/", h.GetAllMetrics)
	r.Post("/update/", h.UpdateJSON)
	r.Post("/value/", h.GetValueJSON)

	return r
}
