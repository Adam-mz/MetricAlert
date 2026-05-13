package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/Adam-mz/MetricAlert/internal/model"
	"github.com/Adam-mz/MetricAlert/internal/storage"
	"go.uber.org/zap"
)

func TestServerAcceptsGzipJSONUpdate(t *testing.T) {
	store := storage.NewMemStorage()
	mux := newMux(store, zap.NewNop())

	value := 2.5
	body := gzipJSON(t, models.Metrics{
		ID:    "temperature",
		MType: models.Gauge,
		Value: &value,
	})

	req := httptest.NewRequest(http.MethodPost, "/update/", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	got, ok := store.GetGauge("temperature")
	if !ok || got != value {
		t.Fatalf("expected stored gauge %v, got %v (exists=%v)", value, got, ok)
	}
}

func TestServerCompressesJSONResponseForGzipClient(t *testing.T) {
	store := storage.NewMemStorage()
	mux := newMux(store, zap.NewNop())

	value := 7.25
	store.UpdateGauge("temperature", value)

	body := bytes.NewBufferString(`{"id":"temperature","type":"gauge"}`)
	req := httptest.NewRequest(http.MethodPost, "/value/", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected Content-Encoding gzip, got %q", got)
	}

	responseBody := gunzipBody(t, rr.Body)
	var metric models.Metrics
	if err := json.Unmarshal(responseBody, &metric); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if metric.ID != "temperature" || metric.MType != models.Gauge {
		t.Fatalf("unexpected metric: %+v", metric)
	}
	if metric.Value == nil || *metric.Value != value {
		t.Fatalf("expected value %v, got %+v", value, metric.Value)
	}
}

func TestServerCompressesHTMLResponseForGzipClient(t *testing.T) {
	store := storage.NewMemStorage()
	mux := newMux(store, zap.NewNop())

	store.UpdateGauge("cpu", 1.2)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected Content-Encoding gzip, got %q", got)
	}

	responseBody := string(gunzipBody(t, rr.Body))
	if !strings.Contains(responseBody, "<h1>Metrics</h1>") {
		t.Fatalf("expected html header in response")
	}
	if !strings.Contains(responseBody, "cpu = 1.2") {
		t.Fatalf("expected gauge metric in response")
	}
}

func gzipJSON(t *testing.T, value any) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		t.Fatalf("failed to encode gzip json: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}

	return &buf
}

func gunzipBody(t *testing.T, body *bytes.Buffer) []byte {
	t.Helper()

	reader, err := gzip.NewReader(body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read gzip body: %v", err)
	}

	return data
}
