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

package nlog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nf-go/nfgo/nconf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newDefaultZapConfig() *zap.Config {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return &zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func mustNewZapLogger(zapConfig *zap.Config) *zap.SugaredLogger {
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatal("fail to new default zap logger: ", err)
	}
	sugar := logger.Sugar()
	return sugar
}

func setLevel(zapConfig *zap.Config, logConf *nconf.LogConfig) {
	level := parseLevel(logConf.Level)
	zapConfig.Level = zap.NewAtomicLevelAt(level.unWrap())
}

func setFormatter(zapConfig *zap.Config, logConf *nconf.LogConfig) {
	var layout string
	if logConf.TimestampFormat == "" {
		layout = "2006-01-02T15:04:05.000Z07:00"
	} else {
		layout = logConf.TimestampFormat
	}
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		type appendTimeEncoder interface {
			AppendTimeLayout(time.Time, string)
		}
		if enc, ok := enc.(appendTimeEncoder); ok {
			enc.AppendTimeLayout(t, layout)
			return
		}
		enc.AppendString(t.Format(layout))
	}
	switch logConf.Format {
	case "text":
		zapConfig.Encoding = "console"
	case "json":
		zapConfig.Encoding = "json"
	}

	if logConf.CallerPrint {
		zapConfig.EncoderConfig.CallerKey = "caller"
	}
}

func setOutput(zapConfig *zap.Config, config *nconf.Config) {
	logConf := config.Log
	if logConf.LogPath == "" {
		return
	}
	logFilename := logConf.LogFilename
	if logFilename == "" {
		hostname, _ := os.Hostname()
		logFilename = fmt.Sprintf("%s.%s.%s.log", config.App.Name, config.App.Profile, hostname)
	}
	logPath := filepath.Join(logConf.LogPath, time.Now().Format("200601"))
	if err := os.MkdirAll(logPath, 0755); err != nil {
		log.Fatal("can't create log dir: ", err)
	}
	fullFilename := filepath.Join(logPath, logFilename)

	zapConfig.OutputPaths = []string{fullFilename}
}
