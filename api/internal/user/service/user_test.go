package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	params := &database.ListUsersParams{
		Limit:          20,
		Offset:         0,
		OnlyRegistered: true,
		OnlyVerified:   true,
	}
	users := entity.Users{
		{
			ID:         "user-id",
			Registered: true,
			CreatedAt:  now,
			UpdatedAt:  now,
			Member: entity.Member{
				AccountID:    "",
				CognitoID:    "cognito-id",
				Username:     "",
				ProviderType: entity.UserAuthProviderTypeEmail,
				Email:        "test-user@and-period.jp",
				PhoneNumber:  "+810000000000",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
				CreatedAt:    now,
				UpdatedAt:    now,
				VerifiedAt:   now,
			},
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *user.ListUsersInput
		expect      entity.Users
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().List(gomock.Any(), params).Return(users, nil)
				mocks.db.User.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListUsersInput{
				Limit:          20,
				Offset:         0,
				OnlyRegistered: true,
				OnlyVerified:   true,
			},
			expect:      users,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &user.ListUsersInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list users",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.User.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListUsersInput{
				Limit:          20,
				Offset:         0,
				OnlyRegistered: true,
				OnlyVerified:   true,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count users",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().List(gomock.Any(), params).Return(users, nil)
				mocks.db.User.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &user.ListUsersInput{
				Limit:          20,
				Offset:         0,
				OnlyRegistered: true,
				OnlyVerified:   true,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, total, err := service.ListUsers(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.ElementsMatch(t, tt.expect, actual)
				assert.Equal(t, tt.expectTotal, total)
			}),
		)
	}
}

func TestMultiGetUsers(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	users := entity.Users{
		{
			ID:         "user-id",
			Registered: true,
			CreatedAt:  now,
			UpdatedAt:  now,
			Member: entity.Member{
				AccountID:    "",
				CognitoID:    "cognito-id",
				Username:     "",
				ProviderType: entity.UserAuthProviderTypeEmail,
				Email:        "test-user@and-period.jp",
				PhoneNumber:  "+810000000000",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
				CreatedAt:    now,
				UpdatedAt:    now,
				VerifiedAt:   now,
			},
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetUsersInput
		expect    entity.Users
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().MultiGet(ctx, []string{"user-id"}).Return(users, nil)
			},
			input: &user.MultiGetUsersInput{
				UserIDs: []string{"user-id"},
			},
			expect: entity.Users{
				{
					ID:         "user-id",
					Registered: true,
					CreatedAt:  now,
					UpdatedAt:  now,
					Member: entity.Member{
						AccountID:    "",
						CognitoID:    "cognito-id",
						Username:     "",
						ProviderType: entity.UserAuthProviderTypeEmail,
						Email:        "test-user@and-period.jp",
						PhoneNumber:  "+810000000000",
						ThumbnailURL: "https://and-period.jp/thumbnail.png",
						CreatedAt:    now,
						UpdatedAt:    now,
						VerifiedAt:   now,
					},
				},
			},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetUsersInput{
				UserIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().
					MultiGet(ctx, []string{"user-id"}).
					Return(nil, assert.AnError)
			},
			input: &user.MultiGetUsersInput{
				UserIDs: []string{"user-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.MultiGetUsers(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.ElementsMatch(t, tt.expect, actual)
			}),
		)
	}
}

func TestMultiGetUserDevices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetUserDevicesInput
		expect    []string
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetUserDevicesInput{
				UserIDs: []string{"user-id"},
			},
			expect:    []string{},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetUserDevicesInput{
				UserIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.MultiGetUserDevices(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.ElementsMatch(t, tt.expect, actual)
			}),
		)
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	u := &entity.User{
		ID:         "user-id",
		Registered: true,
		CreatedAt:  now,
		UpdatedAt:  now,
		Member: entity.Member{
			AccountID:    "",
			CognitoID:    "cognito-id",
			Username:     "",
			ProviderType: entity.UserAuthProviderTypeEmail,
			Email:        "test-user@and-period.jp",
			PhoneNumber:  "+810000000000",
			ThumbnailURL: "https://and-period.jp/thumbnail.png",
			CreatedAt:    now,
			UpdatedAt:    now,
			VerifiedAt:   now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetUserInput
		expect    *entity.User
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(u, nil)
			},
			input: &user.GetUserInput{
				UserID: "user-id",
			},
			expect: &entity.User{
				ID:         "user-id",
				Registered: true,
				CreatedAt:  now,
				UpdatedAt:  now,
				Member: entity.Member{
					AccountID:    "",
					CognitoID:    "cognito-id",
					Username:     "",
					ProviderType: entity.UserAuthProviderTypeEmail,
					Email:        "test-user@and-period.jp",
					PhoneNumber:  "+810000000000",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
					CreatedAt:    now,
					UpdatedAt:    now,
					VerifiedAt:   now,
				},
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetUserInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(nil, assert.AnError)
			},
			input: &user.GetUserInput{
				UserID: "user-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.GetUser(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	m := &entity.User{
		ID:         "user-id",
		Registered: true,
		Status:     entity.UserStatusVerified,
		Member: entity.Member{
			UserID:       "user-id",
			CognitoID:    "cognito-id",
			ProviderType: entity.UserAuthProviderTypeEmail,
			Email:        "test-user@and-period.jp",
			PhoneNumber:  "+810000000000",
			CreatedAt:    now,
			UpdatedAt:    now,
			VerifiedAt:   now,
		},
	}
	g := &entity.User{
		ID:         "user-id",
		Status:     entity.UserStatusGuest,
		Registered: false,
		Guest: entity.Guest{
			UserID:    "user-id",
			Email:     "test-user@and-period.jp",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.DeleteUserInput
		expectErr error
	}{
		{
			name: "success member",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(m, nil)
				mocks.userAuth.EXPECT().DeleteUser(ctx, "cognito-id").Return(nil)
				mocks.db.Member.EXPECT().
					Delete(ctx, "user-id", gomock.Any()).
					DoAndReturn(func(ctx context.Context, userID string, auth func(ctx context.Context) error) error {
						return auth(ctx)
					})
			},
			input: &user.DeleteUserInput{
				UserID: "user-id",
			},
			expectErr: nil,
		},
		{
			name: "success guest",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(g, nil)
				mocks.db.Guest.EXPECT().Delete(ctx, "user-id").Return(nil)
			},
			input: &user.DeleteUserInput{
				UserID: "user-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.DeleteUserInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(nil, assert.AnError)
			},
			input: &user.DeleteUserInput{
				UserID: "user-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to delete member",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(m, nil)
				mocks.db.Member.EXPECT().Delete(ctx, "user-id", gomock.Any()).Return(assert.AnError)
			},
			input: &user.DeleteUserInput{
				UserID: "user-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to delete guest",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(g, nil)
				mocks.db.Guest.EXPECT().Delete(ctx, "user-id").Return(assert.AnError)
			},
			input: &user.DeleteUserInput{
				UserID: "user-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.DeleteUser(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}
