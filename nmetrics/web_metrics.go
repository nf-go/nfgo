package nmetrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/web"
)

func (s *server) regitserWebCollector(config *nconf.Config) error {
	if config.Web != nil {
		s.webMetricsCollector = newWebMetrics()
		if err := s.registry.Register(s.webMetricsCollector); err != nil {
			return err
		}
	}
	return nil
}

type webMetrics struct {
	reqCountTotal      *prometheus.CounterVec
	reqDurationSeconds *prometheus.HistogramVec
	reqSizeBytes       *prometheus.SummaryVec
	respSizeBytes      *prometheus.SummaryVec
}

func newWebMetrics() *webMetrics {
	labelNames := []string{"status_code", "path", "method"}
	return &webMetrics{
		reqCountTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_count_total",
				Help: "Total number of HTTP requests made.",
			}, labelNames),
		reqDurationSeconds: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latencies in seconds.",
				Buckets: prometheus.DefBuckets,
			}, labelNames),
		reqSizeBytes: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_request_size_bytes",
				Help: "HTTP request sizes in bytes.",
			}, labelNames),
		respSizeBytes: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_response_size_bytes",
				Help: "HTTP request sizes in bytes.",
			}, labelNames),
	}
}

func (m *webMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.reqCountTotal.Describe(ch)
	m.reqDurationSeconds.Describe(ch)
	m.reqSizeBytes.Describe(ch)
	m.respSizeBytes.Describe(ch)
}

func (m *webMetrics) Collect(ch chan<- prometheus.Metric) {
	m.reqCountTotal.Collect(ch)
	m.reqDurationSeconds.Collect(ch)
	m.reqSizeBytes.Collect(ch)
	m.respSizeBytes.Collect(ch)
}

func (s *server) WebMetricsMiddleware() web.HandlerFunc {
	return func(c *web.Context) {
		start := time.Now()
		c.Next()
		statusCode := strconv.Itoa(c.Writer.Status())
		path := c.Request.URL.Path
		method := c.Request.Method
		lvs := []string{statusCode, path, method}

		collector := s.webMetricsCollector
		collector.reqCountTotal.WithLabelValues(lvs...).Inc()
		collector.reqDurationSeconds.WithLabelValues(lvs...).Observe(time.Since(start).Seconds())
		collector.reqSizeBytes.WithLabelValues(lvs...).Observe(float64(c.Request.ContentLength))
		collector.respSizeBytes.WithLabelValues(lvs...).Observe(float64(c.Writer.Size()))
	}
}
