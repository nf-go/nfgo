package ndb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"nfgo.ga/nfgo/nconf"
)

func diaContextFunc(addr, pass string, database uint8) func(ctx context.Context) (redis.Conn, error) {
	return func(ctx context.Context) (redis.Conn, error) {
		conn, err := redis.DialContext(ctx, "tcp", addr)
		if err != nil {
			return nil, err
		}

		if pass != "" {
			if _, err := conn.Do("AUTH", pass); err != nil {
				conn.Close()
				return nil, err
			}
		}

		if _, err := conn.Do("SELECT", database); err != nil {
			conn.Close()
			return nil, err
		}

		return conn, nil
	}
}

func testOnBorrowFunc() func(c redis.Conn, t time.Time) error {
	return func(conn redis.Conn, t time.Time) error {
		if time.Since(t) < time.Minute {
			return nil
		}
		_, err := conn.Do("PING")
		return err
	}
}

func newRedisPool(redisConfig *nconf.RedisConfig) (*redis.Pool, error) {
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

// NewRedisPool -
func NewRedisPool(redisConfig *nconf.RedisConfig) (*redis.Pool, error) {
	redisPool, err := newRedisPool(redisConfig)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	redisPool.DialContext = diaContextFunc(addr, redisConfig.Password, redisConfig.Database)

	if redisConfig.TestOnBorrow {
		redisPool.TestOnBorrow = testOnBorrowFunc()
	}
	return redisPool, nil
}

// NewSentinelRedisPool -
func NewSentinelRedisPool(redisConfig *nconf.RedisConfig) (*redis.Pool, error) {
	redisPool, err := newRedisPool(redisConfig)
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
			if !sentinel.TestRole(c, "master") {
				return errors.New("Role check failed")
			}
			return nil
		}
	}

	return redisPool, nil
}
