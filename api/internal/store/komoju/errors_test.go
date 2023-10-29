package komoju

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
