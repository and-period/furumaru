package backoff

import (
	"context"
	"fmt"
)

type DoFunc func() error

type Retryable func(err error) bool

type options struct {
	retryable Retryable
}

type Option func(*options)

func WithRetryablel(retryable Retryable) Option {
	return func(opts *options) {
		opts.retryable = retryable
	}
}

func Retry(ctx context.Context, backoff Backoff, doFunc DoFunc, opts ...Option) (err error) {
	dopts := &options{}
	for i := range opts {
		opts[i](dopts)
	}
	for backoff.Continue() {
		err = doFunc()
		if err == nil {
			return nil
		}
		if dopts.retryable != nil && !dopts.retryable(err) {
			return err
		}
		select {
		case <-backoff.Wait():
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("backoff: retry limit exceeded: %w", err)
}
