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
				logger.WithField("resp", stringer.String()).Debug()
			}
		}

	}()

	if stringer, ok := req.(fmt.Stringer); ok {
		nlog.Logger(ctx).WithField("req", stringer.String()).Info()
	}
	return handler(ctx, req)
}

// LoggingStreamServerInterceptor -
func LoggingStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logger := nlog.Logger(stream.Context())
	logger.WithField("req", srv).Info()
	return handler(srv, stream)
}
