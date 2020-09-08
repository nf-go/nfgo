package ngrace

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
)

// GraceTerminationServer -
type GraceTerminationServer interface {
	MustServe()
}

// Server -
type Server interface {
	Run() error
	MustRun()
	Shutdown(ctx context.Context) error
}

// NewGraceTerminationServer -
func NewGraceTerminationServer(config *nconf.Config, servers ...Server) (GraceTerminationServer, error) {
	if config == nil {
		return nil, errors.New("config is nill")
	}
	return &graceTerminationServer{
		servers: servers,
		config:  config,
	}, nil
}

// MustNewGraceTerminationServer -
func MustNewGraceTerminationServer(config *nconf.Config, servers ...Server) GraceTerminationServer {
	server, err := NewGraceTerminationServer(config, servers...)
	if err != nil {
		nlog.Fatal("fail to init grace termination server: ", err)
	}
	return server
}

type graceTerminationServer struct {
	config  *nconf.Config
	servers []Server
}

// Serve -
func (s *graceTerminationServer) MustServe() {
	for _, server := range s.servers {
		go server.MustRun()
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit

	nlog.Info("the server is going to shutdown...")
	timeout := s.config.GraceTermination.GraceTerminationPeriod
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var errs gracefulErrs
	for _, server := range s.servers {
		if err := server.Shutdown(ctx); err != nil {
			errs.addError(err)
		}
	}
	if errs.isNil() {
		nlog.Info("the server is stopped normally.")
	} else {
		nlog.Fatal("the server is forced to stop.", errs)
	}
}

type gracefulErrs struct {
	errs []error
}

func (e *gracefulErrs) addError(err error) {
	e.errs = append(e.errs, err)
}

func (e *gracefulErrs) isNil() bool {
	return len(e.errs) == 0
}

func (e *gracefulErrs) Error() string {
	var buf bytes.Buffer
	for _, err := range e.errs {
		buf.WriteString(err.Error())
		buf.WriteString(". ")
	}
	return buf.String()
}
