package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
		},
		[]string{"pattern", "method", "status"},
	)

	HttpRequestsCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_inflight_current",
		},
		[]string{},
	)

	HttpRequestsDurationHistorgram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds_historgram",
			Buckets: []float64{
				0.1,  // 100 ms
				0.2,  // 200 ms
				0.25, // 250 ms
				0.5,  // 500 ms
				1,    // 1 s
			},
		},
		[]string{"pattern", "method"},
	)

	TaskCountTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_count_total",
		},
		[]string{},
	)

	TaskCountCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "task_count_current",
		},
		[]string{},
	)

	TaskProcessingDurationHistorgram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "task_processing_duration_seconds_historgram",
			Buckets: []float64{
				30,     // 30 s
				60 * 1, // 1 min
				60 * 3, //3 min
				60 * 5, // 5 min
			},
		},
		[]string{},
	)
)
