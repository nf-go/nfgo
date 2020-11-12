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

package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ValidateUnaryClientInterceptor -
func ValidateUnaryClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if v, ok := req.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return grpc.Errorf(codes.InvalidArgument, err.Error())
		}
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// ValidateStreamClientInterceptor -
func ValidateStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	stream, err := streamer(ctx, desc, cc, method, opts...)
	if err == nil {
		stream = &clientStreamWrapper{
			stream:      stream,
			validateMsg: true,
			method:      method,
		}
	}
	return stream, err
}

// ValidateUnaryServerInterceptor -
func ValidateUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if v, ok := req.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}

// ValidateStreamServerInterceptor -
func ValidateStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	s := &serverStreamWrapper{
		stream:      stream,
		ctx:         stream.Context(),
		validateMsg: true,
	}
	return handler(srv, s)
}
