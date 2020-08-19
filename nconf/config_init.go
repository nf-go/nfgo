package nconf

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

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

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Fatal(err)
	}
	setDefaultValues(
		config.Security,
		config.Metrics,
	)

	// custom config
	if customConfig != nil {
		if err = yaml.Unmarshal(data, customConfig); err != nil {
			log.Fatal(err)
		} else {
			customConfig.SetConfig(config)
		}
	}

	return config
}

func setDefaultValues(configs ...interface{ setDefaultValues() error }) error {
	for _, config := range configs {
		if config != nil {
			if err := config.setDefaultValues(); err != nil {
				return err
			}
		}
	}
	return nil
}
