package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler(t *testing.T) {
	t.Parallel()
	opts := []Option{
		WithAppName("app-name"),
		WithEnvironment("test"),
		WithLogger(zap.NewNop()),
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
	setAuth(ctx, "admin-id", service.AdminRoleAdministrator)
	assert.Equal(t, "admin-id", getAdminID(ctx))
	assert.Equal(t, service.AdminRoleAdministrator, getRole(ctx))
	assert.True(t, currentAdmin(ctx, "admin-id"))
	assert.False(t, currentAdmin(ctx, "other-id"))
}

func TestFilterAccess(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name   string
		role   service.AdminRole
		params *filterAccessParams
		expect error
	}{
		{
			name:   "success administrator",
			role:   service.AdminRoleAdministrator,
			params: &filterAccessParams{},
			expect: nil,
		},
		{
			name: "success coordinator",
			role: service.AdminRoleCoordinator,
			params: &filterAccessParams{
				coordinator: func(_ *gin.Context) (bool, error) {
					return true, nil
				},
			},
			expect: nil,
		},
		{
			name:   "success coordinator for no filter",
			role:   service.AdminRoleCoordinator,
			params: &filterAccessParams{},
			expect: nil,
		},
		{
			name: "failed coordinator for failed to execute function",
			role: service.AdminRoleCoordinator,
			params: &filterAccessParams{
				coordinator: func(_ *gin.Context) (bool, error) {
					return false, assert.AnError
				},
			},
			expect: assert.AnError,
		},
		{
			name: "failed coordinator for invalid coordinator",
			role: service.AdminRoleCoordinator,
			params: &filterAccessParams{
				coordinator: func(_ *gin.Context) (bool, error) {
					return false, nil
				},
			},
			expect: exception.ErrForbidden,
		},
		{
			name: "success producer",
			role: service.AdminRoleProducer,
			params: &filterAccessParams{
				producer: func(_ *gin.Context) (bool, error) {
					return true, nil
				},
			},
			expect: nil,
		},
		{
			name:   "success producer for no filter",
			role:   service.AdminRoleProducer,
			params: &filterAccessParams{},
			expect: nil,
		},
		{
			name: "failed producer for failed to execute function",
			role: service.AdminRoleProducer,
			params: &filterAccessParams{
				producer: func(_ *gin.Context) (bool, error) {
					return false, assert.AnError
				},
			},
			expect: assert.AnError,
		},
		{
			name: "failed producer for invalid producer",
			role: service.AdminRoleProducer,
			params: &filterAccessParams{
				producer: func(_ *gin.Context) (bool, error) {
					return false, nil
				},
			},
			expect: exception.ErrForbidden,
		},
		{
			name:   "failed unknown admin role",
			role:   service.AdminRoleUnknown,
			params: &filterAccessParams{},
			expect: exception.ErrForbidden,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{Header: http.Header{}}
			setAuth(ctx, "admin-id", tt.role)
			assert.ErrorIs(t, filterAccess(ctx, tt.params), tt.expect)
		})
	}
}
