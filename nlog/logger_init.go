package nlog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"nfgo.ga/nfgo/nconf"

	"github.com/sirupsen/logrus"
)

var (
	entry *logrus.Entry
)

// InitLogger -
func InitLogger(config *nconf.Config) {
	fields := logrus.Fields{}
	putNotEmptyVal(fields, "app", config.App.Name)
	putNotEmptyVal(fields, "profile", config.App.Profile)
	entry = logrus.WithFields(fields)

	logConf := config.Log
	if logConf == nil {
		logConf = &nconf.LogConfig{
			Level:  "info",
			Format: "json",
		}
	}

	// set output
	setOutput(entry, config)

	// set formater
	setFormatter(entry, logConf)

	// set level
	setLevel(entry, logConf)

}

func putNotEmptyVal(fields logrus.Fields, key, val string) {
	if val != "" {
		fields[key] = val
	}
}

func setOutput(logEntry *logrus.Entry, config *nconf.Config) {
	logConf := config.Log
	if logConf.LogPath == "" {
		return
	}
	logFilename := logConf.LogFilename
	if logFilename == "" {
		hostname, _ := os.Hostname()
		logFilename = fmt.Sprintf("%s.%s.%s.log", config.App.Name, config.App.Profile, hostname)
	}
	fullLocation := filepath.Join(logConf.LogPath, logFilename)
	logFile, err := os.OpenFile(fullLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("can't open log file: ", err)
	}
	entry.Logger.SetOutput(logFile)
	logrus.RegisterExitHandler(func() {
		if logFile != nil {
			logFile.Close()
		}
	})
}

func setLevel(logEntry *logrus.Entry, logConf *nconf.LogConfig) {
	var level logrus.Level
	switch strings.ToLower(logConf.Level) {
	case "debug":
		level = logrus.DebugLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "info":
		level = logrus.InfoLevel
	}
	logEntry.Logger.Level = level
}

func setFormatter(logEntry *logrus.Entry, logConf *nconf.LogConfig) {
	var timestampFormat string
	if logConf.TimestampFormat == "" {
		timestampFormat = "2006-02-01 15:04:05.000"
	} else {
		timestampFormat = logConf.TimestampFormat
	}

	var formater logrus.Formatter
	switch logConf.Format {
	case "text":
		formater = &logrus.TextFormatter{
			FullTimestamp:   true,
			DisableColors:   true,
			TimestampFormat: timestampFormat,
		}
	case "json":
		formater = &logrus.JSONFormatter{
			PrettyPrint:     logConf.PrettyPrint,
			TimestampFormat: timestampFormat,
		}
	}
	logEntry.Logger.SetFormatter(formater)
}
