package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
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
			sendCounter(baseURL, "PollCount", pollCount)
		}
	}
}

func sendGauge(baseURL, metric string, value float64) {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", baseURL, metric, value)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Printf("Error sending gauge %s: %v\n", metric, err)
		return
	}
	resp.Body.Close()
}

func sendCounter(baseURL, metric string, value int64) {
	url := fmt.Sprintf("%s/update/counter/%s/%d", baseURL, metric, value)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Printf("Error sending counter %s: %v\n", metric, err)
		return
	}
	resp.Body.Close()
}
