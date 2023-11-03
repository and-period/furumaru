package sentry

import (
	"context"
	"time"
)

type fixedMockClient struct{}

func NewFixedMockClient() Client {
	return &fixedMockClient{}
}

func (c *fixedMockClient) ReportError(_ context.Context, _ error, _ ...ReportOption) {}

func (c *fixedMockClient) ReportPanic(_ context.Context, _ error, _ ...ReportOption) {}

func (c *fixedMockClient) ReportMessage(_ context.Context, _ string, _ ...ReportOption) {}

func (c *fixedMockClient) Flush(_ time.Duration) bool {
	return true
}
