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
