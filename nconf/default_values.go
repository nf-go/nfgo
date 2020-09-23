package nconf

import (
	"time"

	"nfgo.ga/nfgo/nutil"
	"nfgo.ga/nfgo/nutil/ntypes"
)

func setDefaultValues(configs ...interface{ SetDefaultValues() }) {
	for _, config := range configs {
		if nutil.IsNotNil(config) {
			config.SetDefaultValues()
		}
	}
}

// SetDefaultValues -
func (conf *Config) SetDefaultValues() {
	if conf.GraceTermination == nil {
		conf.GraceTermination = &GraceTerminationConfig{}
	}
	setDefaultValues(
		conf.DB,
		conf.Redis,
		conf.Web,
		conf.RPC,
		conf.CronConfig,
		conf.Security,
		conf.Metrics,
		conf.GraceTermination,
	)
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
func (conf *SecurityConfig) SetDefaultValues() {
	if conf.TimeWindow == 0 {
		conf.TimeWindow = 30 * time.Minute
	}
	if conf.SignKeyLifeTime == 0 {
		conf.SignKeyLifeTime = 365 * 24 * time.Hour
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
