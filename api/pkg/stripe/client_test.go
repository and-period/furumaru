package stripe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Parallel()
	cli := NewClient(&Params{},
		WithMaxRetries(1),
	)
	assert.NotNil(t, cli)
}

func TestReceiver(t *testing.T) {
	t.Parallel()
	cli := NewReceiver(&Params{})
	assert.NotNil(t, cli)
}
