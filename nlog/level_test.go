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
