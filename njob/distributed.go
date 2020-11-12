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
