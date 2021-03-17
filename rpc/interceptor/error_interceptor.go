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
	"google.golang.org/grpc/status"
	"nfgo.ga/nfgo/nerrors"
	"nfgo.ga/nfgo/nlog"
)

// ErrorHandleUnaryServerInterceptor -
func ErrorHandleUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	err = handleServerError(ctx, err)
	return
}

// ErrorHandleStreamServerInterceptor -
func ErrorHandleStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, stream)
	err = handleServerError(stream.Context(), err)
	return err
}

func handleServerError(ctx context.Context, err error) error {
	if err != nil {
		if bizErr, ok := err.(nerrors.BizError); ok {
			return status.Errorf(codes.Code(bizErr.Code()), bizErr.Msg())
		}
		nlog.Logger(ctx).WithError(err).Error()
		return status.Errorf(codes.Internal, nerrors.ErrInternal.Msg())
	}
	return nil
}
