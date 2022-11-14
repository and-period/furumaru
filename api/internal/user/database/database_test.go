package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type mocks struct {
	db *database.Client
}

func newMocks(ctrl *gomock.Controller) (*mocks, error) {
	setEnv()
	// テスト用Database接続用クライアントの生成
	params := &database.Params{
		Socket:   "tcp",
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: "users",
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	db, err := database.NewClient(params)
	if err != nil {
		return nil, err
	}

	return &mocks{db: db}, nil
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
