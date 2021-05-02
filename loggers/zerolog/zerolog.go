package zerolog

import (
	"context"
	"os"

	"github.com/go-coldbrew/log/loggers"
	"github.com/rs/zerolog"
)

type logger struct {
	logger zerolog.Logger
	level  loggers.Level
}

func (l *logger) GetLevel() loggers.Level {
	return l.level
}

func (l *logger) SetLevel(level loggers.Level) {
	l.logger = l.logger.Level(toZerologLevel(level))
	l.level = level
}

func (l *logger) Log(ctx context.Context, level loggers.Level, skip int, args ...interface{}) {

	logger := l.logger

	var msg string
	//if there are odd number of elements in args, first will be treated as a message and rest will
	//be key value pair to log in json format
	if len(args)%2 != 0 {
		msg = args[0].(string)
		args = args[1:]
	}

	logFunc := logger.Error
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

	var logEvent *zerolog.Event = logFunc()
	ctxFields := loggers.FromContext(ctx)
	if ctxFields != nil {
		logEvent = logFunc().Fields(ctxFields)
	}

	//log each message field as key value pair from the args
	for i := 0; i < len(args); i += 2 {
		logEvent = logEvent.Interface(args[i].(string), args[i+1])
	}

	logEvent.Msg(msg)
}

func toZerologLevel(level loggers.Level) zerolog.Level {
	switch level {
	case loggers.DebugLevel:
		return zerolog.DebugLevel
	case loggers.InfoLevel:
		return zerolog.InfoLevel
	case loggers.WarnLevel:
		return zerolog.WarnLevel
	case loggers.ErrorLevel:
		return zerolog.ErrorLevel
	default:
		return zerolog.ErrorLevel
	}
}

func NewLogger(options ...loggers.Option) loggers.BaseLogger {

	opt := loggers.GetDefaultOptions()
	// read options
	for _, f := range options {
		f(&opt)
	}
	log := zerolog.New(os.Stdout).
		With().
		CallerWithSkipFrameCount(loggers.COLDBREW_CALL_STACK_SIZE + 2).
		Logger().
		Level(toZerologLevel(opt.Level))

	return &logger{
		logger: log,
		level:  opt.Level,
	}
}
