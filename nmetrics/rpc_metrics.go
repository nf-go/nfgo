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
