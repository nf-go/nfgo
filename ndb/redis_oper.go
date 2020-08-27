package ndb

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/gomodule/redigo/redis"
)

// NewRedisOper -
func NewRedisOper(redisPool *redis.Pool) RedisOper {
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

	Delete(key string) error
	Deletes(keys ...string) error
}

type redisOperImpl struct {
	redisPool *redis.Pool
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

func (r *redisOperImpl) Delete(key string) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

func (r *redisOperImpl) Deletes(keys ...string) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	if _, err := conn.Do("MULTI"); err != nil {
		return err
	}
	for _, key := range keys {
		if _, err := conn.Do("DEL", key); err != nil {
			return err
		}

	}
	_, err := conn.Do("EXEC")
	return err
}