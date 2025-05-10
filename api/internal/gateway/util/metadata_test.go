package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ctx    *gin.Context
		expect string
		isErr  bool
	}{
		{
			name: "success",
			ctx: func() *gin.Context {
				gin.SetMode(gin.TestMode)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = &http.Request{Header: http.Header{}}
				c.Request.Header.Set("Authorization", "Bearer xxxxxx")
				return c
			}(),
			expect: "xxxxxx",
			isErr:  false,
		},
		{
			name: "not exists authorization header",
			ctx: func() *gin.Context {
				gin.SetMode(gin.TestMode)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = &http.Request{Header: http.Header{}}
				return c
			}(),
			expect: "",
			isErr:  true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			actual, err := GetAuthToken(tt.ctx)
			assert.Equal(t, tt.isErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
