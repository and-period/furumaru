//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package sentry

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

type Client interface {
	ReportError(ctx context.Context, err error, opts ...ReportOption)
	ReportPanic(ctx context.Context, err error, opts ...ReportOption)
	ReportMessage(ctx context.Context, msg string, opts ...ReportOption)
	Flush(timeout time.Duration) bool
}

type client struct {
	client *sentry.Client
}

func NewZapHookFn(client Client, opts ...ReportOption) func(entry zapcore.Entry) error {
	return func(entry zapcore.Entry) error {
		if entry.Level < zapcore.ErrorLevel {
			return nil
		}
		opts = append(opts, WithLevel("error"))
		msg := entry.Message
		if len(entry.Stack) > 0 {
			msg = fmt.Sprintf("%s\n\nstacktrace:\n%s", entry.Message, entry.Stack)
		}
		client.ReportMessage(context.Background(), msg, opts...)
		return nil
	}
}

func NewClient(opts ...ClientOption) (Client, error) {
	dopts := &sentry.ClientOptions{}
	for i := range opts {
		opts[i](dopts)
	}
	if emptySentryKey(dopts) {
		return NewFixedMockClient(), nil
	}
	cli, err := sentry.NewClient(*dopts)
	if err != nil {
		return nil, err
	}
	return &client{client: cli}, nil
}

func emptySentryKey(opts *sentry.ClientOptions) bool {
	if opts.Dsn != "" {
		return false
	}
	return os.Getenv("SENTRY_DSN") == ""
}

func (c *client) ReportError(ctx context.Context, err error, opts ...ReportOption) {
	hint := &sentry.EventHint{
		Context:           ctx,
		OriginalException: err,
	}
	scope := sentry.NewScope()
	attachReportOption(scope, opts...)
	c.client.CaptureException(err, hint, scope)
}

func (c *client) ReportPanic(ctx context.Context, err error, opts ...ReportOption) {
	hint := &sentry.EventHint{
		Context:            ctx,
		RecoveredException: err,
	}
	scope := sentry.NewScope()
	attachReportOption(scope, opts...)
	c.client.RecoverWithContext(ctx, err, hint, scope)
}

func (c *client) ReportMessage(ctx context.Context, msg string, opts ...ReportOption) {
	hint := &sentry.EventHint{
		Context: ctx,
	}
	scope := sentry.NewScope()
	attachReportOption(scope, opts...)
	c.client.CaptureMessage(msg, hint, scope)
}

func (c *client) Flush(timeout time.Duration) bool {
	return c.client.Flush(timeout)
}
