package exception

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRetryable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "nil",
			err:    nil,
			expect: false,
		},
		{
			name:   "invalid argument",
			err:    ErrInvalidArgument,
			expect: false,
		},
		{
			name:   "unauthenticated",
			err:    ErrUnauthenticated,
			expect: false,
		},
		{
			name:   "forbidden",
			err:    ErrForbidden,
			expect: false,
		},
		{
			name:   "not found",
			err:    ErrNotFound,
			expect: false,
		},
		{
			name:   "already exists",
			err:    ErrAlreadyExists,
			expect: false,
		},
		{
			name:   "failed precondition",
			err:    ErrFailedPrecondition,
			expect: false,
		},
		{
			name:   "resource exhausted",
			err:    ErrResourceExhausted,
			expect: true,
		},
		{
			name:   "not implemented",
			err:    ErrNotImplemented,
			expect: false,
		},
		{
			name:   "internal error",
			err:    ErrInternal,
			expect: true,
		},
		{
			name:   "canceled",
			err:    ErrCanceled,
			expect: false,
		},
		{
			name:   "unavailable",
			err:    ErrUnavailable,
			expect: true,
		},
		{
			name:   "deadline exceeded",
			err:    ErrDeadlineExceeded,
			expect: true,
		},
		{
			name:   "out of range",
			err:    ErrOutOfRange,
			expect: false,
		},
		{
			name:   "unknown",
			err:    ErrUnknown,
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, IsRetryable(tt.err))
		})
	}
}
