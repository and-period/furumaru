package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignInAdmin(t *testing.T) {
	t.Parallel()

	result := &cognito.AuthResult{
		IDToken:      "id-token",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}
	admin := &entity.Admin{
		ID:   "admin-id",
		Type: entity.AdminTypeAdministrator,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.SignInAdminInput
		expect    *entity.AdminAuth
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username").Return(admin, nil)
				mocks.db.Admin.EXPECT().UpdateSignInAt(ctx, "admin-id").Return(nil)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect: &entity.AdminAuth{
				AdminID:      "admin-id",
				Type:         entity.AdminTypeAdministrator,
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.SignInAdminInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to sign in",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "username", "password").Return(nil, assert.AnError)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username").Return(nil, assert.AnError)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username").Return(admin, nil)
				mocks.db.Admin.EXPECT().UpdateSignInAt(ctx, "admin-id").Return(assert.AnError)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.SignInAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestSignOutAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.SignOutAdminInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignOut(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil)
			},
			input: &user.SignOutAdminInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.SignOutAdminInput{
				AccessToken: "",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to sign out",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignOut(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(assert.AnError)
			},
			input: &user.SignOutAdminInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.SignOutAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestGetAdminAuth(t *testing.T) {
	t.Parallel()

	admin := &entity.Admin{
		ID:     "admin-id",
		Type:   entity.AdminTypeAdministrator,
		Status: entity.AdminStatusActivated,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetAdminAuthInput
		expect    *entity.AdminAuth
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username").Return(admin, nil)
			},
			input: &user.GetAdminAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &entity.AdminAuth{
				AdminID:      "admin-id",
				Type:         entity.AdminTypeAdministrator,
				AccessToken:  "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				RefreshToken: "",
				ExpiresIn:    0,
			},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.GetAdminAuthInput{
				AccessToken: "",
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", assert.AnError)
			},
			input: &user.GetAdminAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username").Return(nil, assert.AnError)
			},
			input: &user.GetAdminAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetAdminAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestRefreshAdminToken(t *testing.T) {
	t.Parallel()

	result := &cognito.AuthResult{
		IDToken:      "id-token",
		AccessToken:  "access-token",
		RefreshToken: "",
		ExpiresIn:    3600,
	}
	admin := &entity.Admin{
		ID:   "admin-id",
		Type: entity.AdminTypeAdministrator,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.RefreshAdminTokenInput
		expect    *entity.AdminAuth
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username").Return(admin, nil)
				mocks.db.Admin.EXPECT().UpdateSignInAt(ctx, "admin-id").Return(nil)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &entity.AdminAuth{
				AdminID:      "admin-id",
				Type:         entity.AdminTypeAdministrator,
				AccessToken:  "access-token",
				RefreshToken: "",
				ExpiresIn:    3600,
			},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "",
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to sign in",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil, assert.AnError)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username").Return(nil, assert.AnError)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username").Return(admin, nil)
				mocks.db.Admin.EXPECT().UpdateSignInAt(ctx, "admin-id").Return(assert.AnError)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.RefreshAdminToken(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestRegisterAdminDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.RegisterAdminDeviceInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().UpdateDevice(ctx, "admin-id", "device").Return(nil)
			},
			input: &user.RegisterAdminDeviceInput{
				AdminID: "admin-id",
				Device:  "device",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.RegisterAdminDeviceInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update device",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Admin.EXPECT().UpdateDevice(ctx, "admin-id", "device").Return(assert.AnError)
			},
			input: &user.RegisterAdminDeviceInput{
				AdminID: "admin-id",
				Device:  "device",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.RegisterAdminDevice(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateAdminEmail(t *testing.T) {
	t.Parallel()

	admin := &entity.Admin{
		Email: "test-admin@and-period.jp",
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: "access-token",
		Username:    "username",
		OldEmail:    "test-admin@and-period.jp",
		NewEmail:    "test-other@and-period.jp",
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
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username", "email").Return(admin, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(nil)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
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
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username", "email").Return(nil, assert.AnError)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "does not need to be changed",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username", "email").Return(admin, nil)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-admin@and-period.jp",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username", "email").Return(admin, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(assert.AnError)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateAdminEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyAdminEmail(t *testing.T) {
	t.Parallel()

	admin := &entity.Admin{
		ID:   "admin-id",
		Type: entity.AdminTypeAdministrator,
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: "access-token",
		Username:    "username",
		VerifyCode:  "123456",
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
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username", "id", "role").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test-admin@and-period.jp", nil)
				mocks.db.Admin.EXPECT().UpdateEmail(ctx, "admin-id", "test-admin@and-period.jp").Return(nil)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
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
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", assert.AnError)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username", "id", "role").Return(nil, assert.AnError)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to confirm change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.Admin.EXPECT().GetByCognitoID(ctx, "username", "id", "role").Return(admin, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("", assert.AnError)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyAdminEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateAdminPassword(t *testing.T) {
	t.Parallel()

	params := &cognito.ChangePasswordParams{
		AccessToken: "access-token",
		OldPassword: "12345678",
		NewPassword: "!Qaz2wsx",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateAdminPasswordInput
		expectErr error
	}{
		{
			name: "success to administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().ChangePassword(ctx, params).Return(nil)
			},
			input: &user.UpdateAdminPasswordInput{
				AccessToken:          "access-token",
				OldPassword:          "12345678",
				NewPassword:          "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
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
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().ChangePassword(ctx, params).Return(assert.AnError)
			},
			input: &user.UpdateAdminPasswordInput{
				AccessToken:          "access-token",
				OldPassword:          "12345678",
				NewPassword:          "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateAdminPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestListAdminAuthProviders(t *testing.T) {
	t.Parallel()

	now := time.Now()
	params := &database.ListAdminAuthProvidersParams{
		AdminID: "admin-id",
	}
	providers := entity.AdminAuthProviders{
		{
			AdminID:      "admin-id",
			ProviderType: entity.AdminAuthProviderTypeGoogle,
			AccountID:    "account-id",
			Email:        "account-email",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.ListAdminAuthProvidersInput
		expect    entity.AdminAuthProviders
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuthProvider.EXPECT().List(ctx, params).Return(providers, nil)
			},
			input: &user.ListAdminAuthProvidersInput{
				AdminID: "admin-id",
			},
			expect:    providers,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ListAdminAuthProvidersInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuthProvider.EXPECT().List(ctx, params).Return(nil, assert.AnError)
			},
			input: &user.ListAdminAuthProvidersInput{
				AdminID: "admin-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListAdminAuthProviders(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestInitialGoogleAdminAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.InitialGoogleAdminAuthInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mockInitialAdminAuth(mocks, t)
			},
			input: &user.InitialGoogleAdminAuthInput{
				AdminID:     "admin-id",
				State:       "state",
				RedirectURI: "http://example.com/auth/google/callback",
			},
			expect:    "http://example.com/auth/google/callback",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.InitialGoogleAdminAuthInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.InitialGoogleAdminAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestConnectGoogleAdminAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.ConnectGoogleAdminAuthInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mockConnectAdminAuth(mocks, t)
			},
			input: &user.ConnectGoogleAdminAuthInput{
				AdminID: "admin-id",
				Code:    "code",
				Nonce:   "nonce",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ConnectGoogleAdminAuthInput{},
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ConnectGoogleAdminAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func mockInitialAdminAuth(m *mocks, t *testing.T) {
	params := &cognito.GenerateAuthURLParams{
		State:        "state",
		Nonce:        "nonce",
		ProviderType: cognito.ProviderTypeGoogle,
		RedirectURI:  "http://example.com/auth/google/callback",
	}

	m.db.AdminAuthProvider.EXPECT().Get(gomock.Any(), "admin-id", entity.AdminAuthProviderTypeGoogle).Return(nil, database.ErrNotFound)
	m.cache.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
	m.adminAuth.EXPECT().
		GenerateAuthURL(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, p *cognito.GenerateAuthURLParams) (string, error) {
			p.Nonce = "nonce"
			assert.Equal(t, params, p)
			return "http://example.com/auth/google/callback", nil
		})
}

func TestInitialAdminAuth(t *testing.T) {
	t.Parallel()

	params := &cognito.GenerateAuthURLParams{
		State:        "state",
		Nonce:        "nonce",
		ProviderType: cognito.ProviderTypeGoogle,
		RedirectURI:  "http://example.com/auth/google/callback",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *initialAdminAuthParams
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuthProvider.EXPECT().Get(ctx, "admin-id", entity.AdminAuthProviderTypeGoogle).Return(nil, database.ErrNotFound)
				mocks.cache.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
				mocks.adminAuth.EXPECT().
					GenerateAuthURL(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, p *cognito.GenerateAuthURLParams) (string, error) {
						p.Nonce = "nonce"
						assert.Equal(t, params, p)
						return "http://example.com/auth/google", nil
					})
			},
			input: &initialAdminAuthParams{
				adminID:      "admin-id",
				state:        "state",
				providerType: entity.AdminAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expect:    "http://example.com/auth/google",
			expectErr: nil,
		},
		{
			name: "failed to get auth provider",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuthProvider.EXPECT().Get(ctx, "admin-id", entity.AdminAuthProviderTypeGoogle).Return(nil, assert.AnError)
			},
			input: &initialAdminAuthParams{
				adminID:      "admin-id",
				state:        "state",
				providerType: entity.AdminAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "alraedy connected",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuthProvider.EXPECT().Get(ctx, "admin-id", entity.AdminAuthProviderTypeGoogle).Return(&entity.AdminAuthProvider{}, nil)
			},
			input: &initialAdminAuthParams{
				adminID:      "admin-id",
				state:        "state",
				providerType: entity.AdminAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expect:    "",
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to insert cache",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuthProvider.EXPECT().Get(ctx, "admin-id", entity.AdminAuthProviderTypeGoogle).Return(nil, database.ErrNotFound)
				mocks.cache.EXPECT().Insert(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &initialAdminAuthParams{
				adminID:      "admin-id",
				state:        "state",
				providerType: entity.AdminAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to generate auth url",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.AdminAuthProvider.EXPECT().Get(ctx, "admin-id", entity.AdminAuthProviderTypeGoogle).Return(nil, database.ErrNotFound)
				mocks.cache.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
				mocks.adminAuth.EXPECT().GenerateAuthURL(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &initialAdminAuthParams{
				adminID:      "admin-id",
				state:        "state",
				providerType: entity.AdminAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.initialAdminAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func mockConnectAdminAuth(m *mocks, _ *testing.T) {
	admin := &entity.Admin{
		ID:        "admin-id",
		CognitoID: "cognito-id",
	}
	tokenParams := &cognito.GetAccessTokenParams{
		Code:        "code",
		RedirectURI: "http://example.com/auth/google/callback",
	}
	token := &cognito.AuthResult{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
		ExpiresIn:    3600,
	}
	authUser := &cognito.AuthUser{
		Username:    "google_username",
		Email:       "test@example.com",
		PhoneNumber: "",
	}
	linkParams := &cognito.LinkProviderParams{
		Username:     "cognito-id",
		ProviderType: cognito.ProviderTypeGoogle,
		AccountID:    "username",
	}
	provider := &entity.AdminAuthProvider{
		AdminID:      "admin-id",
		ProviderType: entity.AdminAuthProviderTypeGoogle,
		AccountID:    "username",
		Email:        "test@example.com",
	}

	m.cache.EXPECT().
		Get(gomock.Any(), &entity.AdminAuthEvent{AdminID: "admin-id"}).
		DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
			event.ProviderType = entity.AdminAuthProviderTypeGoogle
			event.Nonce = "nonce"
			return nil
		})
	m.db.Admin.EXPECT().Get(gomock.Any(), "admin-id").Return(admin, nil)
	m.adminAuth.EXPECT().GetAccessToken(gomock.Any(), tokenParams).Return(token, nil)
	m.adminAuth.EXPECT().GetUser(gomock.Any(), "access-token").Return(authUser, nil)
	m.adminAuth.EXPECT().DeleteUser(gomock.Any(), "google_username").Return(nil)
	m.adminAuth.EXPECT().LinkProvider(gomock.Any(), linkParams).Return(nil)
	m.db.AdminAuthProvider.EXPECT().Upsert(gomock.Any(), provider).Return(nil)
}

func TestConnectAdminAuth(t *testing.T) {
	t.Parallel()

	admin := &entity.Admin{
		ID:        "admin-id",
		CognitoID: "cognito-id",
	}
	tokenParams := &cognito.GetAccessTokenParams{
		Code:        "code",
		RedirectURI: "http://example.com/auth/google/callback",
	}
	token := &cognito.AuthResult{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
		ExpiresIn:    3600,
	}
	authUser := &cognito.AuthUser{
		Username:    "google_username",
		Email:       "test@example.com",
		PhoneNumber: "",
	}
	linkParams := &cognito.LinkProviderParams{
		Username:     "cognito-id",
		ProviderType: cognito.ProviderTypeGoogle,
		AccountID:    "username",
	}
	provider := &entity.AdminAuthProvider{
		AdminID:      "admin-id",
		ProviderType: entity.AdminAuthProviderTypeGoogle,
		AccountID:    "username",
		Email:        "test@example.com",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *connectAdminAuthParams
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(authUser, nil)
				mocks.adminAuth.EXPECT().DeleteUser(ctx, "google_username").Return(nil)
				mocks.adminAuth.EXPECT().LinkProvider(ctx, linkParams).Return(nil)
				mocks.db.AdminAuthProvider.EXPECT().Upsert(ctx, provider).Return(nil)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: nil,
		},
		{
			name: "failed to get event",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					Return(assert.AnError)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "invalid nonce",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "invalid-token"
						return nil
					})
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrForbidden,
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(nil, assert.AnError)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get external access token",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(nil, assert.AnError)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get external user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(nil, assert.AnError)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to parse external username",
			setup: func(ctx context.Context, mocks *mocks) {
				authUser := &cognito.AuthUser{
					Username:    "invalid-username",
					Email:       "test@example.com",
					PhoneNumber: "",
				}
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(authUser, nil)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrAlreadyExists,
		},
		{
			name: "failed to delete external user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(authUser, nil)
				mocks.adminAuth.EXPECT().DeleteUser(ctx, "google_username").Return(assert.AnError)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to link provider",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(authUser, nil)
				mocks.adminAuth.EXPECT().DeleteUser(ctx, "google_username").Return(nil)
				mocks.adminAuth.EXPECT().LinkProvider(ctx, linkParams).Return(assert.AnError)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upsert provider",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.AdminAuthEvent{AdminID: "admin-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.AdminAuthEvent) error {
						event.ProviderType = entity.AdminAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.db.Admin.EXPECT().Get(ctx, "admin-id").Return(admin, nil)
				mocks.adminAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.adminAuth.EXPECT().GetUser(ctx, "access-token").Return(authUser, nil)
				mocks.adminAuth.EXPECT().DeleteUser(ctx, "google_username").Return(nil)
				mocks.adminAuth.EXPECT().LinkProvider(ctx, linkParams).Return(nil)
				mocks.db.AdminAuthProvider.EXPECT().Upsert(ctx, provider).Return(assert.AnError)
			},
			input: &connectAdminAuthParams{
				adminID:     "admin-id",
				code:        "code",
				nonce:       "nonce",
				redirectURI: "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.connectAdminAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
