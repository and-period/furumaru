package backoff

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		timeout    time.Duration
		maxRetries int64
		fn         func() error
		retryable  func(err error) bool
		hasErr     bool
	}{
		{
			name:       "non error",
			timeout:    1 * time.Second,
			maxRetries: 1,
			fn: func() error {
				return nil
			},
			retryable: func(err error) bool {
				return false
			},
			hasErr: false,
		},
		{
			name:       "error with retry",
			timeout:    1 * time.Second,
			maxRetries: 1,
			fn: func() error {
				return assert.AnError
			},
			retryable: func(err error) bool {
				return true
			},
			hasErr: true,
		},
		{
			name:       "error without retry",
			timeout:    1 * time.Second,
			maxRetries: 1,
			fn: func() error {
				return assert.AnError
			},
			retryable: func(err error) bool {
				return false
			},
			hasErr: true,
		},
		{
			name:       "context timeout",
			timeout:    10 * time.Millisecond,
			maxRetries: 10,
			fn: func() error {
				return assert.AnError
			},
			retryable: func(err error) bool {
				return true
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), tt.timeout)
			defer cancel()
			backoff := NewFixedIntervalBackoff(10*time.Millisecond, tt.maxRetries)
			err := Retry(ctx, backoff, tt.fn, WithRetryablel(tt.retryable))
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}
