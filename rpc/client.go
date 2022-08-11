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
	"crypto/tls"
	"errors"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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
			interceptor.LoggingUnaryClientInterceptor,
		),
		grpc.WithChainStreamInterceptor(
			interceptor.MDCBindingStreamClientInterceptor,
			interceptor.ValidateStreamClientInterceptor,
			interceptor.LoggingStreamClientInterceptor,
		),
	}
	if ntypes.BoolValue(config.Plaintext) {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		})
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	}

	// Dial connection
	return grpc.Dial(config.Addr, dialOptions...)
}

// DialClientConnPlaintext -
func DialClientConnPlaintext(addr string) (*grpc.ClientConn, error) {
	config := &nconf.RPCClientConfig{Addr: addr}
	config.SetDefaultValues()
	return dialConn(config)
}

// MustDialClientConnPlaintext -
func MustDialClientConnPlaintext(addr string) *grpc.ClientConn {
	conn, err := DialClientConnPlaintext(addr)
	if err != nil {
		nlog.Fatalf("fail to dial rpc client connection %s: %s", addr, err)
	}
	return conn
}
