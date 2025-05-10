package log

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/pkg/sentry"
	sentrygo "github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrSync = errors.New("log: failed to sentry flush")

type sentryCore struct {
	core         zapcore.Core
	client       sentry.Client
	level        zapcore.Level
	sentryLevel  zapcore.Level
	fields       []zapcore.Field
	flushTimeout time.Duration
}

func NewSentryLogger(dsn string, opts ...Option) (*zap.Logger, error) {
	dopts := buildOptions(opts...)
	logger, err := newLogger(dopts)
	if err != nil {
		return nil, err
	}
	sopts := []sentry.ClientOption{
		sentry.WithServerName(dopts.sentryServerName),
		sentry.WithEnvironment(dopts.sentryEnvironment),
		sentry.WithDSN(dsn),
	}
	sclient, err := sentry.NewClient(sopts...)
	if err != nil {
		return nil, err
	}
	logLevel := getLogLevel(dopts.logLevel)
	sentryLevel := getLogLevel(dopts.sentryLevel)
	core := &sentryCore{
		core:         logger.Core(),
		client:       sclient,
		level:        minLevel(sentryLevel, logLevel),
		sentryLevel:  sentryLevel,
		flushTimeout: dopts.sentryFlushTimeout,
	}
	return zap.New(core), nil
}

func minLevel(lv1, lv2 zapcore.Level) zapcore.Level {
	if lv1 < lv2 {
		return lv1
	}
	return lv2
}

func (c *sentryCore) Enabled(level zapcore.Level) bool {
	return level >= c.level
}

func (c *sentryCore) With(fields []zap.Field) zapcore.Core {
	return &sentryCore{
		core:        c.core.With(fields),
		client:      c.client,
		level:       c.level,
		sentryLevel: c.sentryLevel,
		fields:      c.join(fields),
	}
}

func (c *sentryCore) Check(entry zapcore.Entry, centry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if entry.Level >= c.sentryLevel {
		centry = centry.AddCore(entry, c)
	}
	return c.core.Check(entry, centry)
}

func (c *sentryCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	fs := ZapFields(c.join(fields))
	opts := fs.SentryOptions()
	opts = append(opts, sentry.WithLevel(string(getSentryLevel(c.sentryLevel))))
	if ok, err := fs.FirstError(); ok {
		extras := map[string]interface{}{
			"message":    entry.Message,
			"stacktrace": entry.Stack,
		}
		opts = append(opts, sentry.WithExtras(extras))
		c.client.ReportError(context.Background(), err, opts...)
		return nil
	}
	msg := fmt.Sprintf("%s\n\nstacktrace:\n%s", entry.Message, entry.Stack)
	c.client.ReportMessage(context.Background(), msg, opts...)
	return nil // Check()でcore追加しているためここではWrite()を呼ばない。※呼ぶと重複出力されるため
}

func (c *sentryCore) Sync() error {
	if !c.client.Flush(c.flushTimeout) {
		return ErrSync
	}
	return c.core.Sync()
}

func (c *sentryCore) join(fields []zapcore.Field) []zapcore.Field {
	res := make([]zapcore.Field, 0, len(fields)+len(c.fields))
	res = append(res, c.fields...)
	res = append(res, fields...)
	return res
}

type ZapFields []zapcore.Field

func (fs ZapFields) FirstError() (bool, error) {
	for _, f := range fs {
		if f.Type != zapcore.ErrorType {
			continue
		}
		//nolint:forcetypeassert
		return true, f.Interface.(error)
	}
	return false, nil
}

func (fs ZapFields) SentryOptions() []sentry.ReportOption {
	res := make([]sentry.ReportOption, 0, len(fs))
	extras := make(map[string]interface{}, len(fs))
	for _, f := range fs {
		switch f.Type {
		case zapcore.ErrorType:
		case zapcore.StringType:
			extras[f.Key] = f.String
		case zapcore.BoolType:
			// https://github.com/uber-go/zap/blob/v1.26.0/field.go#L62-L69
			extras[f.Key] = f.Integer == 1
		case zapcore.TimeType:
			// https://github.com/uber-go/zap/blob/v1.26.0/field.go#L345-L352
			sec := f.Integer / int64(time.Second)
			nsec := f.Integer % int64(time.Second)
			timestamp := time.Unix(sec, nsec)
			extras[f.Key] = timestamp.Format(time.RFC3339Nano)
		case zapcore.DurationType:
			extras[f.Key] = time.Duration(f.Integer).String()
		case zapcore.Float32Type, zapcore.Float64Type, zapcore.UintptrType,
			zapcore.Int8Type, zapcore.Uint8Type, zapcore.Int16Type, zapcore.Uint16Type,
			zapcore.Int32Type, zapcore.Uint32Type, zapcore.Int64Type, zapcore.Uint64Type:
			extras[f.Key] = f.Integer
		}
		if f.Interface == nil {
			continue
		}
		if opt, ok := sentryOption(f); ok {
			res = append(res, opt)
			continue
		}
		extras[f.Key] = f.Interface
	}
	if len(extras) > 0 {
		res = append(res, sentry.WithExtras(extras))
	}
	return res
}

func sentryOption(f zapcore.Field) (sentry.ReportOption, bool) {
	if f.Type != zapcore.StringerType {
		return nil, false
	}
	opt, ok := f.Interface.(sentry.ReportOption)
	return opt, ok
}

func getSentryLevel(level zapcore.Level) sentrygo.Level {
	switch level {
	case zapcore.DebugLevel:
		return sentrygo.LevelDebug
	case zapcore.InfoLevel:
		return sentrygo.LevelInfo
	case zapcore.WarnLevel:
		return sentrygo.LevelWarning
	case zapcore.ErrorLevel:
		return sentrygo.LevelError
	default:
		return sentrygo.LevelWarning
	}
}
