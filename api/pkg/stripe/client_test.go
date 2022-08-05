package stripe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestClient(t *testing.T) {
	t.Parallel()
	cli := NewClient(&Params{},
		WithLogger(zap.NewNop()),
		WithMaxRetries(1),
	)
	assert.NotNil(t, cli)
}

func TestReceiver(t *testing.T) {
	t.Parallel()
	cli := NewReceiver(&Params{},
		WithLogger(zap.NewNop()),
	)
	assert.NotNil(t, cli)
}
