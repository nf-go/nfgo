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

	"github.com/gomodule/redigo/redis"
	"github.com/robfig/cron/v3"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/ndb"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/ncrypto"
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

type redisMutex struct {
	redisOper   ndb.RedisOper
	lockTimeout time.Duration
}

// NewRedisDistributedMutex -
// This is just an example, and it is recommended to use ETCD for distributed lock.
// In addition, This pattern is discouraged in favor of the Redlock(https://redis.io/topics/distlock) algorithm which is only a bit more complex to implement,
// but offers better guarantees and is fault tolerant.
func NewRedisDistributedMutex(redisOper ndb.RedisOper, lockTimeout time.Duration) DistributedMutex {
	return &redisMutex{
		redisOper:   redisOper,
		lockTimeout: lockTimeout,
	}
}

func (m *redisMutex) TryRunWithMutex(key string, timeout time.Duration, fn func()) error {
	token, err := ncrypto.UUID()
	if err != nil {
		return err
	}

	conn := m.redisOper.Conn()
	_, err = redis.String(conn.Do("SET", key, token, "PX", int64(m.lockTimeout/time.Millisecond), "NX"))
	conn.Close()
	if err != nil {
		if err == redis.ErrNil {
			return nil
		}
		return err
	}

	defer func() {
		// release the lock
		if err := m.redisOper.DeleteByKeyValue(key, token); err != nil {
			nlog.Errorf("fail to release lock %s: %s", key, err)
		}
	}()

	// acquired the lock and run the job
	fn()

	return nil
}
