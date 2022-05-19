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
			t.Parallel()
			const path = "/v1/auth"
			req := newHTTPRequest(t, http.MethodGet, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
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
			t.Parallel()
			const path = "/v1/auth"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
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
			t.Parallel()
			const path = "/v1/auth"
			req := newHTTPRequest(t, http.MethodDelete, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
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
			t.Parallel()
			const path = "/v1/auth/refresh-token"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}
