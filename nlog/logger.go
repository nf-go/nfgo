package nlog

import (
	"context"

	"github.com/sirupsen/logrus"
	"nfgo.ga/nfgo/ncontext"
)

// Logger -
func Logger(ctx context.Context) *logrus.Entry {
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
			return entry.WithFields(fields)
		}
	}
	return entry
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
