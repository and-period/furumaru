package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/golang/mock/gomock"
)

func TestGetAuth(t *testing.T) {
	t.Parallel()

	auth := &uentity.AdminAuth{
		AdminID:      "admin-id",
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
				in := &user.GetAdminAuthInput{AccessToken: tokenmock}
				mocks.user.EXPECT().GetAdminAuth(gomock.Any(), in).Return(auth, nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthResponse{
					Auth: &response.Auth{
						AdminID:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					},
				},
			},
		},
		{
			name: "failed to get admin auth",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetAdminAuthInput{AccessToken: tokenmock}
				mocks.user.EXPECT().GetAdminAuth(gomock.Any(), in).Return(nil, errmock)
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

	auth := &uentity.AdminAuth{
		AdminID:      "admin-id",
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
				in := &user.SignInAdminInput{
					Key:      "username",
					Password: "password",
				}
				mocks.user.EXPECT().SignInAdmin(gomock.Any(), in).Return(auth, nil)
			},
			req: &request.SignInRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthResponse{
					Auth: &response.Auth{
						AdminID:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					},
				},
			},
		},
		{
			name: "failed to sign in admin",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.SignInAdminInput{
					Key:      "username",
					Password: "password",
				}
				mocks.user.EXPECT().SignInAdmin(gomock.Any(), in).Return(nil, errmock)
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
				in := &user.SignOutAdminInput{AccessToken: tokenmock}
				mocks.user.EXPECT().SignOutAdmin(gomock.Any(), in).Return(nil)
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to sign out admin",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.SignOutAdminInput{AccessToken: tokenmock}
				mocks.user.EXPECT().SignOutAdmin(gomock.Any(), in).Return(errmock)
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

	auth := &uentity.AdminAuth{
		AdminID:      "admin-id",
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
				in := &user.RefreshAdminTokenInput{
					RefreshToken: "refresh-token",
				}
				mocks.user.EXPECT().RefreshAdminToken(gomock.Any(), in).Return(auth, nil)
			},
			req: &request.RefreshAuthTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthResponse{
					Auth: &response.Auth{
						AdminID:      "admin-id",
						AccessToken:  "access-token",
						RefreshToken: "",
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					},
				},
			},
		},
		{
			name: "failed to sign in admin",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.RefreshAdminTokenInput{
					RefreshToken: "refresh-token",
				}
				mocks.user.EXPECT().RefreshAdminToken(gomock.Any(), in).Return(nil, errmock)
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

func TestRegisterDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.RegisterAuthDeviceRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.RegisterAdminDeviceInput{
					AdminID: idmock,
					Device:  "device",
				}
				mocks.user.EXPECT().RegisterAdminDevice(gomock.Any(), in).Return(nil)
			},
			req: &request.RegisterAuthDeviceRequest{
				Device: "device",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update device",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.RegisterAdminDeviceInput{
					AdminID: idmock,
					Device:  "device",
				}
				mocks.user.EXPECT().RegisterAdminDevice(gomock.Any(), in).Return(errmock)
			},
			req: &request.RegisterAuthDeviceRequest{
				Device: "device",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/auth/device"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateAuthEmail(t *testing.T) {
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
				in := &user.UpdateAdminEmailInput{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateAdminEmail(gomock.Any(), in).Return(nil)
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
				in := &user.UpdateAdminEmailInput{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateAdminEmail(gomock.Any(), in).Return(errmock)
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
				in := &user.VerifyAdminEmailInput{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyAdminEmail(gomock.Any(), in).Return(nil)
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
				in := &user.VerifyAdminEmailInput{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyAdminEmail(gomock.Any(), in).Return(errmock)
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

func TestUpdateAdminPassword(t *testing.T) {
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
				in := &user.UpdateAdminPasswordInput{
					AccessToken:          tokenmock,
					OldPassword:          "!Qaz2wsx",
					NewPassword:          "!Qaz3edc",
					PasswordConfirmation: "!Qaz3edc",
				}
				mocks.user.EXPECT().UpdateAdminPassword(gomock.Any(), in).Return(nil)
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
				in := &user.UpdateAdminPasswordInput{
					AccessToken:          tokenmock,
					OldPassword:          "!Qaz2wsx",
					NewPassword:          "!Qaz3edc",
					PasswordConfirmation: "!Qaz3edc",
				}
				mocks.user.EXPECT().UpdateAdminPassword(gomock.Any(), in).Return(errmock)
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
