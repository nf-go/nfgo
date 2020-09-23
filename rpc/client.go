package rpc

import (
	"crypto/tls"
	"errors"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
		if err := conn.Close(); err != nil {
			nlog.Error("fail to close grpc client conn: ", err)
		}
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
		grpc.WithChainUnaryInterceptor(
			interceptor.MDCBindingUnaryClientInterceptor,
			interceptor.ValidateUnaryClientInterceptor,
		),
		grpc.WithChainStreamInterceptor(
			interceptor.MDCBindingStreamClientInterceptor,
		),
	}
	if ntypes.BoolValue(config.Plaintext) {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	} else {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		})
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	}

	// Dial connection
	return grpc.Dial(config.Addr, dialOptions...)
}

// MustDialClientConnPlaintext -
func MustDialClientConnPlaintext(addr string) *grpc.ClientConn {
	config := &nconf.RPCClientConfig{Addr: addr}
	config.SetDefaultValues()

	conn, err := dialConn(config)
	if err != nil {
		nlog.Fatalf("fail to dial rpc client connection %s: %s", addr, err)
	}
	return conn
}
