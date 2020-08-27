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
	Run(serverOption *ServerOption) error

	MustRun(serverOption *ServerOption)

	RegisterCollectors(collectors ...prometheus.Collector) error

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
func NewServer(config *nconf.Config) (Server, error) {
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
	}

	if err := s.registerCollectors(config); err != nil {
		return nil, err
	}

	return s, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config) Server {
	server, err := NewServer(config)
	if err != nil {
		nlog.Fatal("fail to init promtheus metrics server: ", err)
	}
	return server
}

type server struct {
	metricsConfig        *nconf.MetricsConfig
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

func (s *server) Run(serverOption *ServerOption) error {

	if serverOption != nil {
		rpcServer := serverOption.RPCServer
		if rpcServer != nil && rpcServer.GRPCServer() != nil {
			grpc_prometheus.Register(rpcServer.GRPCServer())
		}
		db := serverOption.DB
		if db != nil {
			if err := db.Use(s.gormPrometheusPlugin()); err != nil {
				return err
			}
		}
	}

	metricsConfig := s.metricsConfig
	addr := fmt.Sprintf("%s:%d", metricsConfig.Host, metricsConfig.Port)
	serverMux := http.NewServeMux()
	serverMux.Handle(metricsConfig.MetricsPath, promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))
	nlog.Infof("the prometheus metrics server is started and serving on http://%s%s", addr, metricsConfig.MetricsPath)
	if err := http.ListenAndServe(addr, serverMux); err != nil {
		nlog.Error("the prometheus metrics server is stoped with error ", err)
		return err
	}

	return nil
}

func (s *server) MustRun(serverOption *ServerOption) {
	if err := s.Run(serverOption); err != nil {
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
