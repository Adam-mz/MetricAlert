package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	models "github.com/Adam-mz/MetricAlert/internal/model"
)

type Config struct {
	ServerAddr     string
	PollInterval   int
	ReportInterval int
}

func main() {
	// 1. DEFAULT значения
	cfg := Config{
		ServerAddr:     "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}

	// 2. Флаги командной строки
	flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "Server address")
	flag.IntVar(&cfg.PollInterval, "p", cfg.PollInterval, "Poll interval in seconds")
	flag.IntVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "Report interval in seconds")
	flag.Parse()

	// 3. Переменные окружения (имеют приоритет над флагами)
	if addr := os.Getenv("ADDRESS"); addr != "" {
		cfg.ServerAddr = addr
	}
	if poll := os.Getenv("POLL_INTERVAL"); poll != "" {
		if v, err := strconv.Atoi(poll); err == nil {
			cfg.PollInterval = v
		}
	}
	if report := os.Getenv("REPORT_INTERVAL"); report != "" {
		if v, err := strconv.Atoi(report); err == nil {
			cfg.ReportInterval = v
		}
	}

	// Формируем полный URL
	baseURL := cfg.ServerAddr
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	fmt.Printf("Starting agent with config:\n")
	fmt.Printf("  Server: %s\n", baseURL)
	fmt.Printf("  Poll interval: %d seconds\n", cfg.PollInterval)
	fmt.Printf("  Report interval: %d seconds\n", cfg.ReportInterval)

	var mem runtime.MemStats
	var pollCount int64

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			// Обновляем метрики каждые pollInterval секунд
			runtime.ReadMemStats(&mem)
			pollCount++

		case <-reportTicker.C:
			// Отправляем метрики каждые reportInterval секунд
			sendGauge(baseURL, "Alloc", float64(mem.Alloc))
			sendGauge(baseURL, "BuckHashSys", float64(mem.BuckHashSys))
			sendGauge(baseURL, "Frees", float64(mem.Frees))
			sendGauge(baseURL, "GCCPUFraction", mem.GCCPUFraction)
			sendGauge(baseURL, "GCSys", float64(mem.GCSys))
			sendGauge(baseURL, "HeapAlloc", float64(mem.HeapAlloc))
			sendGauge(baseURL, "HeapIdle", float64(mem.HeapIdle))
			sendGauge(baseURL, "HeapInuse", float64(mem.HeapInuse))
			sendGauge(baseURL, "HeapObjects", float64(mem.HeapObjects))
			sendGauge(baseURL, "HeapReleased", float64(mem.HeapReleased))
			sendGauge(baseURL, "HeapSys", float64(mem.HeapSys))
			sendGauge(baseURL, "LastGC", float64(mem.LastGC))
			sendGauge(baseURL, "Lookups", float64(mem.Lookups))
			sendGauge(baseURL, "MCacheInuse", float64(mem.MCacheInuse))
			sendGauge(baseURL, "MCacheSys", float64(mem.MCacheSys))
			sendGauge(baseURL, "MSpanInuse", float64(mem.MSpanInuse))
			sendGauge(baseURL, "MSpanSys", float64(mem.MSpanSys))
			sendGauge(baseURL, "Mallocs", float64(mem.Mallocs))
			sendGauge(baseURL, "NextGC", float64(mem.NextGC))
			sendGauge(baseURL, "NumForcedGC", float64(mem.NumForcedGC))
			sendGauge(baseURL, "NumGC", float64(mem.NumGC))
			sendGauge(baseURL, "OtherSys", float64(mem.OtherSys))
			sendGauge(baseURL, "PauseTotalNs", float64(mem.PauseTotalNs))
			sendGauge(baseURL, "StackInuse", float64(mem.StackInuse))
			sendGauge(baseURL, "StackSys", float64(mem.StackSys))
			sendGauge(baseURL, "Sys", float64(mem.Sys))
			sendGauge(baseURL, "TotalAlloc", float64(mem.TotalAlloc))

			sendGauge(baseURL, "RandomValue", rand.Float64())
			if sendCounter(baseURL, "PollCount", pollCount) {
				pollCount = 0
			}
		}
	}
}

func sendGauge(baseURL, metric string, value float64) {
	sendMetric(baseURL, models.Metrics{
		ID:    metric,
		MType: models.Gauge,
		Value: &value,
	})
}

func sendCounter(baseURL, metric string, delta int64) bool {
	return sendMetric(baseURL, models.Metrics{
		ID:    metric,
		MType: models.Counter,
		Delta: &delta,
	})
}

func sendMetric(baseURL string, metric models.Metrics) bool {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(metric); err != nil {
		fmt.Printf("Error encoding metric %s: %v\n", metric.ID, err)
		return false
	}

	resp, err := http.Post(baseURL+"/update/", "application/json", &buf)
	if err != nil {
		fmt.Printf("Error sending metric %s: %v\n", metric.ID, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error sending metric %s: server returned %d\n", metric.ID, resp.StatusCode)
		return false
	}

	return true
}
