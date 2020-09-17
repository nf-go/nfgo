package njob

import (
	"time"

	"github.com/robfig/cron/v3"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
)

// DistributedMutex -
type DistributedMutex interface {
	TryRunWithMutex(key string, timeout time.Duration, fn func()) error
}

func distributedRunning(conf *nconf.Config, jobName string, mutex DistributedMutex) cron.JobWrapper {

	return func(job cron.Job) cron.Job {

		return cron.FuncJob(func() {
			key := "/nfgo/njob/" + conf.App.Name + "/" + jobName
			err := mutex.TryRunWithMutex(key, time.Second, job.Run)
			if err != nil {
				nlog.Errorf("fail to try run job %s with mutex %s: %s", jobName, key, err)
				return
			}
		})
	}

}
