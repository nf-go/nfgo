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
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/nf-go/nfgo/nconf"
	"github.com/nf-go/nfgo/nlog"
	"github.com/nf-go/nfgo/nutil/graceful"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// Server -
type Server interface {
	graceful.ShutdownServer

	RegisterCollectors(collectors ...prometheus.Collector) error

	MustRegisterCollectors(collectors ...prometheus.Collector)

	GrpcMetricsUnaryServerInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)

	GrpcMetricsStramServerInterceptor() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error

	WebMetricsMiddleware() gin.HandlerFunc
}

// NewServer -
func NewServer(config *nconf.Config, opt ...ServerOption) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	metricsConfig := config.Metrics
	if metricsConfig == nil {
		return nil, errors.New("metrics config is not initialized in the config")
	}

	opts := &serverOptions{}
	for _, o := range opt {
		o(opts)
	}

	s := &server{
		registry:      prometheus.NewRegistry(),
		metricsConfig: metricsConfig,
		opts:          opts,
	}

	serverMux := http.NewServeMux()
	serverMux.Handle(metricsConfig.MetricsPath, promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", metricsConfig.Host, metricsConfig.Port),
		Handler: serverMux,
	}

	if err := s.registerCollectors(config); err != nil {
		return nil, err
	}

	return s, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config, opt ...ServerOption) Server {
	server, err := NewServer(config, opt...)
	if err != nil {
		nlog.Fatal("fail to init promtheus metrics server: ", err)
	}
	return server
}

type server struct {
	metricsConfig        *nconf.MetricsConfig
	opts                 *serverOptions
	httpServer           *http.Server
	registry             *prometheus.Registry
	grpcMetricsCollector *grpc_prometheus.ServerMetrics
	webMetricsCollector  *webMetrics
}

func (s *server) registerCollectors(config *nconf.Config) error {
	if err := s.regitserBuildinCollector(config); err != nil {
		return err
	}
	if err := s.registerRPCCollector(config); err != nil {
		return err
	}
	if err := s.regitserDBCollector(config); err != nil {
		return err
	}
	return s.regitserWebCollector(config)
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *server) Serve() error {
	nlog.Infof("the prometheus metrics server is started and serving on http://%s%s", s.httpServer.Addr, s.metricsConfig.MetricsPath)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		nlog.Error("the prometheus metrics server is stoped with error ", err)
		return err
	}
	return nil
}

func (s *server) MustServe() {
	if err := s.Serve(); err != nil {
		nlog.Fatal("fail to start promtheus metrics server: ", err)
	}
}

func (s *server) RegisterCollectors(collectors ...prometheus.Collector) error {
	for _, collector := range collectors {
		if err := s.registry.Register(collector); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) MustRegisterCollectors(collectors ...prometheus.Collector) {
	if err := s.RegisterCollectors(collectors...); err != nil {
		nlog.Fatal("fail to register collectors: ", err)
	}
}
