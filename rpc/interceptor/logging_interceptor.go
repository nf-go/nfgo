package interceptor

import (
	"context"

	"github.com/sirupsen/logrus"
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
		} else if logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
			logger.WithField("resp", resp).Debug()
		}

	}()

	nlog.Logger(ctx).WithField("req", req).Info()
	return handler(ctx, req)
}

// LoggingStreamServerInterceptor -
func LoggingStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logger := nlog.Logger(stream.Context())
	logger.WithField("req", srv).Info()
	return handler(srv, stream)
}
