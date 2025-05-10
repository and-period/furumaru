package postalcode

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestClient(t *testing.T) {
	t.Parallel()
	h := &http.Client{}
	cli := NewClient(h, WithLogger(zap.NewNop()), WithMaxRetries(1))
	assert.NotNil(t, cli)
}

func TestError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "not error",
			err:    nil,
			expect: nil,
		},
		{
			name:   "invalid argument",
			err:    &apiError{Status: http.StatusBadRequest},
			expect: ErrInvalidArgument,
		},
		{
			name:   "not found",
			err:    &apiError{Status: http.StatusNotFound},
			expect: ErrNotFound,
		},
		{
			name:   "unavailable",
			err:    &apiError{Status: http.StatusBadGateway},
			expect: ErrUnavailable,
		},
		{
			name:   "deadline exceeded",
			err:    &apiError{Status: http.StatusGatewayTimeout},
			expect: ErrDeadlineExceeded,
		},
		{
			name:   "internal",
			err:    &apiError{Status: http.StatusInternalServerError},
			expect: ErrInternal,
		},
		{
			name:   "canceled",
			err:    context.Canceled,
			expect: ErrCanceled,
		},
		{
			name:   "timeout",
			err:    context.DeadlineExceeded,
			expect: ErrTimeout,
		},
		{
			name:   "unknown",
			err:    errors.New("some error"),
			expect: ErrUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli := &client{logger: zap.NewNop()}
			err := cli.newError(tt.err)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}
