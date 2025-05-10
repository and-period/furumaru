package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	t.Parallel()
	opts := []Option{
		WithLogLevel("debug"),
		WithOutput(""),
		WithSentryServerName("server"),
		WithSentryEnvironment("test"),
		WithSentryLevel("debug"),
		WithSentryFlushTimeout(10 * time.Millisecond),
	}
	logger, err := NewLogger(opts...)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestLogger_GetLogLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		level  string
		expect zapcore.Level
	}{
		{
			name:   "debug",
			level:  "debug",
			expect: zapcore.DebugLevel,
		},
		{
			name:   "info",
			level:  "info",
			expect: zapcore.InfoLevel,
		},
		{
			name:   "warn",
			level:  "warn",
			expect: zapcore.WarnLevel,
		},
		{
			name:   "error",
			level:  "error",
			expect: zapcore.ErrorLevel,
		},
		{
			name:   "default",
			level:  "",
			expect: zapcore.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, getLogLevel(tt.level))
		})
	}
}
