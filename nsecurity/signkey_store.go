package nsecurity

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/ndb"
	"nfgo.ga/nfgo/nutil"
)

const (
	// RedisKeySignKey - nfgo:signkey:{appType}:{subject}
	RedisKeySignKey nutil.RedisKey = "nfgo:signkey:%s:%s"
)

// SignKeyStore -
type SignKeyStore interface {
	Store(clientType, subject, signKey string) error
	Get(clientType, subject string) (string, error)
}

// NewRedisSignKeyStore -
func NewRedisSignKeyStore(redisOper ndb.RedisOper, securityConfig *nconf.SecurityConfig) SignKeyStore {
	return &redisSignKeyStore{
		redisOper:      redisOper,
		securityConfig: securityConfig,
	}
}

// redisSignKeyStore -
type redisSignKeyStore struct {
	redisOper      ndb.RedisOper
	securityConfig *nconf.SecurityConfig
}

func (s *redisSignKeyStore) Store(clientType, subject, signKey string) error {
	key := RedisKeySignKey.Key(clientType, subject)
	return s.redisOper.SetStringOpts(key, signKey, false, false, s.securityConfig.SignKeyLifeTime)
}

func (s *redisSignKeyStore) Get(clientType, subject string) (signKey string, err error) {
	conn := s.redisOper.Conn()
	defer conn.Close()

	signKeyRedisKey := RedisKeySignKey.Key(clientType, subject)
	if s.securityConfig.RefreshSignKeyLife {
		conn.Send("EXPIRE", signKeyRedisKey, int64(s.securityConfig.SignKeyLifeTime/time.Second))
		conn.Send("GET", signKeyRedisKey)
		conn.Flush()
		if _, err = conn.Receive(); err != nil {
			return
		}
		signKey, err = redis.String(conn.Receive())
	} else {
		signKey, err = redis.String(conn.Do("GET", signKeyRedisKey))
	}

	if err == redis.ErrNil {
		return "", nil
	}

	return
}
