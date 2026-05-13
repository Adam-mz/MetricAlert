package main

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/Adam-mz/MetricAlert/internal/model"
)

func TestSendMetricSendsGzipRequest(t *testing.T) {
	value := 12.5

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/update/" {
			t.Fatalf("expected path /update/, got %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("expected Content-Type application/json, got %q", got)
		}
		if got := r.Header.Get("Content-Encoding"); got != "gzip" {
			t.Fatalf("expected Content-Encoding gzip, got %q", got)
		}
		if got := r.Header.Get("Accept-Encoding"); got != "gzip" {
			t.Fatalf("expected Accept-Encoding gzip, got %q", got)
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			t.Fatalf("failed to create gzip reader: %v", err)
		}
		defer reader.Close()

		var metric models.Metrics
		if err := json.NewDecoder(reader).Decode(&metric); err != nil {
			t.Fatalf("failed to decode metric: %v", err)
		}

		if metric.ID != "Alloc" || metric.MType != models.Gauge {
			t.Fatalf("unexpected metric: %+v", metric)
		}
		if metric.Value == nil || *metric.Value != value {
			t.Fatalf("expected value %v, got %+v", value, metric.Value)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ok := sendMetric(server.URL, models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &value,
	})

	if !ok {
		t.Fatal("expected sendMetric to return true")
	}
}
