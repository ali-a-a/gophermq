package handler

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

const (
	namespace   = "gophermq"
	labelMethod = "method"
	labelCode   = "code"
)

type usageMetrics struct {
	RPS      *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

//nolint:gochecknoglobals
var metrics = usageMetrics{
	RPS: promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "client_request_rate",
		Help:      "Client request rate.",
	}, []string{labelMethod, labelCode}),
	Duration: promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "request_duration_seconds",
		Help:      "Request latencies in seconds.",
	}, []string{labelMethod, labelCode},
	),
}

func (u *usageMetrics) record(method string, code int, startTime time.Time) {
	u.RPS.With(prometheus.Labels{
		labelMethod: method,
		labelCode:   strconv.Itoa(code),
	}).Inc()

	u.Duration.With(prometheus.Labels{
		labelMethod: method,
		labelCode:   strconv.Itoa(code),
	}).Observe(time.Since(startTime).Seconds())
}
