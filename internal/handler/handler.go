package handler

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strconv"

	models "github.com/Adam-mz/MetricAlert/internal/model"
	"github.com/Adam-mz/MetricAlert/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage *storage.MemStorage
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из URL с помощью chi
	mType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	// Проверяем, что все параметры присутствуют
	if mType == "" || name == "" || value == "" {
		http.Error(w, "invalid request: missing parameters", http.StatusBadRequest)
		return
	}

	switch mType {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "invalid value", http.StatusBadRequest)
			return
		}
		h.Storage.UpdateGauge(name, v)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			// Если не получилось, пробуем как float64 и преобразуем в int64
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(w, "invalid counter value", http.StatusBadRequest)
				return
			}
			v = int64(f) // отбрасывает дробную часть
			// или v = int64(math.Round(f)) для округления
		}
		h.Storage.UpdateCounter(name, v)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

}

func (h *Handler) GetValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	if mType == "" || name == "" {
		http.Error(w, "invalid request: missing parameters", http.StatusBadRequest)
		return
	}

	switch mType {
	case "gauge":
		value, exists := h.Storage.GetGauge(name)
		if !exists {
			http.Error(w, "gauge metric not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%v", value)

	case "counter":
		value, exists := h.Storage.GetCounter(name)
		if !exists {
			http.Error(w, "counter metric not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%v", value)
	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
	}

}

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, h.Storage.GetAll())
}

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	if !isJSONContentType(r) {
		http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	var metric models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if metric.ID == "" {
		http.Error(w, "missing metric id", http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case models.Gauge:
		if metric.Value == nil {
			http.Error(w, "missing gauge value", http.StatusBadRequest)
			return
		}

		h.Storage.UpdateGauge(metric.ID, *metric.Value)
		value, _ := h.Storage.GetGauge(metric.ID)
		metric.Value = &value
		metric.Delta = nil

	case models.Counter:
		if metric.Delta == nil {
			http.Error(w, "missing counter delta", http.StatusBadRequest)
			return
		}

		h.Storage.UpdateCounter(metric.ID, *metric.Delta)
		delta, _ := h.Storage.GetCounter(metric.ID)
		metric.Delta = &delta
		metric.Value = nil

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metric)

}
func (h *Handler) GetValueJSON(w http.ResponseWriter, r *http.Request) {
	if !isJSONContentType(r) {
		http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	var metric models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if metric.ID == "" {
		http.Error(w, "missing metric id", http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case models.Gauge:
		value, ok := h.Storage.GetGauge(metric.ID)
		if !ok {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		metric.Value = &value
		metric.Delta = nil

	case models.Counter:
		delta, ok := h.Storage.GetCounter(metric.ID)
		if !ok {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		metric.Delta = &delta
		metric.Value = nil

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metric)
}

func isJSONContentType(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	return mediaType == "application/json"
}
