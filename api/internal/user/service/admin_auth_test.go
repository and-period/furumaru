package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
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
	auth := &entity.AdminAuth{
		AdminID:   "admin-id",
		CognitoID: "username",
		Role:      entity.AdminRoleAdministrator,
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
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username").Return(auth, nil)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect: &entity.AdminAuth{
				AdminID:      "admin-id",
				CognitoID:    "username",
				Role:         entity.AdminRoleAdministrator,
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
				mocks.adminAuth.EXPECT().SignIn(ctx, "username", "password").Return(nil, errmock)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", errmock)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username").Return(nil, errmock)
			},
			input: &user.SignInAdminInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
				mocks.adminAuth.EXPECT().SignOut(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(errmock)
			},
			input: &user.SignOutAdminInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expectErr: exception.ErrUnknown,
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

	auth := &entity.AdminAuth{
		AdminID:   "admin-id",
		CognitoID: "username",
		Role:      entity.AdminRoleAdministrator,
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
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username").Return(auth, nil)
			},
			input: &user.GetAdminAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &entity.AdminAuth{
				AdminID:      "admin-id",
				CognitoID:    "username",
				Role:         entity.AdminRoleAdministrator,
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
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", errmock)
			},
			input: &user.GetAdminAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username").Return(auth, errmock)
			},
			input: &user.GetAdminAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
	auth := &entity.AdminAuth{
		AdminID:   "admin-id",
		CognitoID: "username",
		Role:      entity.AdminRoleAdministrator,
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
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username").Return(auth, nil)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &entity.AdminAuth{
				AdminID:      "admin-id",
				CognitoID:    "username",
				Role:         entity.AdminRoleAdministrator,
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
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil, errmock)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", errmock)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username").Return(nil, errmock)
			},
			input: &user.RefreshAdminTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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

func TestUpdateAdminEmail(t *testing.T) {
	t.Parallel()

	auth := &entity.AdminAuth{
		AdminID:   "admin-id",
		CognitoID: "username",
		Role:      entity.AdminRoleAdministrator,
	}
	administrator := &entity.Administrator{
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
			name: "success to administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleAdministrator}
				administrator := &entity.Administrator{Email: "test-admin@and-period.jp"}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(administrator, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(nil)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name: "success to coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleCoordinator}
				coordinator := &entity.Coordinator{Email: "test-admin@and-period.jp"}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(coordinator, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(nil)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name: "success to producer",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleProducer}
				producer := &entity.Producer{Email: "test-admin@and-period.jp"}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(producer, nil)
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
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(nil, errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleAdministrator}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleCoordinator}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.db.Coordinator.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleProducer}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.db.Producer.EXPECT().Get(ctx, "admin-id").Return(nil, errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to unknown role",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleUnknown}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "does not need to be changed",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(administrator, nil)
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
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.db.Administrator.EXPECT().Get(ctx, "admin-id").Return(administrator, nil)
				mocks.adminAuth.EXPECT().ChangeEmail(ctx, params).Return(errmock)
			},
			input: &user.UpdateAdminEmailInput{
				AccessToken: "access-token",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
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

	auth := &entity.AdminAuth{
		AdminID:   "admin-id",
		CognitoID: "username",
		Role:      entity.AdminRoleAdministrator,
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
			name: "success to administrator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleAdministrator}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test-admin@and-period.jp", nil)
				mocks.db.Administrator.EXPECT().UpdateEmail(ctx, "admin-id", "test-admin@and-period.jp").Return(nil)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
				VerifyCode:  "123456",
			},
			expectErr: nil,
		},
		{
			name: "success to coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleCoordinator}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test-admin@and-period.jp", nil)
				mocks.db.Coordinator.EXPECT().UpdateEmail(ctx, "admin-id", "test-admin@and-period.jp").Return(nil)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
				VerifyCode:  "123456",
			},
			expectErr: nil,
		},
		{
			name: "success to producer",
			setup: func(ctx context.Context, mocks *mocks) {
				auth := &entity.AdminAuth{AdminID: "admin-id", Role: entity.AdminRoleProducer}
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test-admin@and-period.jp", nil)
				mocks.db.Producer.EXPECT().UpdateEmail(ctx, "admin-id", "test-admin@and-period.jp").Return(nil)
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
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("", errmock)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(nil, errmock)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to confirm change email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.adminAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.AdminAuth.EXPECT().GetByCognitoID(ctx, "username", "admin_id", "role").Return(auth, nil)
				mocks.adminAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("", errmock)
			},
			input: &user.VerifyAdminEmailInput{
				AccessToken: "access-token",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
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
				mocks.adminAuth.EXPECT().ChangePassword(ctx, params).Return(errmock)
			},
			input: &user.UpdateAdminPasswordInput{
				AccessToken:          "access-token",
				OldPassword:          "12345678",
				NewPassword:          "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
			},
			expectErr: exception.ErrUnknown,
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
