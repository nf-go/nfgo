package rpc

import (
	"errors"
	"fmt"
	"net"

	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/rpc/interceptor"

	"google.golang.org/grpc"
)

// Server -
type Server interface {
	Run() error
	GRPCServer() *grpc.Server
}

// NewServer -
func NewServer(config *nconf.Config, interceptors ...grpc.UnaryServerInterceptor) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nill")
	}

	rpcConfig := config.RPC
	if rpcConfig == nil {
		return nil, errors.New("rpc config is not initialized in the config")
	}

	host := rpcConfig.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := rpcConfig.Port
	if port == 0 {
		port = 9090
	}

	var chainInterceptor grpc.UnaryServerInterceptor

	if len(interceptors) == 0 {
		chainInterceptor = interceptor.ChainUnaryServerInterceptor(
			interceptor.MDCBindingUnaryServerInterceptor,
			interceptor.ValidateUnaryServerInterceptor,
			interceptor.LoggingUnaryServerInterceptor,
			interceptor.ErrorHandleUnaryServerInterceptor)
	} else {
		chainInterceptor = interceptor.ChainUnaryServerInterceptor(interceptors...)
	}

	maxRecvMsgSize := int(rpcConfig.MaxRecvMsgSize)
	if maxRecvMsgSize == 0 {
		maxRecvMsgSize = 50 << 20 // 50MiB
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(maxRecvMsgSize),
		grpc.UnaryInterceptor(chainInterceptor),
	}

	grpcServer := grpc.NewServer(opts...)
	return &server{
		Server: grpcServer,
		host:   host,
		port:   port}, nil
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
