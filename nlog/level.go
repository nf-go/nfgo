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
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	// DebugLevel level.
	DebugLevel = Level(zapcore.DebugLevel)
	// InfoLevel level.
	InfoLevel = Level(zapcore.InfoLevel)
	// WarnLevel level.
	WarnLevel = Level(zapcore.WarnLevel)
	// ErrorLevel level.
	ErrorLevel = Level(zapcore.ErrorLevel)
	// PanicLevel level.
	PanicLevel Level = Level(zapcore.PanicLevel)
	// FatalLevel level.
	FatalLevel = Level(zapcore.FatalLevel)
)

// Level -
type Level zapcore.Level

func (l Level) String() string {
	var levelStr string
	switch l.unWrap() {
	case zapcore.DebugLevel:
		levelStr = "debug"
	case zapcore.InfoLevel:
		levelStr = "info"
	case zapcore.WarnLevel:
		levelStr = "warn"
	case zapcore.ErrorLevel:
		levelStr = "error"
	case zapcore.FatalLevel:
		levelStr = "fatal"
	case zapcore.PanicLevel:
		levelStr = "panic"
	}
	return levelStr
}

func parseLevel(levelStr string) Level {
	level := InfoLevel
	switch strings.ToLower(levelStr) {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "panic":
		level = PanicLevel
	case "fatal":
		level = FatalLevel
	}
	return level
}

func (l Level) unWrap() zapcore.Level {
	return zapcore.Level(l)
}
