package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAdminRole(t *testing.T) {
	t.Parallel()

	policiesParams := &database.ListAdminPoliciesParams{}
	policies := entity.AdminPolicies{
		{
			ID:          "policy-id",
			Name:        "ポリシー名",
			Description: "ポリシーの説明",
			Priority:    1,
			Path:        "/health",
			Method:      "GET",
			Action:      entity.AdminPolicyActionAllow,
		},
	}
	rolePoliciesParams := &database.ListAdminRolePoliciesParams{}
	rolePolicies := entity.AdminRolePolicies{
		{
			RoleID:   "role-id",
			PolicyID: "policy-id",
		},
	}
	groupRolesParams := &database.ListAdminGroupRolesParams{}
	groupRoles := entity.AdminGroupRoles{
		{
			GroupID: "group-id",
			RoleID:  "role-id",
		},
	}

	tests := []struct {
		name         string
		setup        func(ctx context.Context, mocks *mocks)
		input        *user.GenerateAdminRoleInput
		expectModel  string
		expectPolicy string
		expectErr    error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminPolicy.EXPECT().
					List(gomock.Any(), policiesParams).
					Return(policies, nil)
				mocks.db.AdminRolePolicy.EXPECT().
					List(gomock.Any(), rolePoliciesParams).
					Return(rolePolicies, nil)
				mocks.db.AdminGroupRole.EXPECT().
					List(gomock.Any(), groupRolesParams).
					Return(groupRoles, nil)
			},
			input:       &user.GenerateAdminRoleInput{},
			expectModel: adminRoleModel,
			expectPolicy: `p,policy-id,/health,GET,allow
g,role-id,policy-id
g,group-id,role-id
`,
			expectErr: nil,
		},
		{
			name: "failed to list admin policies",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminPolicy.EXPECT().
					List(gomock.Any(), policiesParams).
					Return(nil, assert.AnError)
				mocks.db.AdminRolePolicy.EXPECT().
					List(gomock.Any(), rolePoliciesParams).
					Return(rolePolicies, nil)
				mocks.db.AdminGroupRole.EXPECT().
					List(gomock.Any(), groupRolesParams).
					Return(groupRoles, nil)
			},
			input:        &user.GenerateAdminRoleInput{},
			expectModel:  "",
			expectPolicy: "",
			expectErr:    exception.ErrInternal,
		},
		{
			name: "failed to list admin role policies",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminPolicy.EXPECT().
					List(gomock.Any(), policiesParams).
					Return(policies, nil)
				mocks.db.AdminRolePolicy.EXPECT().
					List(gomock.Any(), rolePoliciesParams).
					Return(nil, assert.AnError)
				mocks.db.AdminGroupRole.EXPECT().
					List(gomock.Any(), groupRolesParams).
					Return(groupRoles, nil)
			},
			input:        &user.GenerateAdminRoleInput{},
			expectModel:  "",
			expectPolicy: "",
			expectErr:    exception.ErrInternal,
		},
		{
			name: "failed to list admin group roles",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminPolicy.EXPECT().
					List(gomock.Any(), policiesParams).
					Return(policies, nil)
				mocks.db.AdminRolePolicy.EXPECT().
					List(gomock.Any(), rolePoliciesParams).
					Return(rolePolicies, nil)
				mocks.db.AdminGroupRole.EXPECT().
					List(gomock.Any(), groupRolesParams).
					Return(nil, assert.AnError)
			},
			input:        &user.GenerateAdminRoleInput{},
			expectModel:  "",
			expectPolicy: "",
			expectErr:    exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				model, policy, err := service.GenerateAdminRole(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expectModel, model)
				assert.Equal(t, tt.expectPolicy, policy)
			}),
		)
	}
}
