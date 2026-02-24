package stripe

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrNotSupported(t *testing.T) {
	t.Parallel()
	assert.EqualError(t, ErrNotSupported, "stripe: operation not supported")
}

func TestIsSessionFailed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "nil error",
			err:    nil,
			expect: false,
		},
		{
			name:   "generic error",
			err:    errors.New("some error"),
			expect: true,
		},
		{
			name:   "not supported error",
			err:    ErrNotSupported,
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := isSessionFailed(tt.err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
