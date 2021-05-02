//Package zap provides a BaseLogger implementation for uber/zap
package zap

import (
	"context"

	"github.com/go-coldbrew/log/loggers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logger *zap.SugaredLogger
	opt    loggers.Options
	cfg    zap.Config
}

func (l *logger) Log(ctx context.Context, level loggers.Level, skip int, args ...interface{}) {

	logger := l.logger
	var msg interface{}
	//if there are odd number of elements in args, first will be treated as a message and rest will
	//be key value pair to log in json format
	if len(args)%2 != 0 {
		msg = args[0]
		args = args[1:]
	}
	logger = logger.With(args...)

	ctxFields := loggers.FromContext(ctx)
	if ctxFields != nil {
		for k, v := range ctxFields {
			logger = logger.With(k, v)
		}
	}

	logFunc := l.logger.Error
	switch level {
	case loggers.DebugLevel:
		logFunc = logger.Debug
	case loggers.InfoLevel:
		logFunc = logger.Info
	case loggers.WarnLevel:
		logFunc = logger.Warn
	case loggers.ErrorLevel:
		logFunc = logger.Error
	}
	logFunc(msg)
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

	l = l.WithOptions(zap.AddCallerSkip(loggers.COLBREW_CALL_STACK_SIZE))
	if err != nil {
		//should we fail? will use sugared log here
		l, _ = zap.NewProduction()

	}

	return &logger{
		logger: l.Sugar(),
		opt:    opt,
		cfg:    zapCfg,
	}

}
