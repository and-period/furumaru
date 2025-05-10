package komoju

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect ErrCode
	}{
		{
			name: "komoju error",
			err: &Error{
				Method: "GET",
				Route:  "/route",
				Code:   ErrCodeUnprocessableEntity,
			},
			expect: ErrCodeUnprocessableEntity,
		},
		{
			name:   "unkonown error",
			err:    assert.AnError,
			expect: "",
		},
		{
			name:   "nil",
			err:    nil,
			expect: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewErrCode(tt.err))
		})
	}
}

func TestIsSessionFailed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name: "client error",
			err: &Error{
				Method: "GET",
				Route:  "/route",
				Status: 422,
				Code:   ErrCodeUnprocessableEntity,
			},
			expect: true,
		},
		{
			name: "other komoju error",
			err: &Error{
				Method: "GET",
				Route:  "/route",
				Status: 502,
				Code:   ErrCodeInvalidAccount,
			},
			expect: false,
		},
		{
			name:   "unknown error",
			err:    assert.AnError,
			expect: false,
		},
		{
			name:   "nil",
			err:    nil,
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, IsSessionFailed(tt.err))
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	err := &Error{
		Method:  http.MethodPost,
		Route:   "/api/v1/sessions",
		Status:  400,
		Code:    ErrCodeBadRequest,
		Message: "some error",
	}
	expect := "komoju: method=POST, route=/api/v1/sessions, status=400, code=bad_request, message=some error"
	assert.Equal(t, expect, err.Error())
}
