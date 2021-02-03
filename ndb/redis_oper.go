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
	"bytes"
	"encoding/gob"
	"time"

	"github.com/gomodule/redigo/redis"
)

// NewRedisOper -
func NewRedisOper(redisPool RedisPool) RedisOper {
	return &redisOperImpl{
		redisPool: redisPool,
	}
}

// RedisOper -
type RedisOper interface {
	// Conn - Get a redis connection from the pool.
	// The application must close the returned connection.
	Conn() redis.Conn
	GetString(key string) (string, error)
	SetString(key, value string) error
	SetStringOpts(key, value string, setnx bool, setxx bool, ttl time.Duration) error

	GetObject(key string, model interface{}) (interface{}, error)
	SetObject(key string, value interface{}) error
	SetObjectOpts(key string, value interface{}, setnx bool, setxx bool, ttl time.Duration) error

	Unlink(key string) error
	Unlinks(keys ...string) error

	Delete(key string) error
	Deletes(keys ...string) error
	DeleteByKeyValue(key string, val string) error
}

type redisOperImpl struct {
	redisPool RedisPool
}

func (r *redisOperImpl) Conn() redis.Conn {
	return r.redisPool.Get()
}

func (r *redisOperImpl) GetString(key string) (string, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return "", nil
	}
	return val, err
}

func (r *redisOperImpl) SetString(key, value string) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

func (r *redisOperImpl) SetStringOpts(key, value string, setnx bool, setxx bool, ttl time.Duration) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	args := []interface{}{key, value}
	if ttl > time.Millisecond {
		args = append(args, "PX", int64(ttl/time.Millisecond))
	}
	if setnx {
		args = append(args, "NX")
	}
	if setxx {
		args = append(args, "XX")
	}
	_, err := conn.Do("SET", args...)
	return err
}

func (r *redisOperImpl) GetObject(key string, model interface{}) (interface{}, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(model)
	return model, err
}

func (r *redisOperImpl) SetObject(key string, value interface{}) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	_, err := conn.Do("SET", key, buf.Bytes())
	return err
}

func (r *redisOperImpl) SetObjectOpts(key string, value interface{}, setnx bool, setxx bool, ttl time.Duration) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	args := []interface{}{key, buf.Bytes()}
	if ttl > time.Millisecond {
		args = append(args, "PX", int64(ttl/time.Millisecond))
	}
	if setnx {
		args = append(args, "NX")
	}
	if setxx {
		args = append(args, "XX")
	}
	_, err := conn.Do("SET", args...)
	return err
}

func (r *redisOperImpl) del(delCmd string, key string) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do(delCmd, key)
	return err
}

func (r *redisOperImpl) dels(delCmd string, keys ...string) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	if _, err := conn.Do("MULTI"); err != nil {
		return err
	}
	for _, key := range keys {
		if _, err := conn.Do(delCmd, key); err != nil {
			return err
		}

	}
	_, err := conn.Do("EXEC")
	return err
}

func (r *redisOperImpl) Unlink(key string) error {
	return r.del("UNLINK", key)
}

func (r *redisOperImpl) Unlinks(keys ...string) error {
	return r.dels("UNLINK", keys...)
}

func (r *redisOperImpl) Delete(key string) error {
	return r.del("DEL", key)
}

func (r *redisOperImpl) Deletes(keys ...string) error {
	return r.dels("DEL", keys...)
}

func (r *redisOperImpl) DeleteByKeyValue(key string, val string) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	v, err := redis.String(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil
		}
		return err
	}

	if val == v {
		_, err = conn.Do("DEL", key)
		return err
	}

	return nil
}
