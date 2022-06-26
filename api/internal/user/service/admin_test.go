package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestGetAdmin(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	admin := &entity.Admin{
		ID:            "admin-id",
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		Email:         "test-admin@and-period.jp",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetAdminInput
		expect    *entity.Admin
		expectErr error
	}{
		{
			name: "success to administartor",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleAdministrator}
				administrator := &entity.Administrator{
					ID:            admin.ID,
					Lastname:      admin.Lastname,
					Firstname:     admin.Firstname,
					LastnameKana:  admin.LastnameKana,
					FirstnameKana: admin.FirstnameKana,
					Email:         admin.Email,
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(administrator, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect: func() *entity.Admin {
				admin := *admin
				admin.Role = entity.AdminRoleAdministrator
				return &admin
			}(),
			expectErr: nil,
		},
		{
			name: "success to coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleCoordinator}
				coordinator := &entity.Coordinator{
					ID:            admin.ID,
					Lastname:      admin.Lastname,
					Firstname:     admin.Firstname,
					LastnameKana:  admin.LastnameKana,
					FirstnameKana: admin.FirstnameKana,
					Email:         admin.Email,
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(coordinator, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect: func() *entity.Admin {
				admin := *admin
				admin.Role = entity.AdminRoleCoordinator
				return &admin
			}(),
			expectErr: nil,
		},
		{
			name: "success to producer",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleProducer}
				producer := &entity.Producer{
					ID:            admin.ID,
					Lastname:      admin.Lastname,
					Firstname:     admin.Firstname,
					LastnameKana:  admin.LastnameKana,
					FirstnameKana: admin.FirstnameKana,
					Email:         admin.Email,
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(producer, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect: func() *entity.Admin {
				admin := *admin
				admin.Role = entity.AdminRoleProducer
				return &admin
			}(),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetAdminInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get admin auth",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleAdministrator}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleCoordinator}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleProducer}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to unknown role",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{Role: entity.AdminRoleUnknown}
				mocks.db.AdminAuth.EXPECT().GetByAdminID(ctx, "admin-id", "role").Return(auth, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
