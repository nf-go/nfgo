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
	"io"
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

	// secret config value decrypt
	if secretsConfig, ok := customConfig.(SecretsConfig); ok {
		secretKey = secretsConfig.SecretKey()
		if err := config.decryptSecrets(); err != nil {
			log.Fatal(err)
		}
		if err := secretsConfig.DecryptSecrets(); err != nil {
			log.Fatal()
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
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return MustNewConfigCustom(data, customConfig)
}
