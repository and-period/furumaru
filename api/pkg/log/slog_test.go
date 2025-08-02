package log

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlog(t *testing.T) {
	t.Parallel()
	closer, err := Start(t.Context())
	if err != nil {
		t.Fatalf("failed to start slog: %v", err)
	}
	t.Cleanup(closer)
	t.Run("標準", func(t *testing.T) {
		t.Parallel()
		slog.Info("test info log")
	})
	t.Run("カスタム Error", func(t *testing.T) {
		t.Parallel()
		slog.Error("test error log", Error(assert.AnError))
	})
	t.Run("カスタム Strings", func(t *testing.T) {
		t.Parallel()
		slog.Info("test strings log", Strings("key", []string{"value1", "value2"}))
	})
	t.Run("カスタム Ints", func(t *testing.T) {
		t.Parallel()
		slog.Info("test ints log", Ints("key", []int{1, 2, 3}))
	})
	t.Run("カスタム Int64s", func(t *testing.T) {
		t.Parallel()
		slog.Info("test int64s log", Int64s("key", []int64{1, 2, 3}))
	})
}

func TestSlogLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		level  string
		expect slog.Level
	}{
		{
			name:   "debug",
			level:  "debug",
			expect: slog.LevelDebug,
		},
		{
			name:   "info",
			level:  "info",
			expect: slog.LevelInfo,
		},
		{
			name:   "warn",
			level:  "warn",
			expect: slog.LevelWarn,
		},
		{
			name:   "error",
			level:  "error",
			expect: slog.LevelError,
		},
		{
			name:   "invalid",
			level:  "invalid",
			expect: slog.LevelInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slogLevel(tt.level)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
