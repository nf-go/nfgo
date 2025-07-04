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

package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/nf-go/nfgo/nconf"
	"github.com/nf-go/nfgo/nlog"
	"github.com/nf-go/nfgo/nutil/graceful"
	"github.com/nf-go/nfgo/nutil/ntypes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server -
type Server interface {
	graceful.ShutdownServer

	GRPCServer() *grpc.Server
}

// NewServer -
func NewServer(config *nconf.Config, opt ...ServerOption) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nill")
	}
	rpcConfig := config.RPC
	if rpcConfig == nil {
		return nil, errors.New("rpc config is not initialized in the config")
	}
	opts := &serverOptions{}
	for _, o := range opt {
		o(opts)
	}
	opts.setInterceptors()

	grpcOpts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(rpcConfig.MaxRecvMsgSize)),
		grpc.ChainUnaryInterceptor(opts.unaryServerInterceptors...),
		grpc.ChainStreamInterceptor(opts.streamServerInterceptors...),
	}

	grpcServer := grpc.NewServer(grpcOpts...)

	if opts.metricsServer != nil {
		grpc_prometheus.Register(grpcServer)
	}

	if ntypes.BoolValue(rpcConfig.RegisterHealthServer) {
		healthServer := health.NewServer()
		healthServer.SetServingStatus("grpc.health.v1.Health", 1)
		healthpb.RegisterHealthServer(grpcServer, healthServer)
	}

	if ntypes.BoolValue(rpcConfig.RegisterReflectionServer) {
		reflection.Register(grpcServer)
	}

	return &server{
		Server: grpcServer,
		host:   rpcConfig.Host,
		port:   rpcConfig.Port}, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config, opt ...ServerOption) Server {
	grpcServer, err := NewServer(config, opt...)
	if err != nil {
		nlog.Fatal("fail to init grpc server: ", err)
	}
	return grpcServer
}

type server struct {
	*grpc.Server
	host string
	port int32
}

func (s *server) GRPCServer() *grpc.Server {
	return s.Server
}

func (s *server) Serve() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("fail to listen at %s, %w: ", addr, err)
	}

	nlog.Info("the grpc server is started and serving on ", addr)
	if err = s.Server.Serve(listen); err != nil {
		nlog.Error("the grpc server is stoped  with error ", err)
		return err
	}

	return nil
}

func (s *server) MustServe() {
	if err := s.Serve(); err != nil {
		nlog.Fatal("fail to start grpc server: ", err)
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	s.GracefulStop()
	return nil
}
