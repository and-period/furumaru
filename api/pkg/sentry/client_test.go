package sentry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestClient(t *testing.T) {
	t.Parallel()
	t.Run("with valid dsn", func(t *testing.T) {
		client, err := NewClient(WithDSN("http://dummy@sentry.io/123"))
		assert.NoError(t, err)
		assert.NotNil(t, client)

		ctx := t.Context()
		client.ReportError(ctx, assert.AnError)
		client.ReportPanic(ctx, assert.AnError)
		client.ReportMessage(ctx, "some message")
		client.Flush(10 * time.Second)
	})
	t.Run("with invalid dsn", func(t *testing.T) {
		client, err := NewClient(WithDSN("invalid-dsn"))
		assert.Error(t, err)
		assert.Nil(t, client)
	})
	t.Run("without dsn", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestZapHookFn(t *testing.T) {
	t.Parallel()
	client := NewFixedMockClient()
	logger, err := zap.NewDevelopment(zap.Hooks(NewZapHookFn(client)))
	require.NoError(t, err)
	logger.Warn("warn")
	logger.Error("error")
}
