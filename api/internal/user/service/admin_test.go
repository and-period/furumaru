package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestMultiGetAdmins(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	adminIDs := []string{
		"administrator-id",
		"coordinator-id",
		"producer-id",
	}
	admins := entity.Admins{
		{
			ID:            "administrator-id",
			Type:          entity.AdminTypeAdministrator,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            "coordinator-id",
			Type:          entity.AdminTypeCoordinator,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            "producer-id",
			Type:          entity.AdminTypeProducer,
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetAdminsInput
		expect    entity.Admins
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().MultiGet(ctx, adminIDs).Return(admins, nil)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: []string{
					"administrator-id",
					"coordinator-id",
					"producer-id",
				},
			},
			expect:    admins,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetAdminsInput{
				AdminIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get admin auths",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().MultiGet(ctx, adminIDs).Return(nil, assert.AnError)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: adminIDs,
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetAdmins(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetAdminDevices(t *testing.T) {
	t.Parallel()

	admins := entity.Admins{
		{
			ID:        "admin-id",
			CognitoID: "username",
			Type:      entity.AdminTypeAdministrator,
			Device:    "instance-id",
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetAdminDevicesInput
		expect    []string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().MultiGet(ctx, []string{"admin-id"}, "device").Return(admins, nil)
			},
			input: &user.MultiGetAdminDevicesInput{
				AdminIDs: []string{"admin-id"},
			},
			expect:    []string{"instance-id"},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetAdminDevicesInput{
				AdminIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get admin auths",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().MultiGet(ctx, []string{"admin-id"}, "device").Return(nil, assert.AnError)
			},
			input: &user.MultiGetAdminDevicesInput{
				AdminIDs: []string{"admin-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetAdminDevices(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetAdmin(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	admin := &entity.Admin{
		ID:            "admin-id",
		Type:          entity.AdminTypeAdministrator,
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
			name: "success to administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect: func() *entity.Admin {
				admin := *admin
				admin.Type = entity.AdminTypeAdministrator
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
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(nil, assert.AnError)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestForgotAdminPassword(t *testing.T) {
	t.Parallel()

	admin := &entity.Admin{CognitoID: "cognito-id"}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.ForgotAdminPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().GetByEmail(ctx, "test-admin@and-period.jp", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ForgotPassword(ctx, "cognito-id").Return(nil)
			},
			input: &user.ForgotAdminPasswordInput{
				Email: "test-admin@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ForgotAdminPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to forgot password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().GetByEmail(ctx, "test-admin@and-period.jp", "cognito_id").Return(nil, assert.AnError)
			},
			input: &user.ForgotAdminPasswordInput{
				Email: "test-admin@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to forgot password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().GetByEmail(ctx, "test-admin@and-period.jp", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ForgotPassword(ctx, "cognito-id").Return(assert.AnError)
			},
			input: &user.ForgotAdminPasswordInput{
				Email: "test-admin@and-period.jp",
			},
			expectErr: exception.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ForgotAdminPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyAdminPassword(t *testing.T) {
	t.Parallel()

	admin := &entity.Admin{
		CognitoID: "cognito-id",
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    "cognito-id",
		VerifyCode:  "123456",
		NewPassword: "Passw0rd",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.VerifyAdminPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().GetByEmail(ctx, "test-admin@and-period.jp", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmForgotPassword(ctx, params).Return(nil)
			},
			input: &user.VerifyAdminPasswordInput{
				Email:                "test-admin@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.VerifyAdminPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to forgot password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().GetByEmail(ctx, "test-admin@and-period.jp", "cognito_id").Return(nil, assert.AnError)
			},
			input: &user.VerifyAdminPasswordInput{
				Email:                "test-admin@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to forgot password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().GetByEmail(ctx, "test-admin@and-period.jp", "cognito_id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmForgotPassword(ctx, params).Return(assert.AnError)
			},
			input: &user.VerifyAdminPasswordInput{
				Email:                "test-admin@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyAdminPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
