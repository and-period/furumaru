//go:build integration

package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContainerDB(t *testing.T) {
	ctx := t.Context()

	client, cleanup, err := NewContainerDB(ctx,
		WithContainerDatabase("testdb"),
	)
	require.NoError(t, err)
	defer cleanup()

	require.NotNil(t, client)
	require.NotNil(t, client.DB)

	// テーブル作成
	err = client.DB.Exec(`
		CREATE TABLE IF NOT EXISTS test_items (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	require.NoError(t, err)

	// データ挿入
	err = client.DB.Exec(`INSERT INTO test_items (name) VALUES ('item1'), ('item2'), ('item3')`).Error
	require.NoError(t, err)

	// データ取得
	var count int64
	err = client.DB.Raw(`SELECT COUNT(*) FROM test_items`).Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 名前で検索
	var name string
	err = client.DB.Raw(`SELECT name FROM test_items WHERE id = 1`).Scan(&name).Error
	require.NoError(t, err)
	assert.Equal(t, "item1", name)
}

func TestNewContainerDB_WithOptions(t *testing.T) {
	ctx := t.Context()

	client, cleanup, err := NewContainerDB(ctx,
		WithContainerImage("mysql:8.0"),
		WithContainerDatabase("custom_db"),
		WithContainerUsername("root"),
		WithContainerPassword("testpass"),
	)
	require.NoError(t, err)
	defer cleanup()

	require.NotNil(t, client)

	// 接続確認
	var result int
	err = client.DB.Raw(`SELECT 1`).Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 1, result)
}
