package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Adam-mz/MetricAlert/internal/storage"
)

type Handler struct {
	Storage *storage.MemStorage
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "invalid request", http.StatusNotFound)
		return
	}

	mType := parts[2]
	name := parts[3]
	value := parts[4]

	switch mType {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "invalid value", http.StatusBadRequest)
			return
		}
		h.Storage.UpdateGauge(name, v)

	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "invalid value", http.StatusBadRequest)
			return
		}
		h.Storage.UpdateCounter(name, v)

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
