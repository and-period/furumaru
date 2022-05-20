package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAdmins(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	admins := entity.Admins{
		{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			Role:          entity.AdminRoleAdministrator,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
	params := &database.ListAdminsParams{
		Roles:  []entity.AdminRole{entity.AdminRoleAdministrator},
		Limit:  30,
		Offset: 0,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.ListAdminsInput
		expect    entity.Admins
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().List(ctx, params).Return(admins, nil)
			},
			input: &user.ListAdminsInput{
				Roles:  []int32{int32(entity.AdminRoleAdministrator)},
				Limit:  30,
				Offset: 0,
			},
			expect:    admins,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ListAdminsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid role",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.ListAdminsInput{
				Roles:  []int32{int32(entity.AdminRoleUnknown)},
				Limit:  30,
				Offset: 0,
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get admins",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().List(ctx, params).Return(nil, errmock)
			},
			input: &user.ListAdminsInput{
				Roles:  []int32{int32(entity.AdminRoleAdministrator)},
				Limit:  30,
				Offset: 0,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.ListAdmins(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetAdmins(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	admins := entity.Admins{
		{
			ID:            "admin-id",
			CognitoID:     "cognito-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			Role:          entity.AdminRoleAdministrator,
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
				mocks.db.Admin.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(admins, nil)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: []string{"admin-id"},
			},
			expect:    admins,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.MultiGetAdminsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get admins",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().MultiGet(ctx, []string{"admin-id"}).Return(nil, errmock)
			},
			input: &user.MultiGetAdminsInput{
				AdminIDs: []string{"admin-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.MultiGetAdmins(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetAdmin(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	admin := &entity.Admin{
		ID:        "admin-id",
		CognitoID: "cognito-id",
		Email:     "test-admin@and-period.jp",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetAdminInput
		expect    *entity.Admin
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    admin,
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
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.GetAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateAdminInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &user.CreateAdminInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				Role:          int32(entity.AdminRoleAdministrator),
			},
			expectErr: nil,
		},
		{
			name: "success without notify register admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
				mocks.messenger.EXPECT().NotifyRegisterAdmin(gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateAdminInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				Role:          int32(entity.AdminRoleAdministrator),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateAdminInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid role",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.CreateAdminInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-preiod.jp",
				Role:          int32(entity.AdminRoleUnknown),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &user.CreateAdminInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				Role:          int32(entity.AdminRoleAdministrator),
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to create auth admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(errmock)
			},
			input: &user.CreateAdminInput{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				Role:          int32(entity.AdminRoleAdministrator),
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			_, err := service.CreateAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateAdminEmail(t *testing.T) {
	t.Parallel()

	admin := &entity.Admin{
		ID:    "admin-id",
		Email: "test-admin@and-period.jp",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateAdminEmailInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					OldEmail:    "test-admin@and-period.jp",
					NewEmail:    "test-other@and-period.jp",
				}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id", "email").Return(admin, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(nil)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateAdminEmailInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id", "email").Return(nil, errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to change email",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					OldEmail:    "test-admin@and-period.jp",
					NewEmail:    "test-other@and-period.jp",
				}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id", "email").Return(admin, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			err := service.UpdateAdminEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyAdminEmail(t *testing.T) {
	t.Parallel()

	admin := &entity.Admin{
		ID: "admin-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.VerifyAdminEmailInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					VerifyCode:  "123456",
				}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test-user@and-period.jp", nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "admin-id", "test-user@and-period.jp").Return(nil)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.VerifyAdminEmailInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", errmock)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id").Return(nil, errmock)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to confirm change email",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					VerifyCode:  "123456",
				}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("", errmock)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to update email",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					VerifyCode:  "123456",
				}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test-user@and-period.jp", nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "admin-id", "test-user@and-period.jp").Return(errmock)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			err := service.VerifyAdminEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateAdminPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateAdminPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangePasswordParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					OldPassword: "12345678",
					NewPassword: "12345678",
				}
				mocks.adminAuth.EXPECT().ChangePassword(ctx, params).Return(nil)
			},
			input: &user.UpdateAdminPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateAdminPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid argument for password unmatch",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.UpdateAdminPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "123456789",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to change password",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangePasswordParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					OldPassword: "12345678",
					NewPassword: "12345678",
				}
				mocks.adminAuth.EXPECT().ChangePassword(ctx, params).Return(errmock)
			},
			input: &user.UpdateAdminPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			err := service.UpdateAdminPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
