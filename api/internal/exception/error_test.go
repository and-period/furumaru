package exception

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestInternalError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "err nil",
			err:    nil,
			expect: nil,
		},
		{
			name:   "context canceled",
			err:    context.Canceled,
			expect: ErrCanceled,
		},
		{
			name:   "context deadline exceeded",
			err:    context.DeadlineExceeded,
			expect: ErrDeadlineExceeded,
		},
		{
			name:   "media error",
			err:    media.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "messenger error",
			err:    messenger.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "store error",
			err:    store.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "user error",
			err:    user.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "other error",
			err:    assert.AnError,
			expect: ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, InternalError(tt.err), tt.expect)
		})
	}
}

func TestMediaError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    media.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "not found",
			err:    media.ErrNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    media.ErrAlreadyExists,
			expect: ErrAlreadyExists,
		},
		{
			name:   "failed precondition",
			err:    media.ErrFailedPrecondition,
			expect: ErrFailedPrecondition,
		},
		{
			name:   "canceled",
			err:    media.ErrCanceled,
			expect: ErrCanceled,
		},
		{
			name:   "deadline exceeded",
			err:    media.ErrDeadlineExceeded,
			expect: ErrDeadlineExceeded,
		},
		{
			name:   "other error",
			err:    assert.AnError,
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, mediaError(tt.err), tt.expect)
		})
	}
}

func TestMessengerError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    messenger.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "not found",
			err:    messenger.ErrNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    messenger.ErrAlreadyExists,
			expect: ErrAlreadyExists,
		},
		{
			name:   "forbidden",
			err:    messenger.ErrForbidden,
			expect: ErrForbidden,
		},
		{
			name:   "failed precondition",
			err:    messenger.ErrFailedPrecondition,
			expect: ErrFailedPrecondition,
		},
		{
			name:   "canceled",
			err:    messenger.ErrCanceled,
			expect: ErrCanceled,
		},
		{
			name:   "deadline exceeded",
			err:    messenger.ErrDeadlineExceeded,
			expect: ErrDeadlineExceeded,
		},
		{
			name:   "other error",
			err:    assert.AnError,
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, messengerError(tt.err), tt.expect)
		})
	}
}

func TestStoreError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    store.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "not found",
			err:    store.ErrNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    store.ErrAlreadyExists,
			expect: ErrAlreadyExists,
		},
		{
			name:   "forbidden",
			err:    store.ErrForbidden,
			expect: ErrForbidden,
		},
		{
			name:   "failed precondition",
			err:    store.ErrFailedPrecondition,
			expect: ErrFailedPrecondition,
		},
		{
			name:   "unavailable",
			err:    store.ErrUnavailable,
			expect: ErrUnavailable,
		},
		{
			name:   "canceled",
			err:    store.ErrCanceled,
			expect: ErrCanceled,
		},
		{
			name:   "deadline exceeded",
			err:    store.ErrDeadlineExceeded,
			expect: ErrDeadlineExceeded,
		},
		{
			name:   "other error",
			err:    assert.AnError,
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, storeError(tt.err), tt.expect)
		})
	}
}

func TestUserError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    user.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "unauthenticated",
			err:    user.ErrUnauthenticated,
			expect: ErrUnauthenticated,
		},
		{
			name:   "not found",
			err:    user.ErrNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    user.ErrAlreadyExists,
			expect: ErrAlreadyExists,
		},
		{
			name:   "forbidden",
			err:    user.ErrForbidden,
			expect: ErrForbidden,
		},
		{
			name:   "failed precondition",
			err:    user.ErrFailedPrecondition,
			expect: ErrFailedPrecondition,
		},
		{
			name:   "canceled",
			err:    user.ErrCanceled,
			expect: ErrCanceled,
		},
		{
			name:   "deadline exceeded",
			err:    user.ErrDeadlineExceeded,
			expect: ErrDeadlineExceeded,
		},
		{
			name:   "other error",
			err:    assert.AnError,
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, userError(tt.err), tt.expect)
		})
	}
}
