package nconf

import (
	"strings"

	"nfgo.ga/nfgo/nerrors"
	"nfgo.ga/nfgo/nutil/ncrypto"
	"nfgo.ga/nfgo/nutil/ntypes"
)

var secretKey string

const (
	encEncryptedTextPrefix = "SECRET("
	encEncryptedTextSuffix = ")"
)

type SecretsConfig interface {
	SecretKey() string
	DecryptSecrets() error
}

func DecryptSecretValue(value string) (string, error) {
	if secretKey == "" {
		return "", nerrors.New("secretKey is not set")
	}
	if strings.HasPrefix(value, encEncryptedTextPrefix) && strings.HasSuffix(encEncryptedTextSuffix, ")") {
		value = value[len(encEncryptedTextPrefix) : len(value)-len(encEncryptedTextSuffix)]
		return ncrypto.AESDecryptString(value, secretKey)
	}
	return value, nil
}

func (conf *Config) decryptSecrets() error {
	configs := []interface{ decryptSecrets() error }{
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
			if err := c.decryptSecrets(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (conf *AppConfig) decryptSecrets() error {
	return nil
}

func (conf *LogConfig) decryptSecrets() error {
	return nil
}

func (conf *DbConfig) decryptSecrets() error {
	password, err := DecryptSecretValue(conf.Password)
	if err != nil {
		return err
	}
	conf.Password = password
	return nil
}

func (conf *RedisConfig) decryptSecrets() error {
	password, err := DecryptSecretValue(conf.Password)
	if err != nil {
		return err
	}
	conf.Password = password
	return nil
}

func (conf *WebConfig) decryptSecrets() error {
	return nil
}

func (conf *RPCConfig) decryptSecrets() error {
	return nil
}

func (conf *CronConfig) decryptSecrets() error {
	return nil
}

func (conf *MetricsConfig) decryptSecrets() error {
	return nil
}
