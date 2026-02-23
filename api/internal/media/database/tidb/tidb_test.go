package tidb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/pkg/mysql"
	gmysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	dbClient *mysql.Client
	current  = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
)

func TestMain(m *testing.M) {
	setEnv()
	ctx := context.Background()

	client, cleanup, err := newTestDBClient(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	dbClient = client

	os.Exit(m.Run())
}

func newTestDBClient(ctx context.Context) (*mysql.Client, func(), error) {
	// 外部DBが設定されている場合は従来の外部DB接続を使用
	if !mysql.ShouldUseContainerDB() {
		return newExternalDBClient()
	}

	// コンテナベースのDB接続
	schemaDir, err := schemaDir()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve schema dir: %w", err)
	}

	client, cleanup, err := mysql.NewContainerDB(ctx,
		mysql.WithContainerDatabase("media"),
		mysql.WithSchemaDir(schemaDir),
	)
	if err != nil {
		return nil, nil, err
	}
	return client, cleanup, nil
}

func newExternalDBClient() (*mysql.Client, func(), error) {
	params := &mysql.Params{
		Socket:   "tcp",
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	var client *mysql.Client
	var err error
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		client, err = mysql.NewClient(params)
	case "tidb":
		client, err = mysql.NewTiDBClient(params)
	default:
		return nil, nil, fmt.Errorf("unsupported driver: %s", os.Getenv("DB_DRIVER"))
	}
	if err != nil {
		return nil, nil, err
	}
	noop := func() {}
	return client, noop, nil
}

// schemaDir はスキーマSQLファイルのディレクトリパスを返す。
// runtime.Caller を使ってこのファイルの位置からプロジェクトルートを辿り、
// infra/tidb/schema/media/ を解決する。
func schemaDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}

	// テストファイルの位置から api ディレクトリ（go.mod がある場所）を探す
	apiRoot, err := mysql.FindProjectRoot(filepath.Dir(filename))
	if err != nil {
		return "", err
	}

	// api/ の親がリポジトリルート
	repoRoot := filepath.Dir(apiRoot)
	dir := filepath.Join(repoRoot, "infra", "tidb", "schema", "media")

	// ディレクトリの存在確認
	if _, err := os.Stat(dir); err != nil {
		return "", fmt.Errorf("schema dir not found: %w", err)
	}

	return dir, nil
}

func deleteAll(ctx context.Context) error {
	tables := []string{
		// テストに対応したテーブルから追記(削除順)
		broadcastViewerLogTable,
		broadcastCommentTable,
		broadcastTable,
		videoViewerLogTable,
		videoCommentTable,
		videoProductTable,
		videoExperienceTable,
		videoTable,
	}
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
		os.Setenv("DB_DATABASE", "media")
	}
	if os.Getenv("DB_USERNAME") == "" {
		os.Setenv("DB_USERNAME", "root")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "")
	}
}
