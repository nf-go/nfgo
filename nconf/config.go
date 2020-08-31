package nconf

import (
	"time"

	"nfgo.ga/nfgo/nconst"
)

// Config data
type Config struct {
	App      *AppConfig      `yaml:"app"`
	Log      *LogConfig      `yaml:"log"`
	DB       *DbConfig       `yaml:"db"`
	Redis    *RedisConfig    `yaml:"redis"`
	Web      *WebConfig      `yaml:"web"`
	RPC      *RPCConfig      `yaml:"rpc"`
	Security *SecurityConfig `yaml:"security"`
	Metrics  *MetricsConfig  `yaml:"metrics"`
}

// AppConfig -
type AppConfig struct {
	Name    string    `yaml:"name"`
	Profile string    `yaml:"profile"`
	Ext     ExtConfig `yaml:"ext"`
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
	Host                 string                      `yaml:"host"`
	Port                 int32                       `yaml:"port"`
	MaxRecvMsgSize       int64                       `yaml:"maxRecvMsgSize"`
	RegisterHealthServer *bool                       `yaml:"registerHealthServer"`
	Clients              map[string]*RPCClientConfig `yaml:"clients"`
}

// RPCClientConfig -
type RPCClientConfig struct {
	Addr               string `yaml:"addr"`
	Plaintext          *bool  `yaml:"plaintext"`
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
	MaxOpen                int32         `yaml:"maxOpen"`
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

// SecurityConfig -
type SecurityConfig struct {
	SignKeyLifeTime    time.Duration `yaml:"signKeyLifeTime"`
	RefreshSignKeyLife bool          `yaml:"refreshSignKeyLife"`
	TimeWindow         time.Duration `yaml:"timeWindow"`
	Anons              []string      `yaml:"anons"`
	Model              string        `yaml:"model"`
	Policies           []string      `yaml:"policies"`
}

// MetricsConfig -
type MetricsConfig struct {
	Host        string `yaml:"host"`
	Port        int32  `yaml:"port"`
	MetricsPath string `yaml:"metricsPath"`
}
