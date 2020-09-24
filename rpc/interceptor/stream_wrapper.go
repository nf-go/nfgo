package interceptor

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"nfgo.ga/nfgo/nlog"
)

const (
	fieldNameReq     = "req"
	fieldNameResp    = "resp"
	fieldNameRPCCall = "rpcCall"
)

// serverStreamWrapper -
type serverStreamWrapper struct {
	stream      grpc.ServerStream
	ctx         context.Context
	logMsg      bool
	validateMsg bool
}

func (s *serverStreamWrapper) Context() context.Context {
	return s.ctx
}

func (s *serverStreamWrapper) SetHeader(md metadata.MD) error {
	return s.stream.SetHeader(md)
}

func (s *serverStreamWrapper) SendHeader(md metadata.MD) error {
	return s.stream.SendHeader(md)
}

func (s *serverStreamWrapper) SetTrailer(md metadata.MD) {
	s.stream.SetTrailer(md)
}

func (s *serverStreamWrapper) SendMsg(m interface{}) (err error) {
	if s.logMsg {
		if stringer, ok := m.(fmt.Stringer); ok {
			nlog.Logger(s.ctx).WithField("resp", stringer.String()).Info("server stream send msg.")
		}
	}
	return s.stream.SendMsg(m)
}

func (s *serverStreamWrapper) RecvMsg(m interface{}) error {
	err := s.stream.RecvMsg(m)

	if s.logMsg {
		if stringer, ok := m.(fmt.Stringer); ok {
			nlog.Logger(s.ctx).WithField(fieldNameReq, stringer.String()).Info("server stream recv msg.")
		}
	}

	if s.validateMsg && err == nil {
		if v, ok := m.(interface{ Validate() error }); ok {
			if err = v.Validate(); err != nil {
				err = grpc.Errorf(codes.InvalidArgument, err.Error())
			}
		}
	}

	return err
}

// clientStreamWrapper -
type clientStreamWrapper struct {
	stream      grpc.ClientStream
	logMsg      bool
	validateMsg bool
	method      string
}

func (s *clientStreamWrapper) Header() (metadata.MD, error) {
	return s.stream.Header()
}

func (s *clientStreamWrapper) Trailer() metadata.MD {
	return s.stream.Trailer()
}

func (s *clientStreamWrapper) CloseSend() error {
	return s.stream.CloseSend()
}

func (s *clientStreamWrapper) Context() context.Context {
	return s.stream.Context()
}

func (s *clientStreamWrapper) SendMsg(m interface{}) error {
	if s.logMsg {
		if stringer, ok := m.(fmt.Stringer); ok {
			nlog.Logger(s.Context()).WithFields(nlog.Fields{
				fieldNameReq:     stringer.String(),
				fieldNameRPCCall: s.method,
			}).Info("clent stream send msg.")
		}
	}
	if s.validateMsg {
		if v, ok := m.(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return grpc.Errorf(codes.InvalidArgument, err.Error())
			}
		}
	}
	return s.stream.SendMsg(m)
}

func (s *clientStreamWrapper) RecvMsg(m interface{}) error {
	if s.logMsg {
		if stringer, ok := m.(fmt.Stringer); ok {
			nlog.Logger(s.Context()).WithFields(nlog.Fields{
				fieldNameResp:    stringer.String(),
				fieldNameRPCCall: s.method,
			}).Info("client stream recv msg.")
		}
		err := s.stream.RecvMsg(m)
		if err == io.EOF {
			nlog.Logger(s.Context()).Info("client stream end.")
		}
		return err
	}
	return s.stream.RecvMsg(m)
}
