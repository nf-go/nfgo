package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// ChainUnaryServerInterceptor -
func ChainUnaryServerInterceptor(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)

	if n > 1 {
		lastI := n - 1
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			var (
				chainHandler grpc.UnaryHandler
				curI         int
			)

			chainHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				if curI == lastI {
					return handler(currentCtx, currentReq)
				}
				curI++
				resp, err := interceptors[curI](currentCtx, currentReq, info, chainHandler)
				curI--
				return resp, err
			}

			return interceptors[0](ctx, req, info, chainHandler)
		}
	}

	if n == 1 {
		return interceptors[0]
	}

	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

// ChainStreamingServerInterceptor -
func ChainStreamingServerInterceptor(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	n := len(interceptors)

	if n > 1 {
		lastI := n - 1
		return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			var (
				chainHandler grpc.StreamHandler
				curI         int
			)

			chainHandler = func(currentSrv interface{}, currentStream grpc.ServerStream) error {
				if curI == lastI {
					return handler(currentSrv, currentStream)
				}
				curI++
				err := interceptors[curI](currentSrv, currentStream, info, chainHandler)
				curI--
				return err
			}

			return interceptors[0](srv, stream, info, chainHandler)
		}
	}

	if n == 1 {
		return interceptors[0]
	}

	return func(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, stream)
	}
}
