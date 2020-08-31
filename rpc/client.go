package rpc

import (
	"errors"

	"context"

	"google.golang.org/grpc"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/ntypes"
	"nfgo.ga/nfgo/rpc/interceptor"
)

// ClientConns -
type ClientConns interface {
	//  GetClientConn - A ClientConn can safely be accessed concurrently.
	// https://github.com/grpc/grpc-go/blob/master/Documentation/concurrency.md
	GetClientConn(svcName string) *grpc.ClientConn
	// CloseAll -
	CloseAll()
}

type clientConns map[string]*grpc.ClientConn

func (cs clientConns) GetClientConn(svcName string) *grpc.ClientConn {
	return cs[svcName]
}

func (cs clientConns) CloseAll() {
	for _, conn := range cs {
		conn.Close()
	}
}

// NewClientConns -
func NewClientConns(config *nconf.Config) (ClientConns, error) {
	if config == nil {
		return nil, errors.New("config is nill")
	}
	if config.RPC == nil {
		return nil, errors.New("rpc config is not initialized in the config")
	}
	conns := clientConns(map[string]*grpc.ClientConn{})
	for svcName, clientConf := range config.RPC.Clients {
		conn, err := dialConn(clientConf)
		if err != nil {
			nlog.Logger(context.Background()).WithError(err).Errorf("fail to dial connection for %s", clientConf.Addr)
		}
		conns[svcName] = conn
	}
	return conns, nil
}

// MustNewClientConns -
func MustNewClientConns(config *nconf.Config) ClientConns {
	conns, err := NewClientConns(config)
	if err != nil {
		nlog.Fatal("fail to create rpc client connections: ", err)
	}
	return conns
}

func dialConn(config *nconf.RPCClientConfig) (*grpc.ClientConn, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(config.MaxCallSendMsgSize),
			grpc.MaxCallRecvMsgSize(config.MaxCallRecvMsgSize),
		),
		grpc.WithChainUnaryInterceptor(interceptor.ChainUnaryClientInterceptor(
			interceptor.MDCBindingUnaryClientInterceptor,
			interceptor.ValidateUnaryClientInterceptor,
		)),
		grpc.WithChainStreamInterceptor(interceptor.ChainStreamClientInterceptor(
			interceptor.MDCBindingStreamClientInterceptor,
		)),
	}
	if ntypes.BoolValue(config.Plaintext) {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	// Dial connection
	return grpc.Dial(config.Addr, dialOptions...)
}
