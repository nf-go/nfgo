package nconf

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	// defaultTimeWindowMinite -
	defaultTimeWindowMinite time.Duration = 30 * time.Minute
	// defaultSignKeyLifeTime -
	defaultSignKeyLifeTime time.Duration = 365 * 24 * time.Hour
)

// MustLoadConfig -
func MustLoadConfig(confPath string) *Config {
	if confPath == "" {
		confPath = "app.yaml"
	}

	file, err := os.Open(confPath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{}
	err = yaml.Unmarshal([]byte(data), config)
	if err != nil {
		log.Fatal(err)
	}

	securityConfig := config.Security
	if securityConfig != nil {

		if securityConfig.TimeWindow == 0 {
			securityConfig.TimeWindow = defaultTimeWindowMinite
		}
		if securityConfig.SignKeyLifeTime == 0 {
			securityConfig.SignKeyLifeTime = defaultSignKeyLifeTime
		}
	}

	return config
}
