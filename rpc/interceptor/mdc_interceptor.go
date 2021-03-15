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

	"nfgo.ga/nfgo/nconst"
	"nfgo.ga/nfgo/ncontext"
	"nfgo.ga/nfgo/nutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MDCBindingUnaryClientInterceptor -
func MDCBindingUnaryClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if mdc, err := ncontext.CurrentMDC(ctx); err == nil {
		kv := []string{
			nconst.HeaderTraceID, mdc.TraceID(),
			nconst.HeaderRealIP, mdc.ClientIP(),
			nconst.HeaderClientType, mdc.ClientType(),
			nconst.HeaderSub, mdc.SubjectID(),
		}
		ctx = metadata.AppendToOutgoingContext(ctx, kv...)
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

// MDCBindingStreamClientInterceptor -
func MDCBindingStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if mdc, err := ncontext.CurrentMDC(ctx); err == nil {
		kv := []string{
			nconst.HeaderTraceID, mdc.TraceID(),
			nconst.HeaderRealIP, mdc.ClientIP(),
			nconst.HeaderClientType, mdc.ClientType(),
			nconst.HeaderSub, mdc.SubjectID(),
		}
		ctx = metadata.AppendToOutgoingContext(ctx, kv...)
	}
	return streamer(ctx, desc, cc, method, opts...)
}

// MDCBindingUnaryServerInterceptor -
func MDCBindingUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if ctx, err = bindMDCToContext(ctx, info.FullMethod); err != nil {
		return
	}
	return handler(ctx, req)
}

// MDCBindingStreamServerInterceptor -
func MDCBindingStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	s := &serverStreamWrapper{stream: stream}
	ctx, err := bindMDCToContext(stream.Context(), info.FullMethod)
	if err != nil {
		return err
	}
	s.ctx = ctx
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

	return ncontext.WithMDC(ctx, mdc), nil
}
