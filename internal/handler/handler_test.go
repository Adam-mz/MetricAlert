package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Adam-mz/MetricAlert/internal/storage"
	"github.com/go-chi/chi/v5"
)

func TestUpdate_GaugeOK(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodPost, "/update/gauge/temp/3.14", map[string]string{
		"type":  "gauge",
		"name":  "temp",
		"value": "3.14",
	})
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "OK" {
		t.Fatalf("expected body OK, got %q", rr.Body.String())
	}
	if value, ok := store.GetGauge("temp"); !ok || value != 3.14 {
		t.Fatalf("expected gauge to be updated, got %v (exists=%v)", value, ok)
	}
}

func TestUpdate_CounterIntOK(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodPost, "/update/counter/requests/2", map[string]string{
		"type":  "counter",
		"name":  "requests",
		"value": "2",
	})
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if value, ok := store.GetCounter("requests"); !ok || value != 2 {
		t.Fatalf("expected counter to be updated, got %v (exists=%v)", value, ok)
	}
}

func TestUpdate_CounterFloatOK(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodPost, "/update/counter/requests/2.9", map[string]string{
		"type":  "counter",
		"name":  "requests",
		"value": "2.9",
	})
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if value, ok := store.GetCounter("requests"); !ok || value != 2 {
		t.Fatalf("expected counter to be updated to 2, got %v (exists=%v)", value, ok)
	}
}

func TestUpdate_InvalidMethod(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodGet, "/update/gauge/temp/1", map[string]string{
		"type":  "gauge",
		"name":  "temp",
		"value": "1",
	})
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}

func TestUpdate_MissingParams(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodPost, "/update/gauge/temp/1", nil)
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestUpdate_InvalidGaugeValue(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodPost, "/update/gauge/temp/abc", map[string]string{
		"type":  "gauge",
		"name":  "temp",
		"value": "abc",
	})
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestUpdate_UnknownType(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodPost, "/update/unknown/temp/1", map[string]string{
		"type":  "unknown",
		"name":  "temp",
		"value": "1",
	})
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestGetValue_GaugeOK(t *testing.T) {
	store := storage.NewMemStorage()
	store.UpdateGauge("temp", 1.5)
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodGet, "/value/gauge/temp", map[string]string{
		"type": "gauge",
		"name": "temp",
	})
	rr := httptest.NewRecorder()

	h.GetValue(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "1.5" {
		t.Fatalf("expected body 1.5, got %q", rr.Body.String())
	}
}

func TestGetValue_CounterNotFound(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodGet, "/value/counter/requests", map[string]string{
		"type": "counter",
		"name": "requests",
	})
	rr := httptest.NewRecorder()

	h.GetValue(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestGetValue_InvalidMethod(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodPost, "/value/gauge/temp", map[string]string{
		"type": "gauge",
		"name": "temp",
	})
	rr := httptest.NewRecorder()

	h.GetValue(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}

func TestGetValue_UnknownType(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodGet, "/value/unknown/temp", map[string]string{
		"type": "unknown",
		"name": "temp",
	})
	rr := httptest.NewRecorder()

	h.GetValue(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestGetValue_MissingParams(t *testing.T) {
	store := storage.NewMemStorage()
	h := Handler{Storage: store}

	req := newRequestWithParams(http.MethodGet, "/value/gauge/temp", nil)
	rr := httptest.NewRecorder()

	h.GetValue(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestGetAllMetrics_ReturnsHTML(t *testing.T) {
	store := storage.NewMemStorage()
	store.UpdateGauge("cpu", 1.2)
	store.UpdateCounter("requests", 3)
	h := Handler{Storage: store}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	h.GetAllMetrics(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/html" {
		t.Fatalf("expected Content-Type text/html, got %q", ct)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<h1>Metrics</h1>") {
		t.Fatalf("expected html header in response")
	}
	if !strings.Contains(body, "cpu = 1.2") {
		t.Fatalf("expected gauge metric in response")
	}
	if !strings.Contains(body, "requests = 3") {
		t.Fatalf("expected counter metric in response")
	}
}

func newRequestWithParams(method, path string, params map[string]string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	if params == nil {
		return req
	}

	routeCtx := chi.NewRouteContext()
	for key, value := range params {
		routeCtx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	return req.WithContext(ctx)
}
