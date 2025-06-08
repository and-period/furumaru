package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Parallel()
	t.Run("error", func(t *testing.T) {
		t.Parallel()
		err := &Error{err: errors.New("some error")}
		assert.Equal(t, "some error", err.Error())
	})
	t.Run("unwrap", func(t *testing.T) {
		t.Parallel()
		err := &Error{err: assert.AnError}
		assert.ErrorIs(t, err, assert.AnError)
	})
}
