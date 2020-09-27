package nfgo

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/graceful"
)

// Server -
type Server interface {
	MustServe()
	RegisterOnShutdown(f func())
}

// NewServer -
func NewServer(config *nconf.Config, servers ...graceful.ShutdownServer) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nill")
	}
	return &nfgoServer{
		servers: servers,
		config:  config,
	}, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config, servers ...graceful.ShutdownServer) Server {
	server, err := NewServer(config, servers...)
	if err != nil {
		nlog.Fatal("fail to init grace termination server: ", err)
	}
	return server
}

type nfgoServer struct {
	config     *nconf.Config
	servers    []graceful.ShutdownServer
	onShutdown []func()
	mu         sync.Mutex
}

func (s *nfgoServer) RegisterOnShutdown(f func()) {
	s.mu.Lock()
	s.onShutdown = append(s.onShutdown, f)
	s.mu.Unlock()
}

// Serve -
func (s *nfgoServer) MustServe() {
	for _, server := range s.servers {
		go server.MustServe()
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
	for _, f := range s.onShutdown {
		f()
	}
	if errs.isNil() {
		nlog.Info("the server is stopped normally.")
	} else {
		nlog.Error("the server is forced to stop.", errs)
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
	var sb strings.Builder
	for _, err := range e.errs {
		sb.WriteString(err.Error())
		sb.WriteString(". ")
	}
	return sb.String()
}
