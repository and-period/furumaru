package sentry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
)

func TestReportOption(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	opts := []ReportOption{
		WithLevel("error"),
		WithTraceID("trace-id"),
		WithTag("tag", "value"),
		WithTags(map[string]string{"key": "value"}),
		WithFingerprint("fingerprint"),
		WithRequest(req),
		WithUser(&User{}),
	}
	scope := sentry.NewScope()
	attachReportOption(scope, opts...)
}
