package api

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func TestSignInUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.SignInUserRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.SignInUserRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.SignInUserResponse{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.SignInUser(ctx, tt.req)
		}))
	}
}

func TestSignOutUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.SignOutUserRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.SignOutUserRequest{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.SignOutUserResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.SignOutUserRequest{
				AccessToken: "",
			},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.SignOutUser(ctx, tt.req)
		}))
	}
}

func TestRefreshUserToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.RefreshUserTokenRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.RefreshUserTokenRequest{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.RefreshUserTokenResponse{},
			},
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {},
			req: &user.RefreshUserTokenRequest{
				RefreshToken: "",
			},
			expect: &testResponse{
				code: codes.InvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testGRPC(tt.setup, tt.expect, func(ctx context.Context, service *userService) (proto.Message, error) {
			return service.RefreshUserToken(ctx, tt.req)
		}))
	}
}
