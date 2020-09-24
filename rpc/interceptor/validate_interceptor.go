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
