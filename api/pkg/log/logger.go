package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type options struct {
	logLevel           string
	outputPath         string
	sentryServerName   string
	sentryEnvironment  string
	sentryLevel        string
	sentryFlushTimeout time.Duration
}

type Option func(opts *options)

func buildOptions(opts ...Option) *options {
	dopts := &options{
		logLevel:           "info",
		outputPath:         "",
		sentryServerName:   "",
		sentryEnvironment:  "",
		sentryLevel:        "warn",
		sentryFlushTimeout: 5 * time.Second,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return dopts
}

func WithLogLevel(level string) Option {
	return func(opts *options) {
		opts.logLevel = level
	}
}

func WithSentryServerName(name string) Option {
	return func(opts *options) {
		opts.sentryServerName = name
	}
}

func WithSentryEnvironment(env string) Option {
	return func(opts *options) {
		opts.sentryEnvironment = env
	}
}

func WithSentryLevel(level string) Option {
	return func(opts *options) {
		opts.sentryLevel = level
	}
}

func WithOutput(path string) Option {
	return func(opts *options) {
		opts.outputPath = path
	}
}

func WithSentryFlushTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.sentryFlushTimeout = timeout
	}
}

// NewLogger - ログ出力用クライアントの生成
func NewLogger(opts ...Option) (*zap.Logger, error) {
	dopts := buildOptions(opts...)
	return newLogger(dopts)
}

func newLogger(opts *options) (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 標準出力設定
	level := getLogLevel(opts.logLevel)
	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Path==""のとき、標準出力のみ
	if opts.outputPath == "" {
		logger := zap.New(zapcore.NewTee(consoleCore))
		return logger, nil
	}

	// logPath!==""のとき、ファイル出力も追加
	outputPath := fmt.Sprintf("%s/outputs.log", opts.outputPath)
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o755)
	if err != nil {
		return nil, err
	}

	// ファイル出力設定
	logCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(file), level)
	logger := zap.New(zapcore.NewTee(consoleCore, logCore))
	return logger, nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
