// Package logrus provides a BaseLogger implementation for logrus
package logrus

import (
	"context"
	"fmt"
	stdlog "log"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/netbookai/log/loggers"
)

type logger struct {
	logger *log.Logger
	opt    loggers.Options
}

func toLogrusLogLevel(level loggers.Level) log.Level {
	switch level {
	case loggers.DebugLevel:
		return log.DebugLevel
	case loggers.InfoLevel:
		return log.InfoLevel
	case loggers.WarnLevel:
		return log.WarnLevel
	case loggers.ErrorLevel:
		return log.ErrorLevel
	default:
		return log.ErrorLevel
	}
}

func (l *logger) Log(ctx context.Context, level loggers.Level, skip int, args ...interface{}) {
	fields := make(log.Fields)

	// fetch fields from context and add them to logrus fields
	ctxFields := loggers.FromContext(ctx)
	if ctxFields != nil {
		for k, v := range ctxFields {
			fields[k] = v
		}
	}

	if l.opt.CallerInfo {
		_, file, line := loggers.FetchCallerInfo(skip+1, l.opt.CallerFileDepth)
		fields[l.opt.CallerFieldName] = fmt.Sprintf("%s:%d", file, line)
	}

	logger := l.logger.WithFields(fields)
	switch level {
	case loggers.DebugLevel:
		logger.Debug(args...)
	case loggers.InfoLevel:
		logger.Info(args...)
	case loggers.WarnLevel:
		logger.Warn(args...)
	case loggers.ErrorLevel:
		logger.Error(args...)
	default:
		l.logger.Error(args...)
	}
}

func (l *logger) SetLevel(level loggers.Level) {
	l.logger.SetLevel(toLogrusLogLevel(level))
}

func (l *logger) GetLevel() loggers.Level {
	switch l.logger.Level {
	case log.DebugLevel:
		return loggers.DebugLevel
	case log.InfoLevel:
		return loggers.InfoLevel
	case log.WarnLevel:
		return loggers.WarnLevel
	case log.ErrorLevel:
		return loggers.ErrorLevel
	default:
		return loggers.InfoLevel
	}
}

//NewLogger returns a BaseLogger impl for logrus
func NewLogger(options ...loggers.Option) loggers.BaseLogger {
	// default options
	opt := loggers.GetDefaultOptions()
	// read options
	for _, f := range options {
		f(&opt)
	}

	l := logger{}
	l.logger = log.New()
	l.logger.Out = os.Stdout

	l.logger.SetLevel(toLogrusLogLevel(opt.Level))

	fieldMap := log.FieldMap{
		log.FieldKeyTime:  opt.TimestampFieldName,
		log.FieldKeyLevel: opt.LevelFieldName,
	}
	//check JSON logs
	if opt.JSONLogs {
		l.logger.Formatter = &log.JSONFormatter{
			FieldMap: fieldMap,
		}
	} else {
		l.logger.Formatter = &log.TextFormatter{
			FullTimestamp: true,
		}
	}

	l.opt = opt

	if opt.ReplaceStdLogger {
		stdlog.SetOutput(l.logger.Writer())
	}
	return &l
}
