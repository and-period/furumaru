package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/marche/api/internal/gateway/user/v1/entity"
	"github.com/and-period/marche/api/internal/gateway/user/v1/request"
	"github.com/and-period/marche/api/internal/gateway/user/v1/response"
	"github.com/and-period/marche/api/proto/user"
	"github.com/golang/mock/gomock"
)

func TestSignIn(t *testing.T) {
	t.Parallel()

	auth := &user.Auth{
		UserId:       "user-id",
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
				in := &user.SignInUserRequest{
					Username: "username",
					Password: "password",
				}
				out := &user.SignInUserResponse{Auth: auth}
				mocks.user.EXPECT().SignInUser(gomock.Any(), in).Return(out, nil)
			},
			req: &request.SignInRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthResponse{
					Auth: &entity.Auth{
						UserID:       "user-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
					},
				},
			},
		},
		{
			name: "failed to sign in user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.SignInUserRequest{
					Username: "username",
					Password: "password",
				}
				mocks.user.EXPECT().SignInUser(gomock.Any(), in).Return(nil, errmock)
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

func TestRefreshAuthToken(t *testing.T) {
	t.Parallel()

	auth := &user.Auth{
		UserId:       "user-id",
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
				in := &user.RefreshUserTokenRequest{
					RefreshToken: "refresh-token",
				}
				out := &user.RefreshUserTokenResponse{Auth: auth}
				mocks.user.EXPECT().RefreshUserToken(gomock.Any(), in).Return(out, nil)
			},
			req: &request.RefreshAuthTokenRequest{
				RefreshToken: "refresh-token",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AuthResponse{
					Auth: &entity.Auth{
						UserID:       "user-id",
						AccessToken:  "access-token",
						RefreshToken: "",
						ExpiresIn:    3600,
					},
				},
			},
		},
		{
			name: "failed to sign in user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.RefreshUserTokenRequest{
					RefreshToken: "refresh-token",
				}
				mocks.user.EXPECT().RefreshUserToken(gomock.Any(), in).Return(nil, errmock)
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
