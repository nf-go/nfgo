package nlog

import (
	"context"
	"testing"

	"nfgo.ga/nfgo/nconf"
)

func TestInitLogger(t *testing.T) {
	config := &nconf.Config{
		App: &nconf.AppConfig{
			Name:    "foo-app",
			Profile: "dev",
		},
	}
	config.SetDefaultValues()
	config.Log.CallerPrint = true
	InitLogger(config)
	Info("info log")
	Error("error log")
	logger := Logger(context.Background()).WithField("key1", "val1")
	logger.Info("hello")
	logger.Warn("world")
}
