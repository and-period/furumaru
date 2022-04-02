package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/marche/api/internal/gateway/user/v1/request"
	"github.com/and-period/marche/api/internal/gateway/user/v1/response"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/and-period/marche/api/proto/user"
	"github.com/golang/mock/gomock"
)

func TestGetUserMe(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	u := &user.User{
		Id:           "user-id",
		ProviderType: user.ProviderType_PROVIDER_TYPE_EMAIL,
		Email:        "test@and-period.jp",
		PhoneNumber:  "+819012345678",
		CreatedAt:    now.Unix(),
		UpdatedAt:    now.Unix(),
		VerifiedAt:   now.Unix(),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetUserRequest{UserId: idmock}
				out := &user.GetUserResponse{User: u}
				mocks.user.EXPECT().GetUser(gomock.Any(), in).Return(out, nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.UserMeResponse{
					ID:          "user-id",
					Email:       "test@and-period.jp",
					PhoneNumber: "+819012345678",
				},
			},
		},
		{
			name: "failed to get user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetUserRequest{UserId: idmock}
				mocks.user.EXPECT().GetUser(gomock.Any(), in).Return(nil, errmock)
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
			const path = "/v1/users/me"
			req := newHTTPRequest(t, http.MethodGet, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateUserRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.CreateUserRequest{
					Email:                "test-user@and-period.jp",
					PhoneNumber:          "+819012345678",
					Password:             "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				out := &user.CreateUserResponse{UserId: "user-id"}
				mocks.user.EXPECT().CreateUser(gomock.Any(), in).Return(out, nil)
			},
			req: &request.CreateUserRequest{
				Email:                "test-user@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CreateUserResponse{
					ID: "user-id",
				},
			},
		},
		{
			name: "failed to create user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.CreateUserRequest{
					Email:                "test-user@and-period.jp",
					PhoneNumber:          "+819012345678",
					Password:             "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				mocks.user.EXPECT().CreateUser(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateUserRequest{
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
			t.Parallel()
			const path = "/v1/users"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestVerifyUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.VerifyUserRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserRequest{
					UserId:     "user-id",
					VerifyCode: "123456",
				}
				out := &user.VerifyUserResponse{}
				mocks.user.EXPECT().VerifyUser(gomock.Any(), in).Return(out, nil)
			},
			req: &request.VerifyUserRequest{
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
				in := &user.VerifyUserRequest{
					UserId:     "user-id",
					VerifyCode: "123456",
				}
				mocks.user.EXPECT().VerifyUser(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.VerifyUserRequest{
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
			t.Parallel()
			const path = "/v1/users/verified"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestCreateUserWithOAuth(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	u := &user.User{
		Id:           "user-id",
		ProviderType: user.ProviderType_PROVIDER_TYPE_EMAIL,
		Email:        "test@and-period.jp",
		PhoneNumber:  "+819012345678",
		CreatedAt:    now.Unix(),
		UpdatedAt:    now.Unix(),
		VerifiedAt:   now.Unix(),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.CreateUserWithOAuthRequest{AccessToken: tokenmock}
				out := &user.CreateUserWithOAuthResponse{User: u}
				mocks.user.EXPECT().CreateUserWithOAuth(gomock.Any(), in).Return(out, nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.UserMeResponse{
					ID:          "user-id",
					Email:       "test@and-period.jp",
					PhoneNumber: "+819012345678",
				},
			},
		},
		{
			name: "failed to create user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.CreateUserWithOAuthRequest{AccessToken: tokenmock}
				mocks.user.EXPECT().CreateUserWithOAuth(gomock.Any(), in).Return(nil, errmock)
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
			const path = "/v1/users/oauth"
			req := newHTTPRequest(t, http.MethodPost, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUpdateUserEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.UpdateUserEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateUserEmailRequest{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				out := &user.UpdateUserEmailResponse{}
				mocks.user.EXPECT().UpdateUserEmail(gomock.Any(), in).Return(out, nil)
			},
			req: &request.UpdateUserEmailRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update user email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateUserEmailRequest{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateUserEmail(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.UpdateUserEmailRequest{
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
			t.Parallel()
			const path = "/v1/users/me/email"
			req := newHTTPRequest(t, http.MethodPatch, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestVerifyUserEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.VerifyUserEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserEmailRequest{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				out := &user.VerifyUserEmailResponse{}
				mocks.user.EXPECT().VerifyUserEmail(gomock.Any(), in).Return(out, nil)
			},
			req: &request.VerifyUserEmailRequest{
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to veirify user email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserEmailRequest{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyUserEmail(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.VerifyUserEmailRequest{
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
			t.Parallel()
			const path = "/v1/users/me/email/verified"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUpdateUserPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.UpdateUserPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			req: &request.UpdateUserPasswordRequest{
				OldPassword:          "!Qaz2wsx",
				NewPassword:          "!Qaz3edc",
				PasswordConfirmation: "!Qaz3edc",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/users/me/password"
			req := newHTTPRequest(t, http.MethodPatch, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestForgotUserPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.ForgotUserPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			req: &request.ForgotUserPasswordRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/users/me/forgot-password"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestResetUserPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.ResetUserPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			req: &request.ResetUserPasswordRequest{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				Password:             "!Qaz2wsx",
				PasswordConfirmation: "!Qaz2wsx",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/users/me/forgot-password/verified"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/users/me"
			req := newHTTPRequest(t, http.MethodDelete, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}
