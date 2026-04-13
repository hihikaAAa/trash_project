// Package metrics
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	LatencyAPI = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "paperless_api_latency_seconds",
			Help:    "API endpoint latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status_class"},
	)
)
