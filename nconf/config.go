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

package nconf

import (
	"time"

	"nfgo.ga/nfgo/nutil/nconst"
)

// Config data
type Config struct {
	App        *AppConfig     `yaml:"app"`
	Log        *LogConfig     `yaml:"log"`
	DB         *DbConfig      `yaml:"db"`
	Redis      *RedisConfig   `yaml:"redis"`
	Web        *WebConfig     `yaml:"web"`
	RPC        *RPCConfig     `yaml:"rpc"`
	CronConfig *CronConfig    `yaml:"cron"`
	Metrics    *MetricsConfig `yaml:"metrics"`
}

// AppConfig -
type AppConfig struct {
	Group            string                  `yaml:"group"`
	Name             string                  `yaml:"name"`
	Profile          string                  `yaml:"profile"`
	GOMAXPROCS       int                     `yaml:"goMaxProcs"`
	Ext              ExtConfig               `yaml:"ext"`
	GraceTermination *GraceTerminationConfig `yaml:"graceTermination"`
}

// IsProfileLocal -
func (c *AppConfig) IsProfileLocal() bool {
	return c.Profile == "" || c.Profile == nconst.ProfileLocal
}

// ExtConfig -
type ExtConfig map[string]interface{}

// StrVal -
func (e ExtConfig) StrVal(key string) string {
	if val, ok := e[key].(string); ok {
		return val
	}
	return ""
}

// IntVal -
func (e ExtConfig) IntVal(key string) int {
	if val, ok := e[key].(int); ok {
		return val
	}
	return 0
}

// BoolVal -
func (e ExtConfig) BoolVal(key string) bool {
	if val, ok := e[key].(bool); ok {
		return val
	}
	return false
}

// LogConfig -
type LogConfig struct {
	Level           string `yaml:"level"`
	Format          string `yaml:"format"`
	PrettyPrint     bool   `yaml:"prettyPrint"`
	CallerPrint     bool   `yaml:"callerPrint"`
	TimestampFormat string `yaml:"timestampFormat"`
	LogPath         string `yaml:"logPath"`
	LogFilename     string `yaml:"logFilename"`
}

// WebConfig -
type WebConfig struct {
	Host               string         `yaml:"host"`
	Port               int32          `yaml:"port"`
	Swagger            *SwaggerConfig `yaml:"swagger"`
	MaxMultipartMemory int64          `yaml:"maxMultipartMemory"`
}

// SwaggerConfig -
type SwaggerConfig struct {
	Enabled bool   `yaml:"enabled"`
	URL     string `yaml:"url"`
}

// RPCConfig -
type RPCConfig struct {
	Host                     string                      `yaml:"host"`
	Port                     int32                       `yaml:"port"`
	MaxRecvMsgSize           int64                       `yaml:"maxRecvMsgSize"`
	RegisterHealthServer     *bool                       `yaml:"registerHealthServer"`
	RegisterReflectionServer *bool                       `yaml:"registerReflectionServer"`
	Clients                  map[string]*RPCClientConfig `yaml:"clients"`
}

// RPCClientConfig -
type RPCClientConfig struct {
	Addr               string `yaml:"addr"`
	Plaintext          *bool  `yaml:"plaintext"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
	MaxCallSendMsgSize int    `yaml:"maxCallSendMsgSize"`
	MaxCallRecvMsgSize int    `yaml:"maxCallRecvMsgSize"`
}

// DbConfig -
type DbConfig struct {
	Username               string        `yaml:"username"`
	Password               string        `yaml:"password"`
	Host                   string        `yaml:"host"`
	Port                   int32         `yaml:"port"`
	Database               string        `yaml:"database"`
	Charset                string        `yaml:"charset"`
	MaxIdle                int32         `yaml:"maxIdle"`
	MaxIdleTime            time.Duration `yaml:"maxIdleTime"`
	MaxOpen                int32         `yaml:"maxOpen"`
	MaxLifetime            time.Duration `yaml:"maxLifetime"`
	SlowQueryThreshold     time.Duration `yaml:"slowQueryThreshold"`
	SkipDefaultTransaction *bool         `yaml:"skipDefaultTransaction"`
	PrepareStmt            *bool         `yaml:"prepareStmt"`
}

// RedisConfig -
type RedisConfig struct {
	Password        string               `yaml:"password"`
	Host            string               `yaml:"host"`
	Port            int32                `yaml:"port"`
	Database        uint8                `yaml:"database"`
	IdleTimeout     time.Duration        `yaml:"idleTimeout"`
	MaxConnLifetime time.Duration        `yaml:"maxConnLifetime"`
	MaxIdle         int32                `yaml:"maxIdle"`
	MaxActive       int32                `yaml:"maxActive"`
	TestOnBorrow    bool                 `yaml:"testOnBorrow"`
	Sentinel        *RedisSentinelConfig `yaml:"sentinel"`
	Cluster         *RedisClusterConfig  `yaml:"cluster"`
}

// RedisSentinelConfig -
type RedisSentinelConfig struct {
	// Master - MasterName is a name of Redis master Sentinel servers monitor(Name of the Redis server).
	Master string `yaml:"master"`
	// Addrs - Addrs is a slice of with known Sentinel addresses (A list of "host:port" pairs).
	Addrs []string `yaml:"addrs"`
}

// RedisClusterConfig -
type RedisClusterConfig struct {
	// Addrs - A list of "host:port" pairs to bootstrap from. This represents an "initial" list of cluster nodes and is required to have at least one entry.
	Addrs []string `yaml:"addrs"`
	// MaxRedirects -Maximum number of redirects to follow when executing commands across the cluster.
	MaxRedirects int32 `yaml:"maxRedirects"`
}

// MetricsConfig -
type MetricsConfig struct {
	Host               string `yaml:"host"`
	Port               int32  `yaml:"port"`
	MetricsPath        string `yaml:"metricsPath"`
	BuildInfoCollector *bool  `yaml:"buildInfoCollector"`
	ProcessCollector   *bool  `yaml:"processCollector"`
	GoCollector        *bool  `yaml:"goCollector"`
}

// GraceTerminationConfig -
type GraceTerminationConfig struct {
	GraceTerminationPeriod time.Duration `yaml:"graceTerminationPeriod"`
}

// CronConfig -
type CronConfig struct {
	SkipIfStillRunning *bool            `yaml:"skipIfStillRunning"`
	CronJobs           []*CronJobConfig `yaml:"cronJobs"`
}

// CronJobConfig -
type CronJobConfig struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
}
