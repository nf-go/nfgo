package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"nfgo.ga/nfgo/nerrors"
)

// ErrorHandleUnaryServerInterceptor -
func ErrorHandleUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	err = handleError(err)
	return
}

// ErrorHandleStreamServerInterceptor -
func ErrorHandleStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, stream)
	err = handleError(err)
	return err
}

func handleError(err error) error {
	if err != nil {
		// 处理业务错误
		if bizErr, ok := err.(nerrors.BizError); ok {
			return grpc.Errorf(codes.Code(bizErr.Code()), bizErr.Msg())
		}
		return grpc.Errorf(codes.Internal, nerrors.ErrInternal.Msg())
	}
	return nil
}
