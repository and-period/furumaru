package util

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetParam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx *gin.Context)
		param  string
		expect string
	}{
		{
			name: "success",
			setup: func(ctx *gin.Context) {
				ctx.Params = gin.Params{{Key: "id", Value: "hoge"}}
			},
			param:  "id",
			expect: "hoge",
		},
		{
			name:   "empty",
			setup:  func(ctx *gin.Context) {},
			param:  "id",
			expect: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			tt.setup(ctx)
			assert.Equal(t, tt.expect, GetParam(ctx, tt.param))
		})
	}
}

func TestGetParamInt64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx *gin.Context)
		param  string
		expect int64
		hasErr bool
	}{
		{
			name: "success",
			setup: func(ctx *gin.Context) {
				ctx.Params = gin.Params{{Key: "id", Value: "1"}}
			},
			param:  "id",
			expect: 1,
			hasErr: false,
		},
		{
			name:   "empty",
			setup:  func(ctx *gin.Context) {},
			param:  "id",
			expect: 0,
			hasErr: true,
		},
		{
			name: "invalid param",
			setup: func(ctx *gin.Context) {
				ctx.Params = gin.Params{{Key: "id", Value: "hoge"}}
			},
			param:  "id",
			expect: 0,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			tt.setup(ctx)
			actual, err := GetParamInt64(ctx, tt.param)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestGetQuery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx *gin.Context)
		query  string
		value  string
		expect string
	}{
		{
			name: "success",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?id=hoge")
			},
			query:  "id",
			value:  "",
			expect: "hoge",
		},
		{
			name:   "empty",
			setup:  func(ctx *gin.Context) {},
			query:  "id",
			value:  "",
			expect: "",
		},
		{
			name:   "empty with default value",
			setup:  func(ctx *gin.Context) {},
			query:  "id",
			value:  "hoge",
			expect: "hoge",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{URL: &url.URL{}}
			tt.setup(ctx)
			assert.Equal(t, tt.expect, GetQuery(ctx, tt.query, tt.value))
		})
	}
}

func TestGetQueryInt32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx *gin.Context)
		query  string
		value  int32
		expect int32
		hasErr bool
	}{
		{
			name: "success",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?id=10")
			},
			query:  "id",
			value:  0,
			expect: 10,
			hasErr: false,
		},
		{
			name:   "empty",
			setup:  func(ctx *gin.Context) {},
			query:  "id",
			value:  0,
			expect: 0,
			hasErr: false,
		},
		{
			name:   "empty with default value",
			setup:  func(ctx *gin.Context) {},
			query:  "id",
			value:  10,
			expect: 10,
			hasErr: false,
		},
		{
			name: "invalid query",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?id=hoge")
			},
			query:  "id",
			value:  0,
			expect: 0,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{URL: &url.URL{}}
			tt.setup(ctx)
			actual, err := GetQueryInt32(ctx, tt.query, tt.value)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestGetQueryInt64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx *gin.Context)
		query  string
		value  int64
		expect int64
		hasErr bool
	}{
		{
			name: "success",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?id=10")
			},
			query:  "id",
			value:  0,
			expect: 10,
			hasErr: false,
		},
		{
			name:   "empty",
			setup:  func(ctx *gin.Context) {},
			query:  "id",
			value:  0,
			expect: 0,
			hasErr: false,
		},
		{
			name:   "empty with default value",
			setup:  func(ctx *gin.Context) {},
			query:  "id",
			value:  10,
			expect: 10,
			hasErr: false,
		},
		{
			name: "invalid query",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?id=hoge")
			},
			query:  "id",
			value:  0,
			expect: 0,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{URL: &url.URL{}}
			tt.setup(ctx)
			actual, err := GetQueryInt64(ctx, tt.query, tt.value)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestGetQueryStrings(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx *gin.Context)
		query  string
		expect []string
	}{
		{
			name: "success",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?id=foo,baz,bar")
			},
			query:  "id",
			expect: []string{"foo", "baz", "bar"},
		},
		{
			name:   "empty",
			setup:  func(ctx *gin.Context) {},
			query:  "id",
			expect: []string{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{URL: &url.URL{}}
			tt.setup(ctx)
			assert.Equal(t, tt.expect, GetQueryStrings(ctx, tt.query))
		})
	}
}

func TestGetQueryInt32s(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx *gin.Context)
		query  string
		expect []int32
		hasErr bool
	}{
		{
			name: "success",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?id=1,10,100")
			},
			query:  "id",
			expect: []int32{1, 10, 100},
			hasErr: false,
		},
		{
			name:   "empty",
			setup:  func(ctx *gin.Context) {},
			query:  "id",
			expect: []int32{},
			hasErr: false,
		},
		{
			name: "invalid query",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?id=foo,bar")
			},
			query:  "id",
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{URL: &url.URL{}}
			tt.setup(ctx)
			actual, err := GetQueryInt32s(ctx, tt.query)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestGetOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx *gin.Context)
		query  string
		expect []*Order
	}{
		{
			name: "success",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?orders=foo,-baz,bar")
			},
			expect: []*Order{
				{Key: "foo", Direction: OrderByASC},
				{Key: "baz", Direction: OrderByDesc},
				{Key: "bar", Direction: OrderByASC},
			},
		},
		{
			name:   "empty",
			setup:  func(ctx *gin.Context) {},
			expect: []*Order{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{URL: &url.URL{}}
			tt.setup(ctx)
			assert.Equal(t, tt.expect, GetOrders(ctx))
		})
	}
}
