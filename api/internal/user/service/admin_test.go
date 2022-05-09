package service

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
		input     *GetAdminInput
		expect    *entity.Admin
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
			},
			input: &GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    admin,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &GetAdminInput{},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &GetAdminInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: ErrInternal,
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
		input     *CreateAdminInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(nil)
			},
			input: &CreateAdminInput{
				Email: "test-admin@and-period.jp",
				Role:  entity.AdminRoleAdministrator,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &CreateAdminInput{},
			expectErr: ErrInvalidArgument,
		},
		{
			name:  "invalid role",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &CreateAdminInput{
				Email: "test-admin@and-preiod.jp",
				Role:  entity.AdminRole(-1),
			},
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to create admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &CreateAdminInput{
				Email: "test-admin@and-period.jp",
				Role:  entity.AdminRoleAdministrator,
			},
			expectErr: ErrInternal,
		},
		{
			name: "failed to create auth admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.adminAuth.EXPECT().AdminCreateUser(ctx, gomock.Any()).Return(errmock)
			},
			input: &CreateAdminInput{
				Email: "test-admin@and-period.jp",
				Role:  entity.AdminRoleAdministrator,
			},
			expectErr: ErrInternal,
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
		input     *UpdateAdminEmailInput
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
			input: &UpdateAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &UpdateAdminEmailInput{},
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", errmock)
			},
			input: &UpdateAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id", "email").Return(nil, errmock)
			},
			input: &UpdateAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: ErrInternal,
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
			input: &UpdateAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: ErrInternal,
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
		input     *VerifyAdminEmailInput
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
			input: &VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &VerifyAdminEmailInput{},
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", errmock)
			},
			input: &VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "cognito-id", "id").Return(nil, errmock)
			},
			input: &VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: ErrInternal,
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
			input: &VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: ErrInternal,
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
			input: &VerifyAdminEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: ErrInternal,
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
		input     *UpdateAdminPasswordInput
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
			input: &UpdateAdminPasswordInput{
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
			input:     &UpdateAdminPasswordInput{},
			expectErr: ErrInvalidArgument,
		},
		{
			name:  "invalid argument for password unmatch",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &UpdateAdminPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "123456789",
			},
			expectErr: ErrInvalidArgument,
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
			input: &UpdateAdminPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: ErrInternal,
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
