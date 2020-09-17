package njob

import (
	"context"
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/ngrace"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/ntypes"
)

// Server -
type Server interface {
	ngrace.Server
}

// ServerOption -
type ServerOption struct {
	JobFuncs         map[string]func()
	DistributedMutex DistributedMutex
}

// JobServer -
type jobServer struct {
	config *nconf.Config
	option *ServerOption
	c      *cron.Cron
	stop   chan struct{}
}

// NewServer -
func NewServer(config *nconf.Config, option *ServerOption) (Server, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	cronConfig := config.CronConfig
	if cronConfig == nil {
		return nil, errors.New("config.cron is nil")
	}
	if option == nil {
		option = &ServerOption{
			JobFuncs: map[string]func(){},
		}
	}

	jobWrappers := []cron.JobWrapper{cron.Recover(defaultLogger)}
	if ntypes.BoolValue(cronConfig.SkipIfStillRunning) {
		jobWrappers = append(jobWrappers, cron.SkipIfStillRunning(defaultLogger))
	}

	c := cron.New(cron.WithChain(jobWrappers...))
	return &jobServer{
		option: option,
		config: config,
		c:      c,
		stop:   make(chan struct{}),
	}, nil
}

// MustNewServer -
func MustNewServer(config *nconf.Config, option *ServerOption) Server {
	server, err := NewServer(config, option)
	if err != nil {
		nlog.Fatal("fail to init job server: ", err)
	}
	return server
}

// Serve -
func (s *jobServer) Serve() error {
	if err := s.addJobFuncs(); err != nil {
		return err
	}
	s.c.Start()
	<-s.stop
	return nil
}

func (s *jobServer) addJobFuncs() error {
	cronConf := s.config.CronConfig
	for _, conf := range cronConf.CronJobs {
		if fn, ok := s.option.JobFuncs[conf.Name]; ok {
			if err := s.addJobFunc(conf, fn); err != nil {
				return fmt.Errorf("fail to init croJob %s: %w", conf.Name, err)
			}
		} else {
			nlog.Warnf("Please provide a job func: %s, use serverOption.JobFuncs", conf.Name)
		}

	}
	return nil
}

func (s *jobServer) addJobFunc(conf *nconf.CronJobConfig, fn func()) error {
	var job cron.Job = cron.FuncJob(fn)
	if s.option.DistributedMutex != nil {
		job = cron.NewChain(
			distributedRunning(s.config, conf.Name, s.option.DistributedMutex),
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
