package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAuth(t *testing.T) {
	t.Parallel()

	auth := &uentity.UserAuth{
		UserID:       "user-id",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetUserAuthInput{AccessToken: tokenmock}
				mocks.user.EXPECT().GetUserAuth(gomock.Any(), in).Return(auth, nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthResponse{
					Auth: &response.Auth{
						UserID:       "user-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					},
				},
			},
		},
		{
			name: "failed to get user auth",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetUserAuthInput{AccessToken: tokenmock}
				mocks.user.EXPECT().GetUserAuth(gomock.Any(), in).Return(nil, assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth"
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestSignIn(t *testing.T) {
	t.Parallel()

	auth := &uentity.UserAuth{
		UserID:       "user-id",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.SignInRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.SignInUserInput{
					Key:      "username",
					Password: "password",
				}
				mocks.user.EXPECT().SignInUser(gomock.Any(), in).Return(auth, nil)
			},
			req: &request.SignInRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthResponse{
					Auth: &response.Auth{
						UserID:       "user-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					},
				},
			},
		},
		{
			name: "failed to sign in user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.SignInUserInput{
					Key:      "username",
					Password: "password",
				}
				mocks.user.EXPECT().SignInUser(gomock.Any(), in).Return(nil, assert.AnError)
			},
			req: &request.SignInRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestSignOut(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.SignOutUserInput{AccessToken: tokenmock}
				mocks.user.EXPECT().SignOutUser(gomock.Any(), in).Return(nil)
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to sign out user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.SignOutUserInput{AccessToken: tokenmock}
				mocks.user.EXPECT().SignOutUser(gomock.Any(), in).Return(assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth"
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}

func TestRefreshAuthToken(t *testing.T) {
	t.Parallel()

	auth := &uentity.UserAuth{
		UserID:       "user-id",
		AccessToken:  "access-token",
		RefreshToken: "",
		ExpiresIn:    3600,
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.RefreshAuthTokenRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.RefreshUserTokenInput{
					RefreshToken: "refresh-token",
				}
				mocks.user.EXPECT().RefreshUserToken(gomock.Any(), in).Return(auth, nil)
			},
			req: &request.RefreshAuthTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthResponse{
					Auth: &response.Auth{
						UserID:       "user-id",
						AccessToken:  "access-token",
						RefreshToken: "",
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					},
				},
			},
		},
		{
			name: "failed to sign in user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.RefreshUserTokenInput{
					RefreshToken: "refresh-token",
				}
				mocks.user.EXPECT().RefreshUserToken(gomock.Any(), in).Return(nil, assert.AnError)
			},
			req: &request.RefreshAuthTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/refresh-token"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestInitializeAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.InitializeAuthRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.InitializeUserInput{
					UserID:    idmock,
					AccountID: "account-id",
					Username:  "username",
				}
				mocks.user.EXPECT().InitializeUser(gomock.Any(), in).Return(nil)
			},
			req: &request.InitializeAuthRequest{
				AccountID: "account-id",
				Username:  "username",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to initialize user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.InitializeUserInput{
					UserID:    idmock,
					AccountID: "account-id",
					Username:  "username",
				}
				mocks.user.EXPECT().InitializeUser(gomock.Any(), in).Return(assert.AnError)
			},
			req: &request.InitializeAuthRequest{
				AccountID: "account-id",
				Username:  "username",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/initialized"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateUserAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.UpdateAuthEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateUserEmailInput{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateUserEmail(gomock.Any(), in).Return(nil)
			},
			req: &request.UpdateAuthEmailRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update user email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateUserEmailInput{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateUserEmail(gomock.Any(), in).Return(assert.AnError)
			},
			req: &request.UpdateAuthEmailRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/email"
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestVerifyAuthEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.VerifyAuthEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserEmailInput{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyUserEmail(gomock.Any(), in).Return(nil)
			},
			req: &request.VerifyAuthEmailRequest{
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to veirify user email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserEmailInput{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyUserEmail(gomock.Any(), in).Return(assert.AnError)
			},
			req: &request.VerifyAuthEmailRequest{
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/email/verified"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateAuthPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.UpdateAuthPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateUserPasswordInput{
					AccessToken:          tokenmock,
					OldPassword:          "!Qaz2wsx",
					NewPassword:          "!Qaz3edc",
					PasswordConfirmation: "!Qaz3edc",
				}
				mocks.user.EXPECT().UpdateUserPassword(gomock.Any(), in).Return(nil)
			},
			req: &request.UpdateAuthPasswordRequest{
				OldPassword:          "!Qaz2wsx",
				NewPassword:          "!Qaz3edc",
				PasswordConfirmation: "!Qaz3edc",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update user password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateUserPasswordInput{
					AccessToken:          tokenmock,
					OldPassword:          "!Qaz2wsx",
					NewPassword:          "!Qaz3edc",
					PasswordConfirmation: "!Qaz3edc",
				}
				mocks.user.EXPECT().UpdateUserPassword(gomock.Any(), in).Return(assert.AnError)
			},
			req: &request.UpdateAuthPasswordRequest{
				OldPassword:          "!Qaz2wsx",
				NewPassword:          "!Qaz3edc",
				PasswordConfirmation: "!Qaz3edc",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/password"
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestForgotAuthPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.ForgotAuthPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.ForgotUserPasswordInput{Email: "test-user@and-period.jp"}
				mocks.user.EXPECT().ForgotUserPassword(gomock.Any(), in).Return(nil)
			},
			req: &request.ForgotAuthPasswordRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to forget user password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.ForgotUserPasswordInput{Email: "test-user@and-period.jp"}
				mocks.user.EXPECT().ForgotUserPassword(gomock.Any(), in).Return(assert.AnError)
			},
			req: &request.ForgotAuthPasswordRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/forgot-password"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestResetAuthPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.ResetAuthPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserPasswordInput{
					Email:                "test-user@and-period.jp",
					VerifyCode:           "123456",
					NewPassword:          "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				mocks.user.EXPECT().VerifyUserPassword(gomock.Any(), in).Return(nil)
			},
			req: &request.ResetAuthPasswordRequest{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				Password:             "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to verify user password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserPasswordInput{
					Email:                "test-user@and-period.jp",
					VerifyCode:           "123456",
					NewPassword:          "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				mocks.user.EXPECT().VerifyUserPassword(gomock.Any(), in).Return(assert.AnError)
			},
			req: &request.ResetAuthPasswordRequest{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				Password:             "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/forgot-password/verified"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestGetAuthUser(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	u := &uentity.User{
		ID:        "user-id",
		CreatedAt: now,
		UpdatedAt: now,
		Member: uentity.Member{
			UserID:       "user-id",
			ProviderType: uentity.ProviderTypeEmail,
			Email:        "test@and-period.jp",
			PhoneNumber:  "+819012345678",
			Username:     "username",
			ThumbnailURL: "https://and-period.jp/thumbnail.png",
			VerifiedAt:   now,
		},
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetUserInput{UserID: idmock}
				mocks.user.EXPECT().GetUser(gomock.Any(), in).Return(u, nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthUserResponse{
					AuthUser: &response.AuthUser{
						ID:           "user-id",
						Username:     "username",
						ThumbnailURL: "https://and-period.jp/thumbnail.png",
					},
				},
			},
		},
		{
			name: "failed to get user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetUserInput{UserID: idmock}
				mocks.user.EXPECT().GetUser(gomock.Any(), in).Return(nil, assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/user"
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateAuthRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.CreateUserInput{
					Email:                "test-user@and-period.jp",
					PhoneNumber:          "+819012345678",
					Password:             "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				mocks.user.EXPECT().CreateUser(gomock.Any(), in).Return("user-id", nil)
			},
			req: &request.CreateAuthRequest{
				Email:                "test-user@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CreateAuthResponse{
					ID: "user-id",
				},
			},
		},
		{
			name: "failed to create user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.CreateUserInput{
					Email:                "test-user@and-period.jp",
					PhoneNumber:          "+819012345678",
					Password:             "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				mocks.user.EXPECT().CreateUser(gomock.Any(), in).Return("", assert.AnError)
			},
			req: &request.CreateAuthRequest{
				Email:                "test-user@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/user"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestVerifyAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.VerifyAuthRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserInput{
					UserID:     "user-id",
					VerifyCode: "123456",
				}
				mocks.user.EXPECT().VerifyUser(gomock.Any(), in).Return(nil)
			},
			req: &request.VerifyAuthRequest{
				ID:         "user-id",
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to verify user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserInput{
					UserID:     "user-id",
					VerifyCode: "123456",
				}
				mocks.user.EXPECT().VerifyUser(gomock.Any(), in).Return(assert.AnError)
			},
			req: &request.VerifyAuthRequest{
				ID:         "user-id",
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/user/verified"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestCreateAuthWithOAuth(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	u := &uentity.User{
		ID:        "user-id",
		CreatedAt: now,
		UpdatedAt: now,
		Member: uentity.Member{
			UserID:       "user-id",
			ProviderType: uentity.ProviderTypeEmail,
			Email:        "test@and-period.jp",
			PhoneNumber:  "+819012345678",
			Username:     "username",
			ThumbnailURL: "https://and-period.jp/thumbnail.png",
			VerifiedAt:   now,
		},
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.CreateUserWithOAuthInput{AccessToken: tokenmock}
				mocks.user.EXPECT().CreateUserWithOAuth(gomock.Any(), in).Return(u, nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthUserResponse{
					AuthUser: &response.AuthUser{
						ID:           "user-id",
						Username:     "username",
						ThumbnailURL: "https://and-period.jp/thumbnail.png",
					},
				},
			},
		},
		{
			name: "failed to create user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.CreateUserWithOAuthInput{AccessToken: tokenmock}
				mocks.user.EXPECT().CreateUserWithOAuth(gomock.Any(), in).Return(nil, assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/user/oauth"
			testPost(t, tt.setup, tt.expect, path, nil)
		})
	}
}

func TestDeleteAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.DeleteUserInput{UserID: idmock}
				mocks.user.EXPECT().DeleteUser(gomock.Any(), in).Return(nil)
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.DeleteUserInput{UserID: idmock}
				mocks.user.EXPECT().DeleteUser(gomock.Any(), in).Return(assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/user"
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
