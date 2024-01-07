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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewErrCode(tt.err))
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	err := &Error{
		Method:  http.MethodPost,
		Route:   "/api/v1/sessions",
		Code:    ErrCodeBadRequest,
		Message: "some error",
	}
	expect := "komoju: method=POST, route=/api/v1/sessions, code=bad_request, message=some error"
	assert.Equal(t, expect, err.Error())
}
