package nconf

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// MustNewConfig -
func MustNewConfig(bytes []byte) *Config {
	return MustNewConfigCustom(bytes, nil)
}

// MustNewConfigCustom  -
func MustNewConfigCustom(bytes []byte, customConfig interface{ SetConfig(config *Config) }) *Config {
	config := &Config{}
	if err := yaml.Unmarshal(bytes, config); err != nil {
		log.Fatal(err)
	}
	config.SetDefaultValues()

	// custom config
	if customConfig != nil {
		if err := yaml.Unmarshal(bytes, customConfig); err != nil {
			log.Fatal(err)
		} else {
			customConfig.SetConfig(config)
		}
	}

	return config
}

// MustLoadConfig -
func MustLoadConfig(confPath string) *Config {
	return MustLoadConfigCustom(confPath, nil)
}

// MustLoadConfigCustom  -
func MustLoadConfigCustom(confPath string, customConfig interface{ SetConfig(config *Config) }) *Config {
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

	return MustNewConfigCustom(data, customConfig)
}
