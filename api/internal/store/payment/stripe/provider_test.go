package stripe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProvider(t *testing.T) {
	t.Parallel()
	params := &Params{
		SecretKey: "sk_test_xxx",
	}
	p := NewProvider(params)
	assert.NotNil(t, p)
}
