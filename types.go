package log

import (
	"context"

	"github.com/netbookai/log/loggers"
)

// Logger interface is implemnted by the log implementation
type Logger interface {
	loggers.BaseLogger
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
}
