package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	cacheMisses      *prometheus.CounterVec
	cacheHits        *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec
}

var Reg = prometheus.NewRegistry()

var M = newMetrics()

func newMetrics() *metrics {
	m := &metrics{
		cacheMisses: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses",
				Help: "Number of cache misses.",
			},
			[]string{"method"},
		),

		cacheHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits",
				Help: "Number of cache hits.",
			},
			[]string{"method"},
		),

		latencyHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "request_duration_seconds",
				Help:    "Latency of requests.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method"},
		),
	}

	Reg.MustRegister(m.cacheMisses)
	Reg.MustRegister(m.latencyHistogram)
	return m
}

func (m *metrics) ObserveLatency(method string, start time.Time) {
	m.latencyHistogram.WithLabelValues(method).Observe(time.Since(start).Seconds())
}

func (m *metrics) IncCacheMiss(method string) {
	m.cacheMisses.WithLabelValues(method).Inc()
}

func (m *metrics) IncCacheHit(method string) {
	m.cacheHits.WithLabelValues(method).Inc()
}
