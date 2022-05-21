package exception

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/cognito"
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
