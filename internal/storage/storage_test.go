package storage

import (
	"strings"
	"testing"
)

func TestNewMemStorage_Empty(t *testing.T) {
	store := NewMemStorage()
	if store == nil {
		t.Fatalf("expected non-nil storage")
	}
	if len(store.Gauge) != 0 {
		t.Fatalf("expected empty gauge map, got %d", len(store.Gauge))
	}
	if len(store.Counter) != 0 {
		t.Fatalf("expected empty counter map, got %d", len(store.Counter))
	}
}

func TestMemStorage_UpdateGaugeAndGet(t *testing.T) {
	store := NewMemStorage()
	store.UpdateGauge("cpu", 3.14)

	value, ok := store.GetGauge("cpu")
	if !ok {
		t.Fatalf("expected gauge metric to exist")
	}
	if value != 3.14 {
		t.Fatalf("expected gauge value 3.14, got %v", value)
	}

	if _, ok := store.GetGauge("missing"); ok {
		t.Fatalf("did not expect missing gauge metric to exist")
	}
}

func TestMemStorage_UpdateCounterAccumulates(t *testing.T) {
	store := NewMemStorage()
	store.UpdateCounter("requests", 2)
	store.UpdateCounter("requests", 3)

	value, ok := store.GetCounter("requests")
	if !ok {
		t.Fatalf("expected counter metric to exist")
	}
	if value != 5 {
		t.Fatalf("expected counter value 5, got %v", value)
	}
}

func TestMemStorage_GetAllContainsMetrics(t *testing.T) {
	store := NewMemStorage()
	store.UpdateGauge("cpu", 1.5)
	store.UpdateCounter("requests", 7)

	all := store.GetAll()
	if all == "" {
		t.Fatalf("expected non-empty html")
	}
	if !strings.Contains(all, "<h1>Metrics</h1>") {
		t.Fatalf("expected header to be present")
	}
	if !strings.Contains(all, "<h2>Gauge</h2><ul>") {
		t.Fatalf("expected gauge section to be present")
	}
	if !strings.Contains(all, "<h2>Counter</h2><ul>") {
		t.Fatalf("expected counter section to be present")
	}
	if !strings.Contains(all, "<li>cpu = 1.5</li>") {
		t.Fatalf("expected gauge metric to be listed")
	}
	if !strings.Contains(all, "<li>requests = 7</li>") {
		t.Fatalf("expected counter metric to be listed")
	}
}
