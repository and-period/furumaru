package mysql

import (
	"database/sql/driver"
	"testing"
	"time"

	dmysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestNewTiDBDSN(t *testing.T) {
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
				dialTimeout:          10 * time.Second,
				readTimeout:          30 * time.Second,
				writeTimeout:         30 * time.Second,
			},
			expect: "root:12345678@tcp(127.0.0.1:3306)/test?interpolateParams=true&parseTime=true&readTimeout=30s&timeout=10s&tls=tidb&writeTimeout=30s&maxAllowedPacket=4194304",
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
				dialTimeout:          5 * time.Second,
				readTimeout:          15 * time.Second,
				writeTimeout:         15 * time.Second,
			},
			expect: "root:12345678@tcp(127.0.0.1:3306)/test?interpolateParams=true&loc=Asia%2FTokyo&parseTime=true&readTimeout=15s&timeout=5s&tls=tidb&writeTimeout=15s&maxAllowedPacket=8388608",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, newTiDBDSN(tt.params, tt.options))
		})
	}
}

func TestIsRetryable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "bad connection",
			err:    driver.ErrBadConn,
			expect: true,
		},
		{
			name:   "deadlock",
			err:    &dmysql.MySQLError{Number: 1213, Message: "Deadlock found"},
			expect: true,
		},
		{
			name:   "tiproxy no available instances",
			err:    &dmysql.MySQLError{Number: 1105, Message: "No available TiDB instances, please make sure TiDB is available"},
			expect: true,
		},
		{
			name:   "tiproxy fails to connect",
			err:    &dmysql.MySQLError{Number: 1105, Message: "TiProxy fails to connect to TiDB, please make sure TiDB is available"},
			expect: true,
		},
		{
			name:   "non-retryable mysql error 1105",
			err:    &dmysql.MySQLError{Number: 1105, Message: "other error"},
			expect: false,
		},
		{
			name:   "non-retryable error",
			err:    assert.AnError,
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, IsRetryable(tt.err))
		})
	}
}
