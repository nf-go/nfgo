package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
			return grpc.Errorf(codes.Code(bizErr.Code()), bizErr.Msg())
		}
		nlog.Logger(ctx).WithError(err).Error()
		return grpc.Errorf(codes.Internal, nerrors.ErrInternal.Msg())
	}
	return nil
}
