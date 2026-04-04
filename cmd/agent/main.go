package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"time"
)

var pollCount = 0
var baseURL = "http://localhost:8080"

func main() {
	var mem runtime.MemStats
	tickerTwo := time.NewTicker(2 * time.Second)
	tickerTen := time.NewTicker(10 * time.Second)
	defer tickerTwo.Stop()
	defer tickerTen.Stop()
	for {
		select {
		case <-tickerTen.C:

			sendGauge("Alloc", float64(mem.Alloc))
			sendGauge("BuckHashSys", float64(mem.BuckHashSys))
			sendGauge("Frees", float64(mem.Frees))
			sendGauge("GCCPUFraction", mem.GCCPUFraction)
			sendGauge("GCSys", float64(mem.GCSys))
			sendGauge("HeapAlloc", float64(mem.HeapAlloc))
			sendGauge("HeapIdle", float64(mem.HeapIdle))
			sendGauge("HeapInuse", float64(mem.HeapInuse))
			sendGauge("HeapObjects", float64(mem.HeapObjects))
			sendGauge("HeapReleased", float64(mem.HeapReleased))
			sendGauge("HeapSys", float64(mem.HeapSys))
			sendGauge("LastGC", float64(mem.LastGC))
			sendGauge("Lookups", float64(mem.Lookups))
			sendGauge("MCacheInuse", float64(mem.MCacheInuse))
			sendGauge("MCacheSys", float64(mem.MCacheSys))
			sendGauge("MSpanInuse", float64(mem.MSpanInuse))
			sendGauge("MSpanSys", float64(mem.MSpanSys))
			sendGauge("Mallocs", float64(mem.Mallocs))
			sendGauge("NextGC", float64(mem.NextGC))
			sendGauge("NumForcedGC", float64(mem.NumForcedGC))
			sendGauge("NumGC", float64(mem.NumGC))
			sendGauge("OtherSys", float64(mem.OtherSys))
			sendGauge("PauseTotalNs", float64(mem.PauseTotalNs))
			sendGauge("StackInuse", float64(mem.StackInuse))
			sendGauge("StackSys", float64(mem.StackSys))
			sendGauge("Sys", float64(mem.Sys))
			sendGauge("TotalAlloc", float64(mem.TotalAlloc))

			sendGauge("RandomValue", rand.Float64())
			sendCount("PollCount", int64(pollCount))

		case <-tickerTwo.C:
			runtime.ReadMemStats(&mem)
			pollCount++

		}

	}

}

func sendGauge(metric string, runtime float64) {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", baseURL, metric, runtime)

	fmt.Println(url)

	_, err := http.Post(url, "text/plain", nil)
	if err != nil {
		panic(err)
	}
}

func sendCount(metric string, runtime int64) {
	url := fmt.Sprintf("%s/update/counter/%s/%d", baseURL, metric, runtime)
	fmt.Println(url)
	fmt.Println()
	_, err := http.Post(url, "text/plain", nil)
	if err != nil {
		panic(err)
	}
}
