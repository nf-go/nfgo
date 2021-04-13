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
	"google.golang.org/grpc"
	"nfgo.ga/nfgo/nmetrics"
	"nfgo.ga/nfgo/rpc/interceptor"
)

type serverOptions struct {
	metricsServer            nmetrics.Server
	unaryServerInterceptors  []grpc.UnaryServerInterceptor
	streamServerInterceptors []grpc.StreamServerInterceptor
}

func (opts *serverOptions) setInterceptors() {
	unaryInterceptors := []grpc.UnaryServerInterceptor{interceptor.RecoverUnaryServerInterceptor}
	if opts.metricsServer != nil {
		unaryInterceptors = append(unaryInterceptors, opts.metricsServer.GrpcMetricsUnaryServerInterceptor())
	}
	unaryInterceptors = append(unaryInterceptors,
		interceptor.RecoverUnaryServerInterceptor,
		interceptor.MDCBindingUnaryServerInterceptor,
		interceptor.ValidateUnaryServerInterceptor,
		interceptor.LoggingUnaryServerInterceptor,
		interceptor.ErrorHandleUnaryServerInterceptor)
	if len(opts.unaryServerInterceptors) == 0 {
		opts.unaryServerInterceptors = unaryInterceptors
	} else {
		opts.unaryServerInterceptors = append(unaryInterceptors, opts.unaryServerInterceptors...)
	}

	streamInterceptors := []grpc.StreamServerInterceptor{interceptor.RecoverStreamServerInterceptor}
	if opts.metricsServer != nil {
		streamInterceptors = append(streamInterceptors, opts.metricsServer.GrpcMetricsStramServerInterceptor())
	}
	streamInterceptors = append(streamInterceptors,
		interceptor.MDCBindingStreamServerInterceptor,
		interceptor.ValidateStreamServerInterceptor,
		interceptor.LoggingStreamServerInterceptor,
		interceptor.ErrorHandleStreamServerInterceptor)

	if len(opts.streamServerInterceptors) == 0 {
		opts.streamServerInterceptors = streamInterceptors
	} else {
		opts.streamServerInterceptors = append(streamInterceptors, opts.streamServerInterceptors...)
	}
}

// ServerOption -
type ServerOption func(*serverOptions)

// MetricsServerOption -
func MetricsServerOption(s nmetrics.Server) ServerOption {
	return func(opts *serverOptions) {
		opts.metricsServer = s
	}
}

// UnaryServerInterceptorOption -
func UnaryServerInterceptorOption(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(opts *serverOptions) {
		opts.unaryServerInterceptors = interceptors
	}
}

// StreamServerInterceptorOption -
func StreamServerInterceptorOption(interceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(opts *serverOptions) {
		opts.streamServerInterceptors = interceptors
	}
}
