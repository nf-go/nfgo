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

package njob

import (
	"context"
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nerrors"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/graceful"
	"nfgo.ga/nfgo/nutil/ntypes"
)

// Server -
type Server interface {
	graceful.ShutdownServer
}

// JobServer -
type jobServer struct {
	config *nconf.Config
	opts   *serverOptions
	c      *cron.Cron
	stop   chan struct{}
}

// NewServer -
func NewServer(config *nconf.Config, opt ...ServerOption) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	cronConfig := config.CronConfig
	if cronConfig == nil {
		return nil, errors.New("config.cron is nil")
	}

	opts := &serverOptions{
		funcJobs: map[string]FuncJob{},
		jobs:     map[string]Job{},
	}
	for _, o := range opt {
		o(opts)
	}

	jobWrappers := []cron.JobWrapper{cron.Recover(defaultLogger)}
	if ntypes.BoolValue(cronConfig.SkipIfStillRunning) {
		jobWrappers = append(jobWrappers, cron.SkipIfStillRunning(defaultLogger))
	}

	c := cron.New(cron.WithChain(jobWrappers...))
	return &jobServer{
		opts:   opts,
		config: config,
		c:      c,
		stop:   make(chan struct{}),
	}, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config, opt ...ServerOption) Server {
	server, err := NewServer(config, opt...)
	if err != nil {
		nlog.Fatal("fail to init job server: ", err)
	}
	return server
}

// Serve -
func (s *jobServer) Serve() error {
	if err := s.addJobs(); err != nil {
		return err
	}
	s.c.Start()
	<-s.stop
	return nil
}

func (s *jobServer) addJobs() error {
	cronConf := s.config.CronConfig
	for _, conf := range cronConf.CronJobs {
		if fn, ok := s.opts.funcJobs[conf.Name]; ok {
			if err := s.addJob(conf, fn); err != nil {
				return fmt.Errorf("fail to init croJob %s: %w", conf.Name, err)
			}
		} else if job, ok := s.opts.jobs[conf.Name]; ok {
			if err := s.addJob(conf, job); err != nil {
				return fmt.Errorf("fail to init croJob %s: %w", conf.Name, err)
			}
		} else {
			nlog.Warnf("Please provide a job or a jobfunc for %s, use njob.JobsOption or njob.JobFuncsOption", conf.Name)
		}

	}
	return nil
}

func (s *jobServer) addJob(conf *nconf.CronJobConfig, j Job) error {
	fn := func() {
		ctx := newJobContext(conf.Name)
		if err := j.Run(ctx); err != nil {
			logger := nlog.Logger(ctx).WithError(err)
			if bizErr, ok := err.(nerrors.BizError); ok {
				logger.Infof("job %s completed but an biz error occurred %s", conf.Name, bizErr)
				return
			}
			logger.Errorf("job %s completed but an error occurred %s", conf.Name, err)
		}
	}
	var job cron.Job = cron.FuncJob(fn)
	if s.opts.distributedMutex != nil {
		job = cron.NewChain(
			distributedRunning(s.config, conf.Name, s.opts.distributedMutex),
		).Then(job)
	}
	_, err := s.c.AddJob(conf.Schedule, job)
	return err
}

// MustServe -
func (s *jobServer) MustServe() {
	if err := s.Serve(); err != nil {
		nlog.Fatal("fail to start job server:", err)
	}
}

// Shutdown -
func (s *jobServer) Shutdown(ctx context.Context) error {
	// TODO: wait for running jobs to complete
	s.c.Stop()
	s.stop <- struct{}{}
	return nil
}
