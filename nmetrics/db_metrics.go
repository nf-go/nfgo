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
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
)

func (s *server) regitserDBCollector(config *nconf.Config) error {
	if config.DB != nil {
		s.dbMetricsCollector = newDBMetrics(config.DB)
		if err := s.registry.Register(s.dbMetricsCollector); err != nil {
			return err
		}
	}
	return nil
}

type dbMetrics struct {
	maxOpenConnections            *prometheus.GaugeVec
	openConnections               *prometheus.GaugeVec
	inUseConnections              *prometheus.GaugeVec
	idleConnections               *prometheus.GaugeVec
	waitCountForConnections       *prometheus.GaugeVec
	waitMillSecondsForConnections *prometheus.GaugeVec
	maxIdleClosedConnections      *prometheus.GaugeVec
	maxLifetimeClosedConnections  *prometheus.GaugeVec
}

func newDBMetrics(dbConfig *nconf.DbConfig) *dbMetrics {
	labelNames := []string{}
	return &dbMetrics{
		maxOpenConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_max_open_connections",
				Help: "Maximum number of open connections to the database.",
			},
			labelNames,
		),
		openConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_open_connections",
				Help: "The number of established connections both in use and idle.",
			},
			labelNames,
		),
		inUseConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_in_use_connections",
				Help: "The number of connections currently in use.",
			},
			labelNames,
		),
		idleConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_idle_connections",
				Help: "The number of idle connections.",
			},
			labelNames,
		),
		waitCountForConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_wait_count_for_connections",
				Help: "The total number of connections waited for.",
			},
			labelNames,
		),
		waitMillSecondsForConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_wait_milliseconds_for_connections",
				Help: "The total time blocked waiting for a new connection.",
			},
			labelNames,
		),
		maxIdleClosedConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_max_idle_closed_connections",
				Help: "The total number of connections closed due to SetMaxIdleConns.",
			},
			labelNames,
		),
		maxLifetimeClosedConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_max_life_time_closed_connections",
				Help: "The total number of connections closed due to SetConnMaxLifetime.",
			},
			labelNames,
		),
	}
}

func (m *dbMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.maxOpenConnections.Describe(ch)
	m.openConnections.Describe(ch)
	m.inUseConnections.Describe(ch)
	m.idleConnections.Describe(ch)
	m.waitCountForConnections.Describe(ch)
	m.waitMillSecondsForConnections.Describe(ch)
	m.maxIdleClosedConnections.Describe(ch)
	m.maxLifetimeClosedConnections.Describe(ch)
}

func (m *dbMetrics) Collect(ch chan<- prometheus.Metric) {
	m.maxOpenConnections.Collect(ch)
	m.openConnections.Collect(ch)
	m.inUseConnections.Collect(ch)
	m.idleConnections.Collect(ch)
	m.waitCountForConnections.Collect(ch)
	m.waitMillSecondsForConnections.Collect(ch)
	m.maxIdleClosedConnections.Collect(ch)
	m.maxLifetimeClosedConnections.Collect(ch)
}

func (m *dbMetrics) set(dbStats *sql.DBStats) {
	m.maxOpenConnections.WithLabelValues().Set(float64(dbStats.MaxOpenConnections))
	m.openConnections.WithLabelValues().Set(float64(dbStats.OpenConnections))
	m.inUseConnections.WithLabelValues().Set(float64(dbStats.InUse))
	m.idleConnections.WithLabelValues().Set(float64(dbStats.Idle))
	m.waitCountForConnections.WithLabelValues().Set(float64(dbStats.WaitCount))
	m.waitMillSecondsForConnections.WithLabelValues().Set(float64(dbStats.WaitDuration / time.Millisecond))
	m.maxIdleClosedConnections.WithLabelValues().Set(float64(dbStats.MaxIdleClosed))
	m.maxLifetimeClosedConnections.WithLabelValues().Set(float64(dbStats.MaxLifetimeClosed))
}

// gormPrometheusPlugin -
type gormPrometheusPlugin struct {
	dbMetricsCollector *dbMetrics
	initializeOnce     sync.Once
}

// Name -
func (p *gormPrometheusPlugin) Name() string {
	return "gorm:nfgo:prometheus"
}

// Initialize -
func (p *gormPrometheusPlugin) Initialize(db *gorm.DB) error { //can be called repeatedly
	p.initializeOnce.Do(
		func() {
			go func() {
				for range time.Tick(15 * time.Second) {
					sqlDB, err := db.DB()
					if err != nil {
						nlog.Warn("fail to get sql db in gorm prometheus plugins: ", err)
					}
					dbStatus := sqlDB.Stats()
					p.dbMetricsCollector.set(&dbStatus)
				}
			}()
		})
	return nil
}

// gormPrometheusPlugin -
func (s *server) gormPrometheusPlugin() gorm.Plugin {
	return &gormPrometheusPlugin{
		dbMetricsCollector: s.dbMetricsCollector,
		initializeOnce:     sync.Once{},
	}
}
