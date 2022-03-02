package api

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/and-period/marche/api/proto/user"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	u := &entity.User{
		ID:           "user-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeEmail,
		Email:        "test-user@and-period.jp",
		PhoneNumber:  "+810000000000",
		CreatedAt:    now,
		UpdatedAt:    now,
		VerifiedAt:   now,
	}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.GetUserRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(u, nil)
			},
			req: &user.GetUserRequest{
				UserId: "user-id",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.GetUserResponse{
					User: &user.User{
						Id:           "user-id",
						ProviderType: user.ProviderType_PROVIDER_TYPE_EMAIL,
						Email:        "test-user@and-period.jp",
						PhoneNumber:  "+810000000000",
						CreatedAt:    now.Unix(),
						UpdatedAt:    now.Unix(),
						VerifiedAt:   now.Unix(),
					},
				},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.GetUserRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(nil, errmock)
			},
			req: &user.GetUserRequest{
				UserId: "user-id",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.GetUser(ctx, tt.req)
		}))
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.CreateUserRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.db.User.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.userAuth.EXPECT().SignUp(ctx, gomock.Any()).Return(nil)
			},
			req: &user.CreateUserRequest{
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
			expect: &testResponse{
				code: codes.OK,
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.CreateUserRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name:  "failed to unmatch password",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.CreateUserRequest{
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "12345678",
				PasswordConfirmation: "11111111",
			},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.db.User.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			req: &user.CreateUserRequest{
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.db.User.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.userAuth.EXPECT().SignUp(ctx, gomock.Any()).Return(errmock)
			},
			req: &user.CreateUserRequest{
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.CreateUser(ctx, tt.req)
		}))
	}
}

func TestVerifyUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.VerifyUserRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "user-id", "123456").Return(nil)
				mocks.db.User.EXPECT().UpdateVerified(ctx, "user-id").Return(nil)
			},
			req: &user.VerifyUserRequest{
				UserId:     "user-id",
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.VerifyUserResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.VerifyUserRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "failed to confirm sign up",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "user-id", "123456").Return(errmock)
			},
			req: &user.VerifyUserRequest{
				UserId:     "user-id",
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to update verified",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "user-id", "123456").Return(nil)
				mocks.db.User.EXPECT().UpdateVerified(ctx, "user-id").Return(errmock)
			},
			req: &user.VerifyUserRequest{
				UserId:     "user-id",
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.VerifyUser(ctx, tt.req)
		}))
	}
}

func TestCreateUserWithOAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.CreateUserWithOAuthRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.CreateUserWithOAuthRequest{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.CreateUserWithOAuthResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.CreateUserWithOAuthRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.CreateUserWithOAuth(ctx, tt.req)
		}))
	}
}

func TestUpdateUserEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.UpdateUserEmailRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.UpdateUserEmailRequest{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test@and-period.jp",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.UpdateUserEmailResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.UpdateUserEmailRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.UpdateUserEmail(ctx, tt.req)
		}))
	}
}

func TestVerifyUserEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.VerifyUserEmailRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.VerifyUserEmailRequest{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.VerifyUserEmailResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.VerifyUserEmailRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.VerifyUserEmail(ctx, tt.req)
		}))
	}
}

func TestUpdateUserPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.UpdateUserPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.UpdateUserPasswordRequest{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.UpdateUserPasswordResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.UpdateUserPasswordRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.UpdateUserPassword(ctx, tt.req)
		}))
	}
}

func TestForgotUserPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.ForgotUserPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.ForgotUserPasswordRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.ForgotUserPasswordResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.ForgotUserPasswordRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.ForgotUserPassword(ctx, tt.req)
		}))
	}
}

func TestVerifyUserPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.VerifyUserPasswordRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.VerifyUserPasswordRequest{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.VerifyUserPasswordResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.VerifyUserPasswordRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.VerifyUserPassword(ctx, tt.req)
		}))
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.DeleteUserRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.DeleteUserRequest{
				UserId: "user-id",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.DeleteuserResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req:   &user.DeleteUserRequest{},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.DeleteUser(ctx, tt.req)
		}))
	}
}
