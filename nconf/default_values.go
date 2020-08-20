package nconf

import (
	"time"

	"nfgo.ga/nfgo/nutil"
)

func setDefaultValues(configs ...interface{ setDefaultValues() error }) error {
	for _, config := range configs {
		if nutil.IsNotNil(config) {
			if err := config.setDefaultValues(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (conf *Config) setDefaultValues() error {
	return setDefaultValues(
		conf.DB,
		conf.Redis,
		conf.Web,
		conf.RPC,
		conf.Security,
		conf.Metrics,
	)
}

func (conf *DbConfig) setDefaultValues() error {
	if conf.MaxOpen == 0 {
		conf.MaxOpen = 2
	}
	if !conf.SkipDefaultTransaction {
		conf.SkipDefaultTransaction = true
	}
	if !conf.PrepareStmt {
		conf.PrepareStmt = true
	}
	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}
	return nil
}

func (conf *RedisConfig) setDefaultValues() error {
	if conf.MaxActive == 0 {
		conf.MaxActive = 5
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 5 * time.Minute
	}
	if conf.MaxConnLifetime == 0 {
		conf.MaxConnLifetime = 30 * time.Minute
	}
	return nil
}

func (conf *WebConfig) setDefaultValues() error {
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	if conf.Port == 0 {
		conf.Port = 8080
	}
	if conf.MaxMultipartMemory == 0 {
		conf.MaxMultipartMemory = 50 << 20 // 50MiB
	}
	return nil
}

func (conf *RPCConfig) setDefaultValues() error {
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	if conf.Port == 0 {
		conf.Port = 9090
	}
	if conf.MaxRecvMsgSize == 0 {
		conf.MaxRecvMsgSize = 50 << 20 // 50MiB
	}
	return nil
}

func (conf *SecurityConfig) setDefaultValues() error {
	if conf.TimeWindow == 0 {
		conf.TimeWindow = 30 * time.Minute
	}
	if conf.SignKeyLifeTime == 0 {
		conf.SignKeyLifeTime = 365 * 24 * time.Hour
	}
	return nil
}

func (conf *MetricsConfig) setDefaultValues() error {
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	if conf.Port == 0 {
		conf.Port = 8079
	}
	if conf.MetricsPath == "" {
		conf.MetricsPath = "/metrics"
	}
	return nil
}
