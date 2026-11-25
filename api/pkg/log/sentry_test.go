package log

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTransport struct {
	mu     sync.Mutex
	events []*sentry.Event
}

func (t *testTransport) Flush(time.Duration) bool              { return true }
func (t *testTransport) FlushWithContext(context.Context) bool { return true }
func (t *testTransport) Configure(sentry.ClientOptions)        {}
func (t *testTransport) Close()                                {}
func (t *testTransport) SendEvent(event *sentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.events = append(t.events, event)
}

func newTestSentryLogger(t *testing.T, minLevel slog.Level) (*slog.Logger, *testTransport) {
	t.Helper()

	transport := &testTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{Transport: transport})
	require.NoError(t, err)

	hub := sentry.NewHub(client, sentry.NewScope())
	ctx := sentry.SetHubOnContext(t.Context(), hub)

	return slog.New(newSentryHandler(ctx, minLevel)), transport
}

func TestSentryHandler(t *testing.T) {
	t.Parallel()

	t.Run("最小レベル未満は送信しない", func(t *testing.T) {
		t.Parallel()
		logger, transport := newTestSentryLogger(t, slog.LevelError)
		logger.Warn("warn log")
		logger.Info("info log")
		assert.Empty(t, transport.events)
	})

	t.Run("エラーをイベントとして送信する", func(t *testing.T) {
		t.Parallel()
		logger, transport := newTestSentryLogger(t, slog.LevelError)
		cause := errors.New("something failed")
		logger.With(slog.String("service", "gateway")).
			WithGroup("req").
			Error("error log", Error(cause), slog.Int("status", 500))

		require.Len(t, transport.events, 1)
		event := transport.events[0]
		assert.Equal(t, sentry.LevelError, event.Level)
		assert.Equal(t, "error log", event.Message)
		assert.Equal(t, "slog", event.Logger)
		require.NotEmpty(t, event.Exception)
		assert.Equal(t, "something failed", event.Exception[len(event.Exception)-1].Value)

		attrs := event.Contexts[sentryContextKey]
		assert.Equal(t, "gateway", attrs["service"])
		assert.Equal(t, int64(500), attrs["req.status"])
		assert.Equal(t, "something failed", attrs["req.error"])
		assert.Contains(t, attrs["source"], "sentry_test.go")
	})

	t.Run("Fatal レベルを変換する", func(t *testing.T) {
		t.Parallel()
		logger, transport := newTestSentryLogger(t, slog.LevelWarn)
		logger.Log(t.Context(), levelFatal, "fatal log")
		logger.Warn("warn log")
		require.Len(t, transport.events, 2)
		assert.Equal(t, sentry.LevelFatal, transport.events[0].Level)
		assert.Equal(t, sentry.LevelWarning, transport.events[1].Level)
	})
}

func TestSentryLevel(t *testing.T) {
	t.Parallel()

	tests := map[slog.Level]sentry.Level{
		slog.LevelDebug: sentry.LevelDebug,
		slog.LevelInfo:  sentry.LevelInfo,
		slog.LevelWarn:  sentry.LevelWarning,
		slog.LevelError: sentry.LevelError,
		levelFatal:      sentry.LevelFatal,
	}
	for in, expect := range tests {
		assert.Equal(t, expect, sentryLevel(in), in.String())
	}
}
