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
	beforeSendFn := func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
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
	dopts := &options{
		bind: false,
		opts: sentry.ClientOptions{
			Environment:   "",
			Debug:         false,
			EnableTracing: false,
			BeforeSend:    beforeSendFn,
		},
	}
	for i := range opts {
		opts[i](dopts)
	}
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

func WithTrace(enable bool) ClientOption {
	return func(o *options) {
		o.opts.EnableTracing = enable
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
