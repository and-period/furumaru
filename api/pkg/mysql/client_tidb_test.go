package mysql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
				logger:               zap.NewNop(),
				now:                  time.Now,
				location:             time.UTC,
				charset:              "utf8mb4",
				collation:            "utf8mb4_general_ci",
				enabledTLS:           false,
				allowNativePasswords: true,
				maxAllowedPacket:     4194304, // 4MiB
			},
			expect: "root:12345678@tcp(127.0.0.1:3306)/test?parseTime=true&tls=tidb&maxAllowedPacket=4194304",
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
				logger:               zap.NewNop(),
				now:                  time.Now,
				location:             jst,
				charset:              "utf8mb4",
				collation:            "utf8mb4_0900_ai_ci",
				enabledTLS:           true,
				allowNativePasswords: false,
				maxAllowedPacket:     8388608, // 8MiB
			},
			expect: "root:12345678@tcp(127.0.0.1:3306)/test?loc=Asia%2FTokyo&parseTime=true&tls=tidb&maxAllowedPacket=8388608",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, newTiDBDSN(tt.params, tt.options))
		})
	}
}
