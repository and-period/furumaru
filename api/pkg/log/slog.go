package log

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	slogmulti "github.com/samber/slog-multi"
)

var _ slog.Handler = (*handler)(nil)

type handler struct {
	slog.Handler
}

func middleware(h slog.Handler) slog.Handler {
	return &handler{
		Handler: h,
	}
}

func Start(ctx context.Context, opts ...Option) (func(), error) {
	dopts := buildOptions(opts...)

	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slogLevel(dopts.logLevel),
	})
	handlers := []slog.Handler{
		stdoutHandler,
	}

	stopSentry := startSentry(ctx, dopts, func(handler slog.Handler) {
		handlers = append(handlers, handler)
	})

	stopFn := func() {
		var wg sync.WaitGroup
		wg.Add(1) //nolint:mnd
		go func() {
			defer wg.Done()
			stopSentry()
		}()
		wg.Wait()
	}

	logger := slog.New(slogmulti.Pipe(middleware).Handler(slogmulti.Fanout(handlers...)))
	slog.SetDefault(logger)
	return stopFn, nil
}

func startSentry(ctx context.Context, opts *options, handlerFunc func(slog.Handler)) func() {
	sopts := sentry.ClientOptions{
		Dsn:         opts.sentryDSN,
		ServerName:  opts.sentryServerName,
		Environment: opts.sentryEnvironment,
	}
	sentry.Init(sopts)
	fn := func() {
		const timeout = time.Second * 5
		_ = sentry.Flush(timeout)
	}
	lopts := &sentryslog.Option{
		EventLevel: sentryLevelsFromMinimum(slogLevel(opts.sentryLevel)),
		AddSource:  true,
	}
	handler := lopts.NewSentryHandler(ctx)
	handlerFunc(handler)
	return fn
}

func sentryLevelsFromMinimum(minLevel slog.Level) []slog.Level {
	levels := []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
		sentryslog.LevelFatal,
	}
	var res []slog.Level
	for _, level := range levels {
		if level < minLevel {
			continue
		}
		res = append(res, level)
	}
	return res
}

func slogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
