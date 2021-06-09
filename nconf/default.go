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

	"nfgo.ga/nfgo/nutil/ntypes"
)

// SetDefaultValues -
func (conf *Config) SetDefaultValues() {
	if conf.Log == nil {
		conf.Log = &LogConfig{}
	}
	configs := []interface{ SetDefaultValues() }{
		conf.App,
		conf.DB,
		conf.Redis,
		conf.Web,
		conf.RPC,
		conf.CronConfig,
		conf.Metrics,
	}
	for _, c := range configs {
		if ntypes.IsNotNil(c) {
			c.SetDefaultValues()
		}
	}
}

// SetDefaultValues -
func (conf *AppConfig) SetDefaultValues() {
	if conf.GraceTermination == nil {
		conf.GraceTermination = &GraceTerminationConfig{}
	}
	conf.GraceTermination.SetDefaultValues()
}

// SetDefaultValues -
func (conf *LogConfig) SetDefaultValues() {
	conf.Level = "info"
	conf.Format = "json"
}

// SetDefaultValues -
func (conf *DbConfig) SetDefaultValues() {
	if conf.MaxOpen == 0 {
		conf.MaxOpen = 2
	}
	if conf.SkipDefaultTransaction == nil {
		conf.SkipDefaultTransaction = ntypes.Bool(true)
	}
	if conf.PrepareStmt == nil {
		conf.PrepareStmt = ntypes.Bool(true)
	}
	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}
}

// SetDefaultValues -
func (conf *RedisConfig) SetDefaultValues() {
	if conf.MaxActive == 0 {
		conf.MaxActive = 5
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 5 * time.Minute
	}
	if conf.MaxConnLifetime == 0 {
		conf.MaxConnLifetime = 30 * time.Minute
	}
}

// SetDefaultValues -
func (conf *WebConfig) SetDefaultValues() {
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	if conf.Port == 0 {
		conf.Port = 8080
	}
	if conf.MaxMultipartMemory == 0 {
		conf.MaxMultipartMemory = 50 << 20 // 50MiB
	}
}

// SetDefaultValues -
func (conf *RPCConfig) SetDefaultValues() {
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	if conf.Port == 0 {
		conf.Port = 9090
	}
	if conf.MaxRecvMsgSize == 0 {
		conf.MaxRecvMsgSize = 50 << 20 // 50MiB
	}
	if conf.RegisterHealthServer == nil {
		conf.RegisterHealthServer = ntypes.Bool(true)
	}
	if conf.RegisterReflectionServer == nil {
		conf.RegisterReflectionServer = ntypes.Bool(true)
	}
	for _, clientConf := range conf.Clients {
		if clientConf != nil {
			clientConf.SetDefaultValues()
		}
	}
}

// SetDefaultValues -
func (conf *RPCClientConfig) SetDefaultValues() {
	if conf.MaxCallSendMsgSize == 0 {
		conf.MaxCallSendMsgSize = 50 << 20 // 50MiB
	}
	if conf.MaxCallRecvMsgSize == 0 {
		conf.MaxCallRecvMsgSize = 50 << 20 // 50MiB
	}
	if conf.Plaintext == nil {
		conf.Plaintext = ntypes.Bool(true)
	}
}

// SetDefaultValues -
func (conf *MetricsConfig) SetDefaultValues() {
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	if conf.Port == 0 {
		conf.Port = 8079
	}
	if conf.MetricsPath == "" {
		conf.MetricsPath = "/metrics"
	}
	if conf.BuildInfoCollector == nil {
		conf.BuildInfoCollector = ntypes.Bool(false)
	}
	if conf.ProcessCollector == nil {
		conf.ProcessCollector = ntypes.Bool(false)
	}
	if conf.GoCollector == nil {
		conf.GoCollector = ntypes.Bool(false)
	}
}

// SetDefaultValues -
func (conf *GraceTerminationConfig) SetDefaultValues() {
	if conf.GraceTerminationPeriod == 0 {
		conf.GraceTerminationPeriod = 10 * time.Second
	}
}

// SetDefaultValues -
func (conf *CronConfig) SetDefaultValues() {
	if conf.SkipIfStillRunning == nil {
		conf.SkipIfStillRunning = ntypes.Bool(true)
	}
}
