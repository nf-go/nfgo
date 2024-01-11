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
	"context"

	"github.com/nf-go/nfgo/nconf"
	"github.com/nf-go/nfgo/ncontext"
	"go.uber.org/zap"
)

// NLogger -
type NLogger interface {
	IsLevelEnabled(level Level) bool
	LevelString() string
	WithError(err error) NLogger
	WithField(key string, value interface{}) NLogger
	WithFields(fields Fields) NLogger

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

var (
	logger    *nlogger = newDefaultLogger()
	pkgLogger *nlogger = newPkgLogger(logger)
)

// InitLogger -
func InitLogger(config *nconf.Config) {
	fields := NewFields(
		"group", config.App.Group,
		"app", config.App.Name,
		"profile", config.App.Profile,
	)
	zapConfig := newDefaultZapConfig()
	zapConfig.InitialFields = fields

	logConf := config.Log
	setOutput(zapConfig, config)
	setFormatter(zapConfig, logConf)
	setLevel(zapConfig, logConf)

	zapLogger := mustNewZapLogger(zapConfig)
	logger = &nlogger{zapLogger, zapConfig}
	pkgLogger = newPkgLogger(logger)
}

// Logger -
func Logger(ctx context.Context) NLogger {
	if ctx != nil {
		mdc, _ := ncontext.CurrentMDC(ctx)
		if mdc != nil {
			fields := NewFields(
				"traceID", mdc.TraceID(),
				"subjectID", mdc.SubjectID(),
				"rpcName", mdc.RPCName(),
				"apiName", mdc.APIName(),
				"clientIP", mdc.ClientIP(),
				"clientType", mdc.ClientType(),
			)
			return logger.WithFields(fields)
		}
	}
	return logger
}

type nlogger struct {
	*zap.SugaredLogger
	*zap.Config
}

func newDefaultLogger() *nlogger {
	zapConfig := newDefaultZapConfig()
	return &nlogger{
		mustNewZapLogger(zapConfig),
		zapConfig,
	}
}

func newPkgLogger(logger *nlogger) *nlogger {
	return &nlogger{
		logger.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
		logger.Config,
	}
}

func (l *nlogger) IsLevelEnabled(level Level) bool {
	//nolint:errcheck // ignore errcheck!
	l.Sync()
	return l.Config.Level.Enabled(level.unWrap())
}

func (l *nlogger) WithError(err error) NLogger {
	return &nlogger{l.With("error", err), l.Config}
}

func (l *nlogger) WithField(key string, value interface{}) NLogger {
	return &nlogger{l.With(key, value), l.Config}
}

func (l *nlogger) WithFields(fields Fields) NLogger {
	args := make([]interface{}, 0, len(fields))
	for key, value := range fields {
		args = append(args, key, value)
	}
	return &nlogger{l.With(args...), l.Config}
}

func (l *nlogger) LevelString() string {
	return l.Config.Level.Level().String()
}

// IsLevelEnabled -
func IsLevelEnabled(level Level) bool {
	return pkgLogger.Config.Level.Enabled(level.unWrap())
}

// Debugf -
func Debugf(format string, args ...interface{}) {
	pkgLogger.Debugf(format, args...)
}

// Infof -
func Infof(format string, args ...interface{}) {
	pkgLogger.Infof(format, args...)
}

// Warnf -
func Warnf(format string, args ...interface{}) {
	pkgLogger.Warnf(format, args...)
}

// Errorf -
func Errorf(format string, args ...interface{}) {
	pkgLogger.Errorf(format, args...)
}

// Fatalf -
func Fatalf(format string, args ...interface{}) {
	pkgLogger.Fatalf(format, args...)
}

// Panicf -
func Panicf(format string, args ...interface{}) {
	pkgLogger.Panicf(format, args...)
}

// Debug -
func Debug(args ...interface{}) {
	pkgLogger.Debug(args...)
}

// Info -
func Info(args ...interface{}) {
	pkgLogger.Info(args...)
}

// Warn -
func Warn(args ...interface{}) {
	pkgLogger.Warn(args...)
}

// Error -
func Error(args ...interface{}) {
	pkgLogger.Error(args...)
}

// Fatal -
func Fatal(args ...interface{}) {
	pkgLogger.Fatal(args...)
}

// Panic -
func Panic(args ...interface{}) {
	pkgLogger.Panic(args...)
}

// SetLevel - alters the logging level.
// It lets you safely change the log level of a tree of loggers (the root logger and any children created by adding context) at runtime.
func SetLevel(level Level) {
	pkgLogger.Config.Level.SetLevel(level.unWrap())
}

// Sync - Sync calls the underlying Core's Sync method, flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func Sync() error {
	// ignore "sync /dev/stderr: inappropriate ioctl for device" error
	//nolint:errcheck // ignore errcheck!
	logger.SugaredLogger.Sync()
	//nolint:errcheck // ignore errcheck!
	pkgLogger.SugaredLogger.Sync()
	return nil
}
