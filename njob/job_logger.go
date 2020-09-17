package njob

import (
	"strings"
	"time"

	"nfgo.ga/nfgo/nlog"
)

var (
	defaultLogger = &jobLogger{}
)

type jobLogger struct {
}

func (l *jobLogger) Info(msg string, keysAndValues ...interface{}) {
	if nlog.IsLevelEnabled(nlog.InfoLevel) {
		keysAndValues = formatTimes(keysAndValues)
		nlog.Infof(
			formatString(len(keysAndValues)),
			append([]interface{}{msg}, keysAndValues...)...)
	}
}
func (l *jobLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	nlog.Errorf(
		formatString(len(keysAndValues)+2),
		append([]interface{}{msg, "error", err}, keysAndValues...)...)
}

func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}

func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("%v=%v")
	}
	return sb.String()
}
