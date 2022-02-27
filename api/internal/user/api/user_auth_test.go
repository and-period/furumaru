package api

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func TestSignInUser(t *testing.T) {
	t.Parallel()

	result := &cognito.AuthResult{
		IDToken:      "id-token",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}
	u := &entity.User{ID: "user-id"}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.SignInUserRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(u, nil)
			},
			req: &user.SignInUserRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.SignInUserResponse{
					Auth: &user.Auth{
						UserId:       "user-id",
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    3600,
					},
				},
			},
		},
		{
			name: "failed to sign in",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().SignIn(ctx, "username", "password").Return(nil, errmock)
			},
			req: &user.SignInUserRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: codes.Unauthenticated,
			},
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("", errmock)
			},
			req: &user.SignInUserRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(nil, errmock)
			},
			req: &user.SignInUserRequest{
				Username: "username",
				Password: "password",
			},
			expect: &testResponse{
				code: codes.Internal,
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
			name: "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().SignOut(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil)
			},
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
		{
			name: "failed to sign out",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().SignOut(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(errmock)
			},
			req: &user.SignOutUserRequest{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &testResponse{
				code: codes.Unauthenticated,
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

	result := &cognito.AuthResult{
		IDToken:      "id-token",
		AccessToken:  "access-token",
		RefreshToken: "",
		ExpiresIn:    3600,
	}
	u := &entity.User{ID: "user-id"}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, t *testing.T, mocks *mocks)
		req    *user.RefreshUserTokenRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(u, nil)
			},
			req: &user.RefreshUserTokenRequest{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &testResponse{
				code: codes.OK,
				body: &user.RefreshUserTokenResponse{
					Auth: &user.Auth{
						UserId:       "user-id",
						AccessToken:  "access-token",
						RefreshToken: "",
						ExpiresIn:    3600,
					},
				},
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
		{
			name: "failed to sign in",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil, errmock)
			},
			req: &user.RefreshUserTokenRequest{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &testResponse{
				code: codes.Unauthenticated,
			},
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("", errmock)
			},
			req: &user.RefreshUserTokenRequest{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &testResponse{
				code: codes.Internal,
			},
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.userAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(nil, errmock)
			},
			req: &user.RefreshUserTokenRequest{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &testResponse{
				code: codes.Internal,
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
