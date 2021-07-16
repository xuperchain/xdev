package xchain

import (
	log15 "github.com/xuperchain/log15"
	"github.com/xuperchain/xupercore/lib/logs"
	"os"
)

type mockLogger struct {
	log15.Logger
}

func (*mockLogger) GetLogId() string {
	return ""
}

func (*mockLogger) SetCommField(key string, value interface{}) {

}
func (*mockLogger) SetInfoField(key string, value interface{}) {

}

func NewLogger() logs.Logger {
	logger := log15.New()
	logger.SetHandler(log15.StreamHandler(os.Stderr, log15.LogfmtFormat()))
	return &mockLogger{logger}
}
