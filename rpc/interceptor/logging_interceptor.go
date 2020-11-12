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
	"fmt"

	"google.golang.org/grpc"
	"nfgo.ga/nfgo/nerrors"
	"nfgo.ga/nfgo/nlog"
)

// LoggingUnaryServerInterceptor -
func LoggingUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	defer func() {
		logger := nlog.Logger(ctx)
		if err != nil {
			errLogger := logger.WithError(err)
			if _, ok := err.(nerrors.BizError); ok {
				errLogger.Info()
			} else {
				errLogger.Error()
			}
		} else if logger.IsLevelEnabled(nlog.DebugLevel) {
			if stringer, ok := resp.(fmt.Stringer); ok {
				logger.WithField(fieldNameResp, stringer.String()).Debug()
			}
		}

	}()

	if stringer, ok := req.(fmt.Stringer); ok {
		nlog.Logger(ctx).WithField(fieldNameReq, stringer.String()).Info()
	}
	return handler(ctx, req)
}

// LoggingStreamServerInterceptor -
func LoggingStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	defer func() {
		nlog.Logger(stream.Context()).Info("server stream end.")
	}()

	nlog.Logger(stream.Context()).Info("server stream begin.")
	s := &serverStreamWrapper{
		stream: stream,
		ctx:    stream.Context(),
		logMsg: true,
	}
	return handler(srv, s)
}

// LoggingUnaryClientInterceptor -
func LoggingUnaryClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	defer func() {
		logger := nlog.Logger(ctx)
		if err != nil {
			errLogger := logger.WithError(err)
			if _, ok := err.(nerrors.BizError); ok {
				errLogger.Info()
			} else {
				errLogger.Error()
			}
		} else if logger.IsLevelEnabled(nlog.DebugLevel) {
			if stringer, ok := reply.(fmt.Stringer); ok {
				logger.WithFields(nlog.Fields{
					fieldNameRPCCall: method,
					fieldNameResp:    stringer.String(),
				}).Debug()
			}
		}

	}()
	if stringer, ok := req.(fmt.Stringer); ok {
		nlog.Logger(ctx).WithFields(nlog.Fields{
			fieldNameReq:     stringer.String(),
			fieldNameRPCCall: method,
		}).Info()
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// LoggingStreamClientInterceptor -
func LoggingStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

	nlog.Logger(ctx).WithField(fieldNameRPCCall, method).Info("client stream begin.")
	stream, err := streamer(ctx, desc, cc, method, opts...)
	if err == nil {
		stream = &clientStreamWrapper{
			stream: stream,
			logMsg: true,
			method: method,
		}
	}
	return stream, err
}
