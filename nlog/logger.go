package nlog

import (
	"context"

	"github.com/sirupsen/logrus"
	"nfgo.ga/nfgo/ncontext"
)

const (
	// PanicLevel level, highest level of severity.
	PanicLevel Level = Level(logrus.PanicLevel)
	// FatalLevel level.
	FatalLevel = Level(logrus.FatalLevel)
	// ErrorLevel level.
	ErrorLevel = Level(logrus.ErrorLevel)
	// WarnLevel level.
	WarnLevel = Level(logrus.WarnLevel)
	// InfoLevel level.
	InfoLevel = Level(logrus.InfoLevel)
	// DebugLevel level.
	DebugLevel = Level(logrus.DebugLevel)
	// TraceLevel level.
	TraceLevel = Level(logrus.TraceLevel)
)

// Level -
type Level logrus.Level

// Fields -
type Fields logrus.Fields

// NLogger -
type NLogger interface {
	IsLevelEnabled(level Level) bool
	LevelString() string
	WithError(err error) NLogger
	WithField(key string, value interface{}) NLogger
	WithFields(fields Fields) NLogger

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

type nlogger struct {
	*logrus.Entry
}

func (l *nlogger) IsLevelEnabled(level Level) bool {
	return l.Entry.Logger.IsLevelEnabled(logrus.Level(level))
}

func (l *nlogger) WithError(err error) NLogger {
	return &nlogger{l.Entry.WithError(err)}
}

func (l *nlogger) WithField(key string, value interface{}) NLogger {
	return &nlogger{l.Entry.WithField(key, value)}
}

func (l *nlogger) WithFields(fields Fields) NLogger {
	return &nlogger{l.Entry.WithFields(logrus.Fields(fields))}
}

func (l *nlogger) LevelString() string {
	return l.Entry.Logger.Level.String()
}

// Logger -
func Logger(ctx context.Context) NLogger {
	if ctx != nil {
		mdc, _ := ncontext.CurrentMDC(ctx)
		if mdc != nil {
			var fields = logrus.Fields{}
			putNotEmptyVal(fields, "traceID", mdc.TraceID())
			putNotEmptyVal(fields, "subjectID", mdc.SubjectID())
			putNotEmptyVal(fields, "rpcName", mdc.RPCName())
			putNotEmptyVal(fields, "apiName", mdc.APIName())
			putNotEmptyVal(fields, "clientIP", mdc.ClientIP())
			putNotEmptyVal(fields, "clientType", mdc.ClientType())
			return &nlogger{entry.WithFields(fields)}
		}
	}
	return &nlogger{entry}
}

// IsLevelEnabled -
func IsLevelEnabled(level Level) bool {
	return entry.Logger.IsLevelEnabled(logrus.Level(level))
}

// Tracef -
func Tracef(format string, args ...interface{}) {
	entry.Logf(logrus.TraceLevel, format, args...)
}

// Debugf -
func Debugf(format string, args ...interface{}) {
	entry.Logf(logrus.DebugLevel, format, args...)
}

// Infof -
func Infof(format string, args ...interface{}) {
	entry.Logf(logrus.InfoLevel, format, args...)
}

// Warnf -
func Warnf(format string, args ...interface{}) {
	entry.Logf(logrus.WarnLevel, format, args...)
}

// Errorf -
func Errorf(format string, args ...interface{}) {
	entry.Logf(logrus.ErrorLevel, format, args...)
}

// Fatalf -
func Fatalf(format string, args ...interface{}) {
	entry.Logf(logrus.FatalLevel, format, args...)
	entry.Logger.Exit(1)
}

// Panicf -
func Panicf(format string, args ...interface{}) {
	entry.Logf(logrus.PanicLevel, format, args...)
}

// Trace -
func Trace(args ...interface{}) {
	entry.Log(logrus.TraceLevel, args...)
}

// Debug -
func Debug(args ...interface{}) {
	entry.Log(logrus.DebugLevel, args...)
}

// Info -
func Info(args ...interface{}) {
	entry.Log(logrus.InfoLevel, args...)
}

// Warn -
func Warn(args ...interface{}) {
	entry.Log(logrus.WarnLevel, args...)
}

// Error -
func Error(args ...interface{}) {
	entry.Log(logrus.ErrorLevel, args...)
}

// Fatal -
func Fatal(args ...interface{}) {
	entry.Log(logrus.FatalLevel, args...)
	entry.Logger.Exit(1)
}

// Panic -
func Panic(args ...interface{}) {
	entry.Log(logrus.PanicLevel, args...)
}
