package interceptor

import (
	"context"

	"nfgo.ga/nfgo/nconst"
	"nfgo.ga/nfgo/ncontext"
	"nfgo.ga/nfgo/nutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MDCBindingUnaryServerInterceptor -
func MDCBindingUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if ctx, err = bindMDCToContext(ctx, info.FullMethod); err != nil {
		return
	}
	return handler(ctx, req)
}

// MDCBindingStreamServerInterceptor -
func MDCBindingStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	s := &nstream{
		Stream: stream,
	}

	ctx, cancel := context.WithCancel(context.Background())

	s.Cancel = cancel
	ctx, err := bindMDCToContext(ctx, info.FullMethod)
	if err != nil {
		return err
	}
	s.SetContext(ctx)

	return handler(srv, s)
}

func getHeader(md metadata.MD, name string) string {
	values := md.Get(name)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func bindMDCToContext(ctx context.Context, fullMethodName string) (context.Context, error) {
	var traceID string
	var clinetIP string
	var clientType string
	var subject string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		traceID = getHeader(md, nconst.HeaderTraceID)
		clinetIP = getHeader(md, nconst.HeaderRealIP)
		clientType = getHeader(md, nconst.HeaderClientType)
		subject = getHeader(md, nconst.HeaderSub)
	}
	if traceID == "" {
		var err error
		if traceID, err = nutil.UUID(); err != nil {
			return nil, err
		}
	}

	mdc := ncontext.NewMDC()
	mdc.SetTraceID(traceID)
	mdc.SetClientIP(clinetIP)
	mdc.SetClientType(clientType)
	mdc.SetRPCName(fullMethodName)
	mdc.SetSubjectID(subject)

	return ncontext.BindMDCToContext(ctx, mdc), nil
}
