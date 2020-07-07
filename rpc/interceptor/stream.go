package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type nstream struct {
	Stream grpc.ServerStream
	ctx    context.Context
	Cancel context.CancelFunc
}

func (s *nstream) Context() context.Context {
	defer func() {
		select {
		case <-s.Stream.Context().Done():
			s.Cancel()
			return
		default:
		}
	}()
	return s.ctx
}

func (s *nstream) SetContext(new context.Context) {
	s.ctx = new
}

func (s *nstream) SetHeader(md metadata.MD) error {
	return s.Stream.SetHeader(md)
}

func (s *nstream) SendHeader(md metadata.MD) error {
	return s.Stream.SendHeader(md)
}

func (s *nstream) SetTrailer(md metadata.MD) {
	s.Stream.SetTrailer(md)
}

func (s *nstream) SendMsg(m interface{}) (err error) {
	return s.Stream.SendMsg(m)
}

func (s *nstream) RecvMsg(m interface{}) error {
	return s.Stream.RecvMsg(m)
}
