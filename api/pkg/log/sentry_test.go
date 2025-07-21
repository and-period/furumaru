package log

import (
	"testing"
	"time"

	mock_sentry "github.com/and-period/furumaru/api/mock/pkg/sentry"
	"github.com/and-period/furumaru/api/pkg/sentry"
	sentrygo "github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestSentryLogger(t *testing.T) {
	t.Parallel()
	logger, err := NewSentryLogger("")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestSentryOptions(t *testing.T) {
	t.Parallel()
	now := time.Now()
	fields := ZapFields{
		zap.Error(assert.AnError),
		zap.String("string", "str"),
		zap.Bool("string", true),
		zap.Bool("string", false),
		zap.Time("time", now),
		zap.Duration("duration", 10*time.Second),
		zap.Any("struct", struct{}{}),
		zap.Any("sentryValue", sentry.WithTag("key", "value")),
	}
	opts := fields.SentryOptions()
	assert.Len(t, opts, 2)
}

func TestLoggerWithSentry(t *testing.T) {
	t.Parallel()
	const timeout = 10 * time.Millisecond
	tests := []struct {
		name  string
		setup func(mock *mock_sentry.MockClient)
		exec  func(logger *zap.Logger)
	}{
		{
			name: "report error",
			setup: func(mock *mock_sentry.MockClient) {
				mock.EXPECT().Flush(timeout).Return(true)
				mock.EXPECT().ReportError(gomock.Any(), gomock.Any(), gomock.Any())
			},
			exec: func(logger *zap.Logger) {
				logger.Error("some message", zap.Error(assert.AnError))
			},
		},
		{
			name: "report message",
			setup: func(mock *mock_sentry.MockClient) {
				mock.EXPECT().Flush(timeout).Return(true)
				mock.EXPECT().ReportMessage(gomock.Any(), gomock.Any(), gomock.Any())
			},
			exec: func(logger *zap.Logger) {
				logger.Error(
					"some message",
					zap.String("string", "str"),
					zap.Any("sentry", sentry.WithTag("key", "value")),
				)
			},
		},
		{
			name: "with",
			setup: func(mock *mock_sentry.MockClient) {
				mock.EXPECT().Flush(timeout).Return(false)
				mock.EXPECT().ReportMessage(gomock.Any(), gomock.Any(), gomock.Any()).Times(3)
			},
			exec: func(logger *zap.Logger) {
				logger.Error("some message 1")
				logger = logger.With(zap.Int64("num", 64))
				logger.Warn("some message 2")
				logger.Error("some message 3")
			},
		},
		{
			name: "no call",
			setup: func(mock *mock_sentry.MockClient) {
				mock.EXPECT().Flush(timeout).Return(false)
			},
			exec: func(logger *zap.Logger) {
				logger.Info("some message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_sentry.NewMockClient(ctrl)
			core := &sentryCore{
				core:         zap.NewNop().Core(),
				client:       mock,
				level:        zapcore.InfoLevel,
				sentryLevel:  zapcore.WarnLevel,
				flushTimeout: timeout,
			}

			logger := zap.New(core)
			defer logger.Sync()

			tt.setup(mock)
			tt.exec(logger)
		})
	}
}

func TestZapFields_FirstError(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		fs := make(ZapFields, 0, 1)
		ok, err := fs.FirstError()
		assert.False(t, ok)
		assert.NoError(t, err)
	})
	t.Run("exist", func(t *testing.T) {
		t.Parallel()
		fs := make(ZapFields, 0, 1)
		fs = append(fs, zap.String("hoge", "fuga"), zap.Error(assert.AnError))
		ok, res := fs.FirstError()
		require.True(t, ok)
		require.ErrorIs(t, res, assert.AnError)
	})
}

func TestLogger_GetSentryLovel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		level  zapcore.Level
		expect sentrygo.Level
	}{
		{
			name:   "debug",
			level:  zapcore.DebugLevel,
			expect: sentrygo.LevelDebug,
		},
		{
			name:   "info",
			level:  zapcore.InfoLevel,
			expect: sentrygo.LevelInfo,
		},
		{
			name:   "warn",
			level:  zapcore.WarnLevel,
			expect: sentrygo.LevelWarning,
		},
		{
			name:   "error",
			level:  zapcore.ErrorLevel,
			expect: sentrygo.LevelError,
		},
		{
			name:   "default",
			level:  zapcore.InvalidLevel,
			expect: sentrygo.LevelWarning,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, getSentryLevel(tt.level))
		})
	}
}
