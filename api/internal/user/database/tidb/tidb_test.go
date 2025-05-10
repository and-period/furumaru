package tidb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/pkg/mysql"
	gmysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	dbClient *mysql.Client
	current  = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
)

func TestMain(m *testing.M) {
	setEnv()

	client, err := newTestDBClient()
	if err != nil {
		panic(err)
	}
	dbClient = client

	os.Exit(m.Run())
}

func newTestDBClient() (*mysql.Client, error) {
	setEnv()
	// テスト用Database接続用クライアントの生成
	params := &mysql.Params{
		Socket:   "tcp",
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	logger, _ := zap.NewDevelopment()
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return mysql.NewClient(params, mysql.WithLogger(logger))
	case "tidb":
		return mysql.NewTiDBClient(params, mysql.WithLogger(logger))
	default:
		return nil, fmt.Errorf("unsupported driver: %s", os.Getenv("DB_DRIVER"))
	}
}

func deleteAll(ctx context.Context) error {
	tables := []string{
		// テストに対応したテーブルから追記(削除順)
		adminGroupUserTable,
		adminGroupRoleTable,
		adminGroupTable,
		adminRolePolicyTable,
		adminRoleTable,
		adminPolicyTable,
		adminAuthProviderTable,
		addressRevisionTable,
		addressTable,
		producerTable,
		coordinatorTable,
		administratorTable,
		adminTable,
		guestTable,
		memberTable,
		userNotificationTable,
		userTable,
	}
	if err := dbClient.DB.WithContext(ctx).Exec("SET foreign_key_checks = 0").Error; err != nil {
		return err
	}
	defer func() {
		if err := dbClient.DB.WithContext(ctx).Exec("SET foreign_key_checks = 1").Error; err != nil {
			fmt.Printf("mysql: failed to delete all: %s", err.Error())
		}
	}()
	return delete(ctx, tables...)
}

func delete(ctx context.Context, tables ...string) error {
	for _, table := range tables {
		sql := fmt.Sprintf("DELETE FROM %s", table)
		if err := dbClient.DB.WithContext(ctx).Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

func TestDatabase(t *testing.T) {
	t.Parallel()
	require.NotNil(t, NewDatabase(nil))
}

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
			err:    &gmysql.MySQLError{Number: 1062},
			expect: database.ErrAlreadyExists,
		},
		{
			name:   "other mysql error",
			err:    &gmysql.MySQLError{},
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

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := dbError(tt.err)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}

func setEnv() {
	if os.Getenv("DB_DRIVER") == "" {
		os.Setenv("DB_DRIVER", "mysql")
	}
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "127.0.0.1")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "3306")
	}
	if os.Getenv("DB_DATABASE") == "" {
		os.Setenv("DB_DATABASE", "users")
	}
	if os.Getenv("DB_USERNAME") == "" {
		os.Setenv("DB_USERNAME", "root")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "")
	}
}
