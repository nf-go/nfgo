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

package ndb

import (
	"context"
	"errors"
	"fmt"
	"time"

	sentinel "github.com/FZambia/sentinel/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"github.com/nf-go/nfgo/nconf"
	"github.com/nf-go/nfgo/nerrors"
	"github.com/nf-go/nfgo/nlog"
)

// RedisPool -
type RedisPool interface {
	Get() redis.Conn
	Close() error
}

// NewRedisPool -
func NewRedisPool(redisConfig *nconf.RedisConfig) (RedisPool, error) {
	if redisConfig.Sentinel != nil {
		return newSentinelRedisPool(redisConfig)
	}
	if redisConfig.Cluster != nil {
		return newClusterRedisPool(redisConfig)
	}
	return newRedisPool(redisConfig)
}

// MustNewRedisPool -
func MustNewRedisPool(redisConfig *nconf.RedisConfig) RedisPool {
	pool, err := NewRedisPool(redisConfig)
	if err != nil {
		nlog.Fatal("fail to new redis pool: ", err)
	}
	return pool
}

func diaContextFunc(addr, pass string, database uint8) func(ctx context.Context) (redis.Conn, error) {
	return func(ctx context.Context) (redis.Conn, error) {
		dialOptions := make([]redis.DialOption, 0, 2)
		if pass != "" {
			dialOptions = append(dialOptions, redis.DialPassword(pass))
		}
		if database > 0 {
			dialOptions = append(dialOptions, redis.DialDatabase(int(database)))
		}
		return redis.DialContext(ctx, "tcp", addr, dialOptions...)
	}
}

func testOnBorrow(conn redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := conn.Do("PING")
	return err
}

func newPool(redisConfig *nconf.RedisConfig) (*redis.Pool, error) {
	if redisConfig == nil {
		return nil, errors.New("redisConfig is nil")
	}
	return &redis.Pool{
		MaxIdle:         int(redisConfig.MaxIdle),
		MaxActive:       int(redisConfig.MaxActive),
		IdleTimeout:     redisConfig.IdleTimeout,
		MaxConnLifetime: redisConfig.MaxConnLifetime,
		Wait:            false,
	}, nil
}

// newRedisPool  - stand-alone redis
func newRedisPool(redisConfig *nconf.RedisConfig) (*redis.Pool, error) {
	redisPool, err := newPool(redisConfig)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	redisPool.DialContext = diaContextFunc(addr, redisConfig.Password, redisConfig.Database)

	if redisConfig.TestOnBorrow {
		redisPool.TestOnBorrow = testOnBorrow
	}
	return redisPool, nil
}

// newSentinelRedisPool - sentinel redis
func newSentinelRedisPool(redisConfig *nconf.RedisConfig) (*redis.Pool, error) {
	redisPool, err := newPool(redisConfig)
	if err != nil {
		return nil, err
	}
	sentinelConfig := redisConfig.Sentinel
	if sentinelConfig == nil {
		return nil, errors.New("redisConfig's sentinel config is nil")
	}
	sntnl := &sentinel.Sentinel{
		Addrs:      sentinelConfig.Addrs,
		MasterName: sentinelConfig.Master,
		Dial: func(addr string) (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	redisPool.DialContext = func(ctx context.Context) (redis.Conn, error) {
		masterAddr, err := sntnl.MasterAddr()
		if err != nil {
			return nil, err
		}
		return diaContextFunc(masterAddr, redisConfig.Password, redisConfig.Database)(ctx)
	}

	if redisConfig.TestOnBorrow {
		redisPool.TestOnBorrow = func(c redis.Conn, t time.Time) error {
			ok, err := sentinel.TestRole(c, "master")
			if err != nil {
				return nerrors.WithStack(err)
			}
			if !ok {
				return errors.New("role check failed")
			}
			return nil
		}
	}

	return redisPool, nil
}

// newClusterRedisPool - redis cluster
func newClusterRedisPool(redisConfig *nconf.RedisConfig) (*redisc.Cluster, error) {
	dialOptions := []redis.DialOption{}
	if redisConfig.Password != "" {
		dialOptions = append(dialOptions, redis.DialPassword(redisConfig.Password))
	}
	cluster := &redisc.Cluster{
		StartupNodes: redisConfig.Cluster.Addrs,
		DialOptions:  dialOptions,
		CreatePool: func(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
			pool, err := newPool(redisConfig)
			if err != nil {
				return nil, err
			}
			pool.Dial = func() (redis.Conn, error) {
				return redis.Dial("tcp", addr, opts...)
			}
			if redisConfig.TestOnBorrow {
				pool.TestOnBorrow = testOnBorrow
			}
			return pool, nil
		},
	}
	if err := cluster.Refresh(); err != nil {
		return nil, err
	}
	return cluster, nil
}
