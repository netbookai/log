package zerolog

import (
	"context"

	"github.com/go-coldbrew/log/loggers"
)

type logger struct {
}

func (l *logger) GetLevel() loggers.Level {
	return loggers.DebugLevel
}

func (l *logger) SetLevel(loggers.Level) {
}

func (l *logger) Log(ctx context.Context, level loggers.Level, skip int, args ...interface{}) {

}

func NewLogger(options ...loggers.Option) loggers.BaseLogger {
	return &logger{}
}
