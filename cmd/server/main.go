package main

import (
	"net/http"

	"github.com/Adam-mz/MetricAlert/internal/handler"
	"github.com/Adam-mz/MetricAlert/internal/storage"
)

func main() {
	storage := storage.NewMemStorage()

	h := handler.Handler{Storage: storage}

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, h.Update)
	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err)
	}
}
