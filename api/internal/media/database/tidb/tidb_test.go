package tidb

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDBError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "not error",
			err:    nil,
			expect: nil,
		},
		{
			name:   "database error",
			err:    database.ErrFailedPrecondition,
			expect: database.ErrFailedPrecondition,
		},
		{
			name:   "context canceled",
			err:    context.Canceled,
			expect: database.ErrCanceled,
		},
		{
			name:   "context deadline exceeded",
			err:    context.DeadlineExceeded,
			expect: database.ErrDeadlineExceeded,
		},
		{
			name:   "mysql already exists",
			err:    &mysql.MySQLError{Number: 1062},
			expect: database.ErrAlreadyExists,
		},
		{
			name:   "other mysql error",
			err:    &mysql.MySQLError{},
			expect: database.ErrInternal,
		},
		{
			name:   "invalid argument",
			err:    gorm.ErrInvalidValue,
			expect: database.ErrInvalidArgument,
		},
		{
			name:   "not found",
			err:    gorm.ErrRecordNotFound,
			expect: database.ErrNotFound,
		},
		{
			name:   "internal",
			err:    gorm.ErrUnsupportedRelation,
			expect: database.ErrInternal,
		},
		{
			name:   "unknown",
			err:    assert.AnError,
			expect: database.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := dbError(tt.err)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}
