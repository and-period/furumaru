package database

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/stretchr/testify/require"
)

var (
	dbClient *database.Client
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

func newTestDBClient() (*database.Client, error) {
	setEnv()
	// テスト用Database接続用クライアントの生成
	params := &database.Params{
		Socket:   "tcp",
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: "messengers",
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	return database.NewClient(params)
}

func deleteAll(ctx context.Context) error {
	tables := []string{
		// テストに対応したテーブルから追記(削除順)
		scheduleTable,
		notificationTable,
		receivedQueueTable,
		reportTemplateTable,
		pushTemplateTable,
		messageTemplateTable,
		messageTable,
		contactTable,
		contactCategoryTable,
	}
	return delete(ctx, tables...)
}

func delete(ctx context.Context, tables ...string) error {
	for _, table := range tables {
		sql := fmt.Sprintf("DELETE FROM %s", table)
		if err := dbClient.DB.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

func TestDatabase(t *testing.T) {
	t.Parallel()
	require.NotNil(t, NewDatabase(&Params{}))
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
