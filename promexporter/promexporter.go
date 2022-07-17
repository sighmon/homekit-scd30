package promexporter

import (
	"flag"
	"net/http"
	"text/template"

	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	temperatureGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "temperature",
			Help: "Temperature measured (°C)",
		},
	)

	humidityGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "humidity",
			Help: "Relative humidity measured (%)",
		},
	)

	co2Gauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "co2",
			Help: "CO2 measured (ppm)",
		},
	)

	index = template.Must(template.New("index").Parse(
		`<!doctype html>
	 <title>SCD-30 Prometheus Exporter</title>
	 <h1>SCD-30 Prometheus Exporter</h1>
	 <a href="/metrics">Metrics</a>
	 <p>
	 `))
)

type Exporter struct {
	address string
}

func New(address string) *Exporter {
	return &Exporter{address: address}
}

func (e *Exporter) Start() {
	flag.Parse()
	log.Printf("Prometheus Exporter starting on port %s\n", e.address)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		index.Execute(w, "")
	})
	if err := http.ListenAndServe(e.address, nil); err != http.ErrServerClosed {
		panic(err)
	}
}

func (e *Exporter) UpdateReadings(temperature float64, humidity float64, co2 float64) {
	temperatureGauge.Set(temperature)
	humidityGauge.Set(humidity)
	co2Gauge.Set(co2)
}
