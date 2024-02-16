package sentry

import (
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
)

type options struct {
	bind bool
	opts sentry.ClientOptions
}

type ClientOption func(*options)

func buildOptions(opts ...ClientOption) *options {
	dopts := &options{
		bind: false,
		opts: sentry.ClientOptions{
			ServerName:         "",
			Environment:        "",
			Debug:              false,
			EnableTracing:      false,
			TracesSampleRate:   0.0,
			ProfilesSampleRate: 0.0,
		},
	}
	for i := range opts {
		opts[i](dopts)
	}
	dopts.opts.BeforeSend = func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
		for i := range event.Exception {
			exception := &event.Exception[i]
			if !strings.Contains(exception.Type, "wrapError") {
				continue
			}
			strs := strings.SplitN(exception.Value, ":", 2)
			if len(strs) != 2 {
				continue
			}
			exception.Type, exception.Value = strs[0], strs[1]
		}
		return event
	}
	dopts.opts.TracesSampler = sentry.TracesSampler(func(ctx sentry.SamplingContext) float64 {
		if ctx.Span.Name == "GET /health" {
			return 0.0
		}
		return dopts.opts.TracesSampleRate
	})
	return dopts
}

func WithBind(bind bool) ClientOption {
	return func(o *options) {
		o.bind = bind
	}
}

func WithDSN(dsn string) ClientOption {
	return func(o *options) {
		o.opts.Dsn = dsn
	}
}

func WithServerName(name string) ClientOption {
	return func(o *options) {
		o.opts.ServerName = name
	}
}

func WithEnvironment(env string) ClientOption {
	return func(o *options) {
		o.opts.Environment = env
	}
}

func WithDebug(debug bool) ClientOption {
	return func(o *options) {
		o.opts.Debug = debug
	}
}

func WithTracesSampleRate(rate float64) ClientOption {
	return func(o *options) {
		if rate > 0.0 {
			o.opts.EnableTracing = true
		}
		o.opts.TracesSampleRate = rate
	}
}

func WithProfilesSampleRate(rate float64) ClientOption {
	return func(o *options) {
		o.opts.ProfilesSampleRate = rate
	}
}

type ReportOption func(*sentry.Scope)

func (o ReportOption) String() string {
	return "sentry:ReportOption"
}

func attachReportOption(scope *sentry.Scope, opts ...ReportOption) {
	for i := range opts {
		opts[i](scope)
	}
}

func WithLevel(level string) ReportOption {
	return func(scope *sentry.Scope) {
		scope.SetLevel(sentry.Level(level))
	}
}

func WithTraceID(traceID string) ReportOption {
	return func(scope *sentry.Scope) {
		scope.SetTag("trace_id", traceID)
	}
}

func WithTag(key, value string) ReportOption {
	return func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	}
}

func WithTags(tags map[string]string) ReportOption {
	return func(scope *sentry.Scope) {
		scope.SetTags(tags)
	}
}

func WithExtras(extras map[string]interface{}) ReportOption {
	return func(scope *sentry.Scope) {
		scope.SetExtras(extras)
	}
}

func WithRequest(req *http.Request) ReportOption {
	return func(scope *sentry.Scope) {
		scope.SetRequest(req)
	}
}

func WithFingerprint(fingerprint ...string) ReportOption {
	return func(scope *sentry.Scope) {
		scope.SetFingerprint(fingerprint)
	}
}

type User struct {
	ID        string
	Email     string
	IPAddress string
	Username  string
	Name      string
	Data      map[string]string
}

func WithUser(user *User) ReportOption {
	return func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:        user.ID,
			Email:     user.Email,
			IPAddress: user.IPAddress,
			Username:  user.Username,
			Data:      user.Data,
		})
	}
}
