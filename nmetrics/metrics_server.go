package nmetrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/rpc"
	"nfgo.ga/nfgo/web"
)

// Server -
type Server interface {
	Serve() error

	MustServe()

	Shutdown(ctx context.Context) error

	RegisterCollectors(collectors ...prometheus.Collector) error

	MustRegisterCollectors(collectors ...prometheus.Collector)

	GrpcMetricsUnaryServerInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)

	GrpcMetricsStramServerInterceptor() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error

	WebMetricsMiddleware() web.HandlerFunc
}

// ServerOption -
type ServerOption struct {
	RPCServer rpc.Server
	DB        *gorm.DB
}

// NewServer -
func NewServer(config *nconf.Config, serverOption *ServerOption) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	metricsConfig := config.Metrics
	if metricsConfig == nil {
		return nil, errors.New("metrics config is not initialized in the config")
	}

	s := &server{
		registry:      prometheus.NewRegistry(),
		metricsConfig: metricsConfig,
		serverOption:  serverOption,
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
func MustNewServer(config *nconf.Config, serverOption *ServerOption) Server {
	server, err := NewServer(config, serverOption)
	if err != nil {
		nlog.Fatal("fail to init promtheus metrics server: ", err)
	}
	return server
}

type server struct {
	metricsConfig        *nconf.MetricsConfig
	serverOption         *ServerOption
	httpServer           *http.Server
	registry             *prometheus.Registry
	grpcMetricsCollector *grpc_prometheus.ServerMetrics
	webMetricsCollector  *webMetrics
	dbMetricsCollector   *dbMetrics
}

func (s *server) registerCollectors(config *nconf.Config) error {
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

	if s.serverOption != nil {
		rpcServer := s.serverOption.RPCServer
		if rpcServer != nil && rpcServer.GRPCServer() != nil {
			grpc_prometheus.Register(rpcServer.GRPCServer())
		}
		db := s.serverOption.DB
		if db != nil {
			if err := db.Use(s.gormPrometheusPlugin()); err != nil {
				return err
			}
		}
	}

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
