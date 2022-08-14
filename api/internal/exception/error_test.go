package exception

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
	"github.com/and-period/furumaru/api/pkg/line"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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
			name:   "internl error",
			err:    ErrInvalidArgument,
			expect: ErrInvalidArgument,
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
			name:   "validation error",
			err:    validator.ValidationErrors{},
			expect: ErrInvalidArgument,
		},
		{
			name:   "database error",
			err:    gorm.ErrInvalidField,
			expect: ErrInvalidArgument,
		},
		{
			name:   "auth error",
			err:    cognito.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "storage error",
			err:    storage.ErrInvalidURL,
			expect: ErrInvalidArgument,
		},
		{
			name:   "mailer error",
			err:    mailer.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "messaging error",
			err:    messaging.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "notifier error",
			err:    line.ErrInvalidArgument,
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

func TestRetryable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "err nil",
			err:    nil,
			expect: false,
		},
		{
			name:   "internl",
			err:    ErrInternal,
			expect: false,
		},
		{
			name:   "resource exhausted",
			err:    ErrResourceExhausted,
			expect: true,
		},
		{
			name:   "canceled",
			err:    ErrCanceled,
			expect: true,
		},
		{
			name:   "unablailable",
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
			expect: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, Retryable(tt.err))
		})
	}
}

func TestValidationError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    validator.ValidationErrors{},
			expect: ErrInvalidArgument,
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
			assert.ErrorIs(t, validationError(tt.err), tt.expect)
		})
	}
}

func TestDBError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "already exists",
			err:    &mysql.MySQLError{Number: 1062},
			expect: ErrAlreadyExists,
		},
		{
			name:   "other mysql error",
			err:    &mysql.MySQLError{Number: 0},
			expect: ErrInternal,
		},
		{
			name:   "invalid argument",
			err:    gorm.ErrInvalidField,
			expect: ErrInvalidArgument,
		},
		{
			name:   "not found",
			err:    gorm.ErrRecordNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "not implemented",
			err:    gorm.ErrNotImplemented,
			expect: ErrNotImplemented,
		},
		{
			name:   "internal",
			err:    gorm.ErrUnsupportedRelation,
			expect: ErrInternal,
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
			assert.ErrorIs(t, dbError(tt.err), tt.expect)
		})
	}
}

func TestAuthError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    cognito.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "unauthenticated",
			err:    cognito.ErrUnauthenticated,
			expect: ErrUnauthenticated,
		},
		{
			name:   "not found",
			err:    cognito.ErrNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    cognito.ErrAlreadyExists,
			expect: ErrAlreadyExists,
		},
		{
			name:   "resource exhausted",
			err:    cognito.ErrResourceExhausted,
			expect: ErrResourceExhausted,
		},
		{
			name:   "canceled",
			err:    cognito.ErrCanceled,
			expect: ErrCanceled,
		},
		{
			name:   "timeout",
			err:    cognito.ErrTimeout,
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
			assert.ErrorIs(t, authError(tt.err), tt.expect)
		})
	}
}

func TestStorageError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    storage.ErrInvalidURL,
			expect: ErrInvalidArgument,
		},
		{
			name:   "not found",
			err:    storage.ErrNotFound,
			expect: ErrNotFound,
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
			assert.ErrorIs(t, storageError(tt.err), tt.expect)
		})
	}
}

func TestMailerError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    mailer.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "unauthenticated",
			err:    mailer.ErrUnauthenticated,
			expect: ErrUnauthenticated,
		},
		{
			name:   "permission denied",
			err:    mailer.ErrPermissionDenied,
			expect: ErrForbidden,
		},
		{
			name:   "payload too long",
			err:    mailer.ErrPayloadTooLong,
			expect: ErrResourceExhausted,
		},
		{
			name:   "not found",
			err:    mailer.ErrNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "internal",
			err:    mailer.ErrInternal,
			expect: ErrInternal,
		},
		{
			name:   "unavailable",
			err:    mailer.ErrUnavailable,
			expect: ErrUnavailable,
		},
		{
			name:   "canceled",
			err:    mailer.ErrCanceled,
			expect: ErrCanceled,
		},
		{
			name:   "timeout",
			err:    mailer.ErrTimeout,
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
			assert.ErrorIs(t, mailerError(tt.err), tt.expect)
		})
	}
}

func TestMessagingError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    messaging.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "unauthenticated",
			err:    messaging.ErrUnauthenticated,
			expect: ErrUnauthenticated,
		},
		{
			name:   "not found",
			err:    messaging.ErrNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "resource exhausted",
			err:    messaging.ErrResourceExhausted,
			expect: ErrResourceExhausted,
		},
		{
			name:   "internal",
			err:    messaging.ErrInternal,
			expect: ErrInternal,
		},
		{
			name:   "unavailable",
			err:    messaging.ErrUnavailable,
			expect: ErrUnavailable,
		},
		{
			name:   "canceled",
			err:    messaging.ErrCanceled,
			expect: ErrCanceled,
		},
		{
			name:   "timeout",
			err:    messaging.ErrTimeout,
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
			assert.ErrorIs(t, messagingError(tt.err), tt.expect)
		})
	}
}

func TestNotifierError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "invalid argument",
			err:    line.ErrInvalidArgument,
			expect: ErrInvalidArgument,
		},
		{
			name:   "unauthenticated",
			err:    line.ErrUnauthenticated,
			expect: ErrUnauthenticated,
		},
		{
			name:   "permission denied",
			err:    line.ErrPermissionDenied,
			expect: ErrForbidden,
		},
		{
			name:   "payload too long",
			err:    line.ErrPayloadTooLong,
			expect: ErrResourceExhausted,
		},
		{
			name:   "not found",
			err:    line.ErrNotFound,
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    line.ErrAlreadyExists,
			expect: ErrAlreadyExists,
		},
		{
			name:   "internal",
			err:    line.ErrInternal,
			expect: ErrInternal,
		},
		{
			name:   "unavailable",
			err:    line.ErrUnavailable,
			expect: ErrUnavailable,
		},
		{
			name:   "resource exhausted",
			err:    line.ErrResourceExhausted,
			expect: ErrResourceExhausted,
		},
		{
			name:   "canceled",
			err:    line.ErrCanceled,
			expect: ErrCanceled,
		},
		{
			name:   "timeout",
			err:    line.ErrTimeout,
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
			assert.ErrorIs(t, notifierError(tt.err), tt.expect)
		})
	}
}
