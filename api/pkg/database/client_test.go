package database

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestNewClient(t *testing.T) {
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
			client, err := NewClient(tt.params)
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
	setEnv()
	ctx, cancel := context.WithCancel(context.Background())
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
	setEnv()
	ctx, cancel := context.WithCancel(context.Background())
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
	data, err := client.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		return "data", nil
	})
	require.NoError(t, err)
	require.NotNil(t, data)
}

func TestNewDSN(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		params  *Params
		options *options
		expect  string
	}{
		{
			name: "tcp socket",
			params: &Params{
				Socket:   "tcp",
				Host:     "127.0.0.1",
				Port:     "3306",
				Database: "test",
				Username: "root",
				Password: "12345678",
			},
			options: &options{
				logger:     zap.NewNop(),
				timezone:   "",
				enabledTLS: false,
			},
			expect: "root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Asia%2FTokyo",
		},
		{
			name: "tcp socket with options",
			params: &Params{
				Socket:   "tcp",
				Host:     "127.0.0.1",
				Port:     "3306",
				Database: "test",
				Username: "root",
				Password: "12345678",
			},
			options: &options{
				logger:     zap.NewNop(),
				timezone:   "UTC",
				enabledTLS: true,
			},
			expect: "root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&tls=true&loc=UTC",
		},
		{
			name: "unix socket",
			params: &Params{
				Socket:   "unix",
				Host:     "127.0.0.1",
				Database: "test",
				Username: "root",
				Password: "12345678",
			},
			options: &options{
				logger:     zap.NewNop(),
				timezone:   "",
				enabledTLS: false,
			},
			expect: "root:12345678@unix(127.0.0.1)/test?charset=utf8mb4&parseTime=true",
		},
		{
			name: "unix socket with options",
			params: &Params{
				Socket:   "unix",
				Host:     "127.0.0.1",
				Database: "test",
				Username: "root",
				Password: "12345678",
			},
			options: &options{
				logger:     zap.NewNop(),
				timezone:   "UTC",
				enabledTLS: true,
			},
			expect: "root:12345678@unix(127.0.0.1)/test?charset=utf8mb4&parseTime=true&tls=true",
		},
		{
			name:   "invalid socket type",
			params: &Params{},
			expect: "",
		},
	}

	for _, tt := range tests {
		tt := tt
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
		os.Setenv("DB_PORT", "3326")
	}
	if os.Getenv("DB_USERNAME") == "" {
		os.Setenv("DB_USERNAME", "root")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "12345678")
	}
}
