package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func mustRegisterCollector(c prometheus.Collector) {
	prometheus.MustRegister(c)
}

// MetricsEndpoint is a handler that generates a web page with all details
func createMetricsEndpoint() http.Handler {
	return promhttp.Handler()
}
