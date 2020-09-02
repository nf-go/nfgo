package rpc

import (
	"errors"
	"fmt"
	"net"

	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/ntypes"
	"nfgo.ga/nfgo/rpc/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Server -
type Server interface {
	Run() error
	MustRun()
	GRPCServer() *grpc.Server
}

// ServerOption -
type ServerOption struct {
	UnaryServerInterceptors  []grpc.UnaryServerInterceptor
	StreamServerInterceptors []grpc.StreamServerInterceptor
}

func (o *ServerOption) setDefaultValue() {
	if len(o.UnaryServerInterceptors) == 0 {
		o.UnaryServerInterceptors = []grpc.UnaryServerInterceptor{
			interceptor.RecoverUnaryServerInterceptor,
			interceptor.MDCBindingUnaryServerInterceptor,
			interceptor.ValidateUnaryServerInterceptor,
			interceptor.LoggingUnaryServerInterceptor,
			interceptor.ErrorHandleUnaryServerInterceptor,
		}
	} else {
		intercepts := []grpc.UnaryServerInterceptor{
			interceptor.RecoverUnaryServerInterceptor,
		}
		o.UnaryServerInterceptors = append(intercepts, o.UnaryServerInterceptors...)
	}

	if len(o.StreamServerInterceptors) == 0 {
		o.StreamServerInterceptors = []grpc.StreamServerInterceptor{}
	} else {
		intercepts := []grpc.StreamServerInterceptor{
			interceptor.RecoverStreamServerInterceptor,
		}
		o.StreamServerInterceptors = append(intercepts, o.StreamServerInterceptors...)
	}
}

// NewServer -
func NewServer(config *nconf.Config, option *ServerOption) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nill")
	}
	rpcConfig := config.RPC
	if rpcConfig == nil {
		return nil, errors.New("rpc config is not initialized in the config")
	}
	if option == nil {
		option = &ServerOption{}
	}
	option.setDefaultValue()

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(rpcConfig.MaxRecvMsgSize)),
		grpc.ChainUnaryInterceptor(option.UnaryServerInterceptors...),
		grpc.ChainStreamInterceptor(option.StreamServerInterceptors...),
	}

	grpcServer := grpc.NewServer(opts...)

	if ntypes.BoolValue(rpcConfig.RegisterHealthServer) {
		healthServer := health.NewServer()
		healthServer.SetServingStatus("grpc.health.v1.Health", 1)
		healthpb.RegisterHealthServer(grpcServer, healthServer)
	}

	return &server{
		Server: grpcServer,
		host:   rpcConfig.Host,
		port:   rpcConfig.Port}, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config, option *ServerOption) Server {
	grpcServer, err := NewServer(config, option)
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

func (s *server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("fail to listen at %s, %w: ", addr, err)
	}

	nlog.Info("the grpc server is started and serving on ", addr)
	if err = s.Serve(listen); err != nil {
		nlog.Error("the grpc server is stoped  with error ", err)
		return err
	}

	return nil
}

func (s *server) MustRun() {
	if err := s.Run(); err != nil {
		nlog.Fatal("fail to start grpc server: ", err)
	}
}
