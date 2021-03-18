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

package nfgo

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nerrors"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/graceful"
)

// Server -
type Server interface {
	MustServe()
	RegisterOnShutdown(f func() error)
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
	onShutdown []func() error
	mu         sync.Mutex
}

func (s *nfgoServer) RegisterOnShutdown(f func() error) {
	s.mu.Lock()
	s.onShutdown = append(s.onShutdown, f)
	s.mu.Unlock()
}

// Serve -
func (s *nfgoServer) MustServe() {
	s.autoSetMaxProcs()
	s.serve()
	s.shutdown()
}

func (s *nfgoServer) autoSetMaxProcs() {
	if maxProcs := s.config.App.GOMAXPROCS; maxProcs > 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		undo, err := maxprocs.Set()
		defer undo()
		if err != nil {
			nlog.Fatal("fail to auto set max procs", err)
		}
	}
	nlog.Infof("auto max procs, procs=%d", runtime.GOMAXPROCS(-1))
}

func (s *nfgoServer) serve() {
	for _, server := range s.servers {
		go server.MustServe()
	}
	// Wait for interrupt signal to shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
}

func (s *nfgoServer) shutdown() {
	nlog.Info("the server is going to shutdown...")
	timeout := s.config.GraceTermination.GraceTerminationPeriod
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cleaned := make(chan error)
	go s.cleanup(ctx, cleaned)

	select {
	case err := <-cleaned:
		if err != nil {
			nlog.Error("the server is stopped with error: ", err)
		} else {
			nlog.Info("the server is stopped normally.")
		}
	case <-ctx.Done():
		nlog.Error("shutdown timeout, the server is stopped forcibly.")
	}
}

func (s *nfgoServer) cleanup(ctx context.Context, cleaned chan<- error) {
	var err error
	for _, server := range s.servers {
		err = nerrors.Append(err, server.Shutdown(ctx))
	}
	for _, f := range s.onShutdown {
		err = nerrors.Append(err, f())
	}
	cleaned <- err
}
