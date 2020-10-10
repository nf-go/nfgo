package nlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLevelString(t *testing.T) {
	assert.Equal(t, "debug", DebugLevel.String())
	assert.Equal(t, "info", InfoLevel.String())
	assert.Equal(t, "warn", WarnLevel.String())
	assert.Equal(t, "error", ErrorLevel.String())
	assert.Equal(t, "panic", PanicLevel.String())
	assert.Equal(t, "fatal", FatalLevel.String())
}

func TestLevelUnWrap(t *testing.T) {
	assert.Equal(t, zapcore.DebugLevel, DebugLevel.unWrap())
	assert.Equal(t, zapcore.InfoLevel, InfoLevel.unWrap())
	assert.Equal(t, zapcore.WarnLevel, WarnLevel.unWrap())
	assert.Equal(t, zapcore.ErrorLevel, ErrorLevel.unWrap())
	assert.Equal(t, zapcore.PanicLevel, PanicLevel.unWrap())
	assert.Equal(t, zapcore.FatalLevel, FatalLevel.unWrap())
}

func TestLevelParseLevel(t *testing.T) {
	assert.Equal(t, DebugLevel, parseLevel("debug"))
	assert.Equal(t, InfoLevel, parseLevel("info"))
	assert.Equal(t, WarnLevel, parseLevel("warn"))
	assert.Equal(t, ErrorLevel, parseLevel("error"))
	assert.Equal(t, PanicLevel, parseLevel("panic"))
	assert.Equal(t, FatalLevel, parseLevel("fatal"))
}
