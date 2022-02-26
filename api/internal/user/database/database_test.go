package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	mock_cognito "github.com/and-period/marche/api/mock/pkg/cognito"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var errmock = errors.New("some error")

type mocks struct {
	db       *database.Client
	userAuth *mock_cognito.MockClient
}

func newMocks(ctrl *gomock.Controller) (*mocks, error) {
	setEnv()
	// テスト用Database接続用クライアントの生成
	params := &database.Params{
		Socket:        "tcp",
		Host:          os.Getenv("DB_HOST"),
		Port:          os.Getenv("DB_PORT"),
		Database:      "users",
		Username:      os.Getenv("DB_USERNAME"),
		Password:      os.Getenv("DB_PASSWORD"),
		DisableLogger: true,
	}
	db, err := database.NewClient(params)
	if err != nil {
		return nil, err
	}

	return &mocks{
		db:       db,
		userAuth: mock_cognito.NewMockClient(ctrl),
	}, nil
}

func (m *mocks) dbDelete(ctx context.Context, tables ...string) error {
	for _, table := range tables {
		sql := fmt.Sprintf("DELETE FROM %s", table)
		if err := m.db.DB.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

func TestDatabase(t *testing.T) {
	t.Parallel()
	require.NotNil(t, NewDatabase(&Params{}))
}

func TestDBError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "error is nil",
			err:    nil,
			expect: nil,
		},
		{
			name:   "mysql error",
			err:    &mysql.MySQLError{},
			expect: ErrUnknown,
		},
		{
			name:   "record duplicated",
			err:    &mysql.MySQLError{Number: 1062},
			expect: ErrAlreadyExists,
		},
		{
			name:   "invalid argument",
			err:    gorm.ErrEmptySlice,
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
			err:    gorm.ErrUnsupportedDriver,
			expect: ErrInternal,
		},
		{
			name:   "other error",
			err:    errors.New("some error"),
			expect: ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, dbError(tt.err), tt.expect)
		})
	}
}

func setEnv() {
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "127.0.0.1")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "3326")
	}
	if os.Getenv("DB_USERNAME") == "" {
		os.Setenv("DB_USERNAME", "root")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "12345678")
	}
}
