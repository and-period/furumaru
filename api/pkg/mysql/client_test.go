//nolint:paralleltest
package mysql

import (
	"context"
	"os"
	"testing"
	"time"

	dmysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewClient(t *testing.T) {
	t.Skip() // 確認したい時のみコメントアウト

	setEnv()

	tests := []struct {
		name   string
		params *Params
		isErr  bool
	}{
		{
			name: "success",
			params: &Params{
				Socket:   "tcp",
				Host:     os.Getenv("DB_HOST"),
				Port:     os.Getenv("DB_PORT"),
				Database: "users",
				Username: os.Getenv("DB_USERNAME"),
				Password: os.Getenv("DB_PASSWORD"),
			},
			isErr: false,
		},
		{
			name: "failed to connect mysql",
			params: &Params{
				Socket:   "tcp",
				Host:     "127.0.0.1",
				Port:     "80",
				Database: "users",
				Username: "",
				Password: "",
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(
				tt.params,
				WithNow(time.Now),
				WithLocation(time.UTC),
				WithCharset("utf8mb4"),
				WithCollation("utf8mb4_general_ci"),
				WithTLS(false),
				WithNativePasswords(true),
				WithMaxAllowedPacket(4194304),
			)
			if tt.isErr {
				assert.Error(t, err)
				assert.Nil(t, client)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}

func TestBeginAndClose(t *testing.T) {
	t.Skip() // 確認したい時のみコメントアウト

	setEnv()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	params := &Params{
		Socket:   "tcp",
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: "users",
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	client, err := NewClient(params)
	require.NoError(t, err)
	tx, err := client.Begin(ctx)
	require.NoError(t, err)
	f := client.Close(tx)
	require.NotNil(t, f)
}

func TestTransaction(t *testing.T) {
	t.Skip() // 確認したい時のみコメントアウト

	setEnv()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	params := &Params{
		Socket:   "tcp",
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: "users",
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	client, err := NewClient(params)
	require.NoError(t, err)
	t.Run("success", func(t *testing.T) {
		err := client.Transaction(ctx, func(tx *gorm.DB) error {
			return nil
		})
		require.NoError(t, err)
	})
	t.Run("failure", func(t *testing.T) {
		err := client.Transaction(ctx, func(tx *gorm.DB) error {
			return assert.AnError
		})
		require.Error(t, err)
	})
}

func TestRetryable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "dead lock",
			err:    &dmysql.MySQLError{Number: 1213},
			expect: true,
		},
		{
			name:   "non retryable error",
			err:    assert.AnError,
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := Retryable(tt.err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestNewDSN(t *testing.T) {
	t.Parallel()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	tests := []struct {
		name    string
		params  *Params
		options *options
		expect  string
	}{
		{
			name: "success",
			params: &Params{
				Socket:   "tcp",
				Host:     "127.0.0.1",
				Port:     "3306",
				Database: "test",
				Username: "root",
				Password: "12345678",
			},
			options: &options{
				now:                  time.Now,
				location:             time.UTC,
				charset:              "utf8mb4",
				collation:            "utf8mb4_general_ci",
				enabledTLS:           false,
				allowNativePasswords: true,
				maxAllowedPacket:     4194304, // 4MiB
			},
			expect: "root:12345678@tcp(127.0.0.1:3306)/test?collation=utf8mb4_general_ci&parseTime=true&maxAllowedPacket=4194304&charset=utf8mb4",
		},
		{
			name: "success with options",
			params: &Params{
				Socket:   "tcp",
				Host:     "127.0.0.1",
				Port:     "3306",
				Database: "test",
				Username: "root",
				Password: "12345678",
			},
			options: &options{
				now:                  time.Now,
				location:             jst,
				charset:              "utf8mb4",
				collation:            "utf8mb4_0900_ai_ci",
				enabledTLS:           true,
				allowNativePasswords: false,
				maxAllowedPacket:     8388608, // 8MiB
			},
			expect: "root:12345678@tcp(127.0.0.1:3306)/test?allowNativePasswords=false&collation=utf8mb4_0900_ai_ci&loc=Asia%2FTokyo&parseTime=true&tls=true&maxAllowedPacket=8388608&charset=utf8mb4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, newDSN(tt.params, tt.options))
		})
	}
}

func setEnv() {
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "127.0.0.1")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "3306")
	}
	if os.Getenv("DB_USERNAME") == "" {
		os.Setenv("DB_USERNAME", "root")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "")
	}
}
