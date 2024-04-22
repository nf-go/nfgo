// Copyright 2021 The nfgo Authors. All Rights Reserved.
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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type FooConfig struct {
	*Config
	Foo string `yaml:"foo"`
	Bar string `yaml:"bar"`
}

func (c *FooConfig) SetConfig(config *Config) {
	c.Config = config
}

func TestMustNewConfigCustom(t *testing.T) {

	f := `
---
app:
  name: foo
  profile: local
  ext:
    someIds:
    - 1
    - 2
    strV: str1
    intV: 1
    boolV: true
log:
  format: text
  level: info
  callerPrint: true
rpc:
  port: 9090
web:
  port: 8080
  swagger:
    enabled: true
    url: http://localhost:8080/apidoc/doc.json
  sensitiveURLPaths:
    "/login": {}
    "/api/secret": {}
metrics:
  port: 8079
cron:
  cronJobs:
  - name: demoJob
    schedule: "* * * * *"
db:
  host: 127.0.0.1
  port: 3306
  database: test
  username: root
  password: ""
  maxIdle: 0
  maxOpen: 2
  skipDefaultTransaction: true
  prepareStmt: true
redis:
  host: 127.0.0.1
  port: 6379
  database: 1
  password: ""
  maxIdle: 1
  maxActive: 5
  testOnBorrow: true
  idleTimeout: 5m
  maxConnLifetime : 30m
foo: f1
bar: b1
`

	fooConfig := &FooConfig{}
	config := MustNewConfigCustom([]byte(f), fooConfig)
	a := assert.New(t)
	a.Equal(config, fooConfig.Config)
	a.Equal(10*time.Second, config.App.GraceTermination.GraceTerminationPeriod)
	a.Equal("f1", fooConfig.Foo)
	a.Equal("b1", fooConfig.Bar)
	ext := fooConfig.App.Ext
	a.Equal("str1", ext.StrVal("strV"))
	a.Equal(1, ext.IntVal("intV"))
	a.True(ext.BoolVal("boolV"))
	someIds, ok := ext["someIds"].([]interface{})
	a.True(ok)
	a.Equal([]interface{}{1, 2}, someIds)
	_, ok = config.Web.SensitiveURLPaths["/login"]
	a.True(ok)
}
