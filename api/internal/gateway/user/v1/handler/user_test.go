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
)

func TestGetUserMe(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	u := &uentity.User{
		ID:           "user-id",
		ProviderType: uentity.ProviderTypeEmail,
		Email:        "test@and-period.jp",
		PhoneNumber:  "+819012345678",
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
		CreatedAt:    now,
		UpdatedAt:    now,
		VerifiedAt:   now,
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
				body: &response.UserMeResponse{
					ID:           "user-id",
					Email:        "test@and-period.jp",
					PhoneNumber:  "+819012345678",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
				},
			},
		},
		{
			name: "failed to get user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetUserInput{UserID: idmock}
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
				in := &user.CreateUserInput{
					Email:                "test-user@and-period.jp",
					PhoneNumber:          "+819012345678",
					Password:             "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				mocks.user.EXPECT().CreateUser(gomock.Any(), in).Return("user-id", nil)
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
				in := &user.CreateUserInput{
					Email:                "test-user@and-period.jp",
					PhoneNumber:          "+819012345678",
					Password:             "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				mocks.user.EXPECT().CreateUser(gomock.Any(), in).Return("", errmock)
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
				in := &user.VerifyUserInput{
					UserID:     "user-id",
					VerifyCode: "123456",
				}
				mocks.user.EXPECT().VerifyUser(gomock.Any(), in).Return(nil)
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
				in := &user.VerifyUserInput{
					UserID:     "user-id",
					VerifyCode: "123456",
				}
				mocks.user.EXPECT().VerifyUser(gomock.Any(), in).Return(errmock)
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
	u := &uentity.User{
		ID:           "user-id",
		ProviderType: uentity.ProviderTypeEmail,
		Email:        "test@and-period.jp",
		PhoneNumber:  "+819012345678",
		CreatedAt:    now,
		UpdatedAt:    now,
		VerifiedAt:   now,
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
				in := &user.CreateUserWithOAuthInput{AccessToken: tokenmock}
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
				in := &user.UpdateUserEmailInput{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateUserEmail(gomock.Any(), in).Return(nil)
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
				in := &user.UpdateUserEmailInput{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateUserEmail(gomock.Any(), in).Return(errmock)
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
				in := &user.VerifyUserEmailInput{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyUserEmail(gomock.Any(), in).Return(nil)
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
				in := &user.VerifyUserEmailInput{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyUserEmail(gomock.Any(), in).Return(errmock)
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
			req: &request.UpdateUserPasswordRequest{
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
				mocks.user.EXPECT().UpdateUserPassword(gomock.Any(), in).Return(errmock)
			},
			req: &request.UpdateUserPasswordRequest{
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
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.ForgotUserPasswordInput{Email: "test-user@and-period.jp"}
				mocks.user.EXPECT().ForgotUserPassword(gomock.Any(), in).Return(nil)
			},
			req: &request.ForgotUserPasswordRequest{
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
				mocks.user.EXPECT().ForgotUserPassword(gomock.Any(), in).Return(errmock)
			},
			req: &request.ForgotUserPasswordRequest{
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
		{
			name: "failed to verify user password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyUserPasswordInput{
					Email:                "test-user@and-period.jp",
					VerifyCode:           "123456",
					NewPassword:          "!Qaz2wsx",
					PasswordConfirmation: "!Qaz2wsx",
				}
				mocks.user.EXPECT().VerifyUserPassword(gomock.Any(), in).Return(errmock)
			},
			req: &request.ResetUserPasswordRequest{
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
				mocks.user.EXPECT().DeleteUser(gomock.Any(), in).Return(errmock)
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
			req := newHTTPRequest(t, http.MethodDelete, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}
