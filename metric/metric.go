package metric

import (
	"github.com/ali-a-a/gophermq/config"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func StartPrometheusServer(monitoring config.Monitoring) {
	if monitoring.Enable {
		metricServer := http.NewServeMux()
		metricServer.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(monitoring.Port, metricServer); err != nil {
			logrus.Panicf("failed to start prometheus metrics server %s", err.Error())
		}
	}
}
