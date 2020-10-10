package nlog

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/ncontext"
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
	logger = &nlogger{zapLogger, zapConfig.Level.Level()}
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
	zapcore.Level
}

func newDefaultLogger() *nlogger {
	return &nlogger{
		mustNewZapLogger(newDefaultZapConfig()),
		zapcore.InfoLevel,
	}
}

func newPkgLogger(logger *nlogger) *nlogger {
	return &nlogger{
		logger.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
		logger.Level,
	}
}

func (l *nlogger) IsLevelEnabled(level Level) bool {
	return l.Level.Enabled(level.unWrap())
}

func (l *nlogger) WithError(err error) NLogger {
	return &nlogger{l.With("error", err), l.Level}
}

func (l *nlogger) WithField(key string, value interface{}) NLogger {
	return &nlogger{l.With(key, value), l.Level}
}

func (l *nlogger) WithFields(fields Fields) NLogger {
	args := make([]interface{}, 0, len(fields))
	for key, value := range fields {
		args = append(args, key, value)
	}
	return &nlogger{l.With(args...), l.Level}
}

func (l *nlogger) LevelString() string {
	return Level(l.Level).String()
}

// IsLevelEnabled -
func IsLevelEnabled(level Level) bool {
	return pkgLogger.Level.Enabled(level.unWrap())
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
