package xchain

import (
	"os"
	"strings"

	"github.com/xuperchain/xupercore/lib/logs"
)

type mockLogger struct {
	// log15.Logger
}

func (*mockLogger) GetLogId() string {
	return ""
}

func (*mockLogger) SetCommField(key string, value interface{}) {

}
func (*mockLogger) SetInfoField(key string, value interface{}) {

}

func (l *mockLogger) Debug(msg string, ctx ...interface{}) {

}
func (l *mockLogger) Trace(msg string, ctx ...interface{}) {

}
func (l *mockLogger) Info(msg string, ctx ...interface{}) {
	if len(msg) == 0 {
		return
	}
	os.Stdout.WriteString(msg)
	//  27 is escape character
	if msg[0] != byte(27) && !strings.HasSuffix(msg, "\n") {
		os.Stdout.WriteString("\n")
	}
}
func (l *mockLogger) Warn(msg string, ctx ...interface{}) {

}
func (l *mockLogger) Error(msg string, ctx ...interface{}) {

}
func (l *mockLogger) Crit(msg string, ctx ...interface{}) {

}

func NewLogger() logs.Logger {
	return &mockLogger{}
}
