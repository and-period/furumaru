package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstitutions(t *testing.T) {
	t.Parallel()
	params := map[string]string{"key": "value"}
	expect := map[string]interface{}{"key": "value"}
	assert.Equal(t, expect, NewSubstitutions(params))
}
