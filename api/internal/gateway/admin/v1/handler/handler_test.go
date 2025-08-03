package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Parallel()
	opts := []Option{
		WithAppName("app-name"),
		WithEnvironment("test"),
		WithSentry(sentry.NewFixedMockClient()),
	}
	h := NewHandler(&Params{}, opts...)
	assert.NotNil(t, h)
}

func TestSetAuth(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{Header: http.Header{}}
	auth := &service.Auth{Auth: response.Auth{
		AdminID: "admin-id",
		Type:    service.AdminTypeAdministrator.Response(),
	}}
	setAuth(ctx, auth)
	assert.Equal(t, "admin-id", getAdminID(ctx))
	assert.Equal(t, service.AdminTypeAdministrator, getAdminType(ctx))
	assert.True(t, currentAdmin(ctx, "admin-id"))
	assert.False(t, currentAdmin(ctx, "other-id"))
}

func TestFilterAccess(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name   string
		role   service.AdminType
		params *filterAccessParams
		expect error
	}{
		{
			name:   "success administrator",
			role:   service.AdminTypeAdministrator,
			params: &filterAccessParams{},
			expect: nil,
		},
		{
			name: "success coordinator",
			role: service.AdminTypeCoordinator,
			params: &filterAccessParams{
				coordinator: func(_ *gin.Context) (bool, error) {
					return true, nil
				},
			},
			expect: nil,
		},
		{
			name:   "success coordinator for no filter",
			role:   service.AdminTypeCoordinator,
			params: &filterAccessParams{},
			expect: nil,
		},
		{
			name: "failed coordinator for failed to execute function",
			role: service.AdminTypeCoordinator,
			params: &filterAccessParams{
				coordinator: func(_ *gin.Context) (bool, error) {
					return false, assert.AnError
				},
			},
			expect: assert.AnError,
		},
		{
			name: "failed coordinator for invalid coordinator",
			role: service.AdminTypeCoordinator,
			params: &filterAccessParams{
				coordinator: func(_ *gin.Context) (bool, error) {
					return false, nil
				},
			},
			expect: exception.ErrForbidden,
		},
		{
			name: "success producer",
			role: service.AdminTypeProducer,
			params: &filterAccessParams{
				producer: func(_ *gin.Context) (bool, error) {
					return true, nil
				},
			},
			expect: nil,
		},
		{
			name:   "success producer for no filter",
			role:   service.AdminTypeProducer,
			params: &filterAccessParams{},
			expect: nil,
		},
		{
			name: "failed producer for failed to execute function",
			role: service.AdminTypeProducer,
			params: &filterAccessParams{
				producer: func(_ *gin.Context) (bool, error) {
					return false, assert.AnError
				},
			},
			expect: assert.AnError,
		},
		{
			name: "failed producer for invalid producer",
			role: service.AdminTypeProducer,
			params: &filterAccessParams{
				producer: func(_ *gin.Context) (bool, error) {
					return false, nil
				},
			},
			expect: exception.ErrForbidden,
		},
		{
			name:   "failed unknown admin role",
			role:   service.AdminTypeUnknown,
			params: &filterAccessParams{},
			expect: exception.ErrForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{Header: http.Header{}}
			auth := &service.Auth{Auth: response.Auth{
				AdminID: "admin-id",
				Type:    int32(tt.role),
			}}
			setAuth(ctx, auth)
			assert.ErrorIs(t, filterAccess(ctx, tt.params), tt.expect)
		})
	}
}
