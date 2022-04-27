package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGinMiddleware(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		options []Option
		isErr   bool
	}{
		{
			name: "success",
			options: []Option{
				WithLogLevel("debug"),
				WithOutput(""),
			},
			isErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger, err := NewGinMiddleware(tt.options...)
			if tt.isErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}
