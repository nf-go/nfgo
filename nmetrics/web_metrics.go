// Copyright 2020 The nfgo Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nmetrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nf-go/nfgo/nconf"
	"github.com/prometheus/client_golang/prometheus"
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

func (s *server) WebMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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
