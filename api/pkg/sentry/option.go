package sentry

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

type ClientOption func(*sentry.ClientOptions)

func WithDSN(dsn string) ClientOption {
	return func(opts *sentry.ClientOptions) {
		opts.Dsn = dsn
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
