package nmetrics

import (
	"context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"nfgo.ga/nfgo/nconf"
)

func (s *server) registerRPCCollector(config *nconf.Config) error {
	if config.RPC != nil {
		s.grpcMetricsCollector = grpc_prometheus.NewServerMetrics()
		s.grpcMetricsCollector.EnableHandlingTimeHistogram()
		if err := s.registry.Register(s.grpcMetricsCollector); err != nil {
			return err
		}
	}
	return nil
}

// GrpcMetricsUnaryServerInterceptor - may return nil
func (s *server) GrpcMetricsUnaryServerInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return s.grpcMetricsCollector.UnaryServerInterceptor()
}

// GrpcMetricsUnaryServerInterceptor - may return nil
func (s *server) GrpcMetricsStramServerInterceptor() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return s.grpcMetricsCollector.StreamServerInterceptor()
}
