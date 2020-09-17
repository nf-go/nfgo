package ndb

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
)

// https://github.com/go-gorm/gorm/blob/master/logger/logger.go
type dbLogger struct {
	LogLevel                            logger.LogLevel
	SlowThreshold                       time.Duration
	traceStr, traceErrStr, traceWarnStr string
}

func newLogger(config *nconf.DbConfig) *dbLogger {
	nlogLevel := nlog.Logger(context.Background()).LevelString()
	logLevel := logger.Silent
	switch nlogLevel {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Warn
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	}

	slowThreshold := config.SlowQueryThreshold
	if slowThreshold == 0 {
		slowThreshold = 500 * time.Millisecond
	}

	return &dbLogger{
		SlowThreshold: slowThreshold,
		LogLevel:      logLevel,
	}
}

func (l *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		nlog.Logger(ctx).Infof(msg, data)
	}

}

func (l *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		nlog.Logger(ctx).Warnf(msg, data)
	}
}

func (l *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		nlog.Logger(ctx).Errorf(msg, data)
	}
}

func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			fileWithLineNum := utils.FileWithLineNum()
			logEnry(ctx, fileWithLineNum, elapsed, fc).WithError(err).Error()
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			fileWithLineNum := utils.FileWithLineNum()
			logEnry(ctx, fileWithLineNum, elapsed, fc).Warn()
		case l.LogLevel >= logger.Info:
			fileWithLineNum := utils.FileWithLineNum()
			logEnry(ctx, fileWithLineNum, elapsed, fc).Debug()
		}
	}
}

func logEnry(ctx context.Context, fileWithLineNum string, elapsed time.Duration, fc func() (string, int64)) nlog.NLogger {
	sql, rows := fc()
	return nlog.Logger(ctx).WithFields(nlog.Fields{
		"fileWithLineNum": fileWithLineNum,
		"rowsAffected":    rows,
		"elapsed":         elapsed.Milliseconds(),
		"sql":             sql,
	})
}
