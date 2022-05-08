package service

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/stretchr/testify/assert"
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
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *SignInUserInput
		expect    *entity.UserAuth
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(u, nil)
			},
			input: &SignInUserInput{
				Key:      "username",
				Password: "password",
			},
			expect: &entity.UserAuth{
				UserID:       "user-id",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &SignInUserInput{},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to sign in",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().SignIn(ctx, "username", "password").Return(nil, errmock)
			},
			input: &SignInUserInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: ErrUnauthenticated,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("", errmock)
			},
			input: &SignInUserInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().SignIn(ctx, "username", "password").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(nil, errmock)
			},
			input: &SignInUserInput{
				Key:      "username",
				Password: "password",
			},
			expect:    nil,
			expectErr: ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.SignInUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestSignOutUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *SignOutUserInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().SignOut(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil)
			},
			input: &SignOutUserInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &SignOutUserInput{
				AccessToken: "",
			},
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to sign out",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().SignOut(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(errmock)
			},
			input: &SignOutUserInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expectErr: ErrUnauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			err := service.SignOutUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestGetUserAuth(t *testing.T) {
	t.Parallel()

	u := &entity.User{ID: "user-id"}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *GetUserAuthInput
		expect    *entity.UserAuth
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(u, nil)
			},
			input: &GetUserAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &entity.UserAuth{
				UserID:       "user-id",
				AccessToken:  "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				RefreshToken: "",
				ExpiresIn:    0,
			},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &GetUserAuthInput{
				AccessToken: "",
			},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", errmock)
			},
			input: &GetUserAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: ErrUnauthenticated,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(u, errmock)
			},
			input: &GetUserAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: ErrUnauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.GetUserAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
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
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *RefreshUserTokenInput
		expect    *entity.UserAuth
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(u, nil)
			},
			input: &RefreshUserTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect: &entity.UserAuth{
				UserID:       "user-id",
				AccessToken:  "access-token",
				RefreshToken: "",
				ExpiresIn:    3600,
			},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &RefreshUserTokenInput{
				RefreshToken: "",
			},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to sign in",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil, errmock)
			},
			input: &RefreshUserTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: ErrUnauthenticated,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("", errmock)
			},
			input: &RefreshUserTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().RefreshToken(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(result, nil)
				mocks.userAuth.EXPECT().GetUsername(ctx, "access-token").Return("username", nil)
				mocks.db.User.EXPECT().GetByCognitoID(ctx, "username", "id").Return(nil, errmock)
			},
			input: &RefreshUserTokenInput{
				RefreshToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expect:    nil,
			expectErr: ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *userService) {
			actual, err := service.RefreshUserToken(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
