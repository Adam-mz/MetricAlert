package storage

import (
	"fmt"
	"strings"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateGauge(name string, value float64) {
	s.Gauge[name] = value
}

func (s *MemStorage) UpdateCounter(name string, value int64) {
	s.Counter[name] += value
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {

	value, exists := s.Gauge[name]
	return value, exists
}

func (s *MemStorage) GetCounter(name string) (int64, bool) {

	value, exists := s.Counter[name]
	return value, exists
}

func (s *MemStorage) GetAll() string {

	var builder strings.Builder

	builder.WriteString("<h1>Metrics</h1>")

	builder.WriteString("<h2>Gauge</h2><ul>")
	for key, value := range s.Gauge {
		builder.WriteString(fmt.Sprintf("<li>%s = %v</li>", key, value))
	}
	builder.WriteString("</ul>")

	builder.WriteString("<h2>Counter</h2><ul>")
	for key, value := range s.Counter {
		builder.WriteString(fmt.Sprintf("<li>%s = %v</li>", key, value))
	}
	builder.WriteString("</ul>")

	return builder.String()
}
