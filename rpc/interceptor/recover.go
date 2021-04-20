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
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"nfgo.ga/nfgo/nlog"
)

// RecoverUnaryServerInterceptor -
func RecoverUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	panicked := true
	defer func() {
		if r := recover(); r != nil || panicked {
			err = recoverFrom(ctx, r)
		}
	}()
	resp, err = handler(ctx, req)
	panicked = false
	return
}

// RecoverStreamServerInterceptor -
func RecoverStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	panicked := true
	defer func() {
		if r := recover(); r != nil || panicked {
			err = recoverFrom(stream.Context(), r)
		}
	}()
	err = handler(srv, stream)
	panicked = false
	return err
}

func recoverFrom(ctx context.Context, r interface{}) error {
	stackTrace := debug.Stack()
	nlog.Logger(ctx).Error(string(stackTrace))
	return status.Errorf(codes.Unknown, "panic triggered: %v", r)
}
