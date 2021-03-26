//Package zap provides a BaseLogger implementation for uber/zap
package zap

import (
	"context"
	"fmt"

	"github.com/go-coldbrew/log/loggers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logger *zap.Logger
	opt    loggers.Options
	cfg    zap.Config
}

func (l *logger) Log(ctx context.Context, level loggers.Level, skip int, args ...interface{}) {

	msg := fmt.Sprint(args...)
	logger := l.logger
	switch level {
	case loggers.DebugLevel:
		logger.Debug(msg)
	case loggers.InfoLevel:
		logger.Info(msg)
	case loggers.WarnLevel:
		logger.Warn(msg)
	case loggers.ErrorLevel:
		logger.Error(msg)
	default:
		l.logger.Error(msg)
	}
}

func (l *logger) GetLevel() loggers.Level {
	return l.opt.Level
}

func (l *logger) SetLevel(level loggers.Level) {
	l.opt.Level = level
	l.cfg.Level.SetLevel(toZapLevel(level))
}

func toZapLevel(level loggers.Level) zapcore.Level {

	switch level {
	case loggers.DebugLevel:
		return zapcore.DebugLevel
	case loggers.InfoLevel:
		return zap.InfoLevel
	case loggers.WarnLevel:
		return zap.WarnLevel
	case loggers.ErrorLevel:
		return zap.ErrorLevel
	default:
		return zapcore.ErrorLevel
	}
}

func NewLogger(options ...loggers.Option) loggers.BaseLogger {

	opt := loggers.GetDefaultOptions()
	// read options
	for _, f := range options {
		f(&opt)
	}

	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(toZapLevel(opt.Level)),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    opt.LevelFieldName,
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    opt.TimestampFieldName,
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    opt.CallerFieldName,
			EncodeCaller: zapcore.FullCallerEncoder,
		},
	}

	if opt.JSONLogs {
		zapCfg.Encoding = "json"
	} else {
		zapCfg.Encoding = "console"
	}
	l, err := zapCfg.Build()

	if err != nil {
		//should we fail? will use sugared log here
		l, _ = zap.NewProduction()

	}
	return &logger{
		logger: l,
		opt:    opt,
		cfg:    zapCfg,
	}

}
