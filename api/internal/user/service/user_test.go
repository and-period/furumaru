package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	params := &database.ListUsersParams{
		Limit:  20,
		Offset: 0,
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
				ProviderType: entity.ProviderTypeEmail,
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
				Limit:  20,
				Offset: 0,
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
				mocks.db.User.EXPECT().List(gomock.Any(), params).Return(nil, errmock)
				mocks.db.User.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListUsersInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
		{
			name: "failed to count users",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().List(gomock.Any(), params).Return(users, nil)
				mocks.db.User.EXPECT().Count(gomock.Any(), params).Return(int64(0), errmock)
			},
			input: &user.ListUsersInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListUsers(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
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
				ProviderType: entity.ProviderTypeEmail,
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
						ProviderType: entity.ProviderTypeEmail,
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
				mocks.db.User.EXPECT().MultiGet(ctx, []string{"user-id"}).Return(nil, errmock)
			},
			input: &user.MultiGetUsersInput{
				UserIDs: []string{"user-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetUsers(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
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
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetUserDevices(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
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
			ProviderType: entity.ProviderTypeEmail,
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
					ProviderType: entity.ProviderTypeEmail,
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
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(nil, errmock)
			},
			input: &user.GetUserInput{
				UserID: "user-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateUserInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectUser := &entity.User{
					Registered: true,
					Member: entity.Member{
						ProviderType: entity.ProviderTypeEmail,
						Email:        "test@and-period.jp",
						PhoneNumber:  "+819012345678",
					},
				}
				mocks.db.Member.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *entity.User, m *entity.Member) error {
						expectUser.ID = u.ID
						expectUser.Member.UserID, expectUser.Member.CognitoID = m.UserID, m.CognitoID
						assert.Equal(t, expectUser, u)
						assert.Equal(t, expectUser.Member, *m)
						return nil
					})
				mocks.userAuth.EXPECT().SignUp(ctx, gomock.Any()).Return(nil)
			},
			input: &user.CreateUserInput{
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateUserInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "failed to unmatch password",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.CreateUserInput{
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "12345678",
				PasswordConfirmation: "11111111",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateUserInput{
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(nil)
				mocks.userAuth.EXPECT().SignUp(ctx, gomock.Any()).Return(errmock)
			},
			input: &user.CreateUserInput{
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.VerifyUserInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "user-id", "123456").Return(nil)
				mocks.db.Member.EXPECT().UpdateVerified(ctx, "user-id").Return(nil)
			},
			input: &user.VerifyUserInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.VerifyUserInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to confirm sign up",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "user-id", "123456").Return(errmock)
			},
			input: &user.VerifyUserInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to update verified",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "user-id", "123456").Return(nil)
				mocks.db.Member.EXPECT().UpdateVerified(ctx, "user-id").Return(errmock)
			},
			input: &user.VerifyUserInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestCreateUserWithOAuth(t *testing.T) {
	t.Parallel()

	auth := &cognito.AuthUser{
		Username:    "cognito-id",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+810000000000",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateUserWithOAuthInput
		expect    *entity.User
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectUser := &entity.User{
					Registered: true,
					Member: entity.Member{
						ProviderType: entity.ProviderTypeOAuth,
						Email:        "test-user@and-period.jp",
						PhoneNumber:  "+810000000000",
					},
				}
				mocks.userAuth.EXPECT().GetUser(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(auth, nil)
				mocks.db.Member.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *entity.User, m *entity.Member) error {
						expectUser.ID = u.ID
						expectUser.Member.UserID, expectUser.CognitoID = m.UserID, m.CognitoID
						assert.Equal(t, expectUser, u)
						assert.Equal(t, expectUser.Member, *m)
						return nil
					})
			},
			input: &user.CreateUserWithOAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateUserWithOAuthInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUser(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil, errmock)
			},
			input: &user.CreateUserWithOAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to create user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUser(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(auth, nil)
				mocks.db.Member.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(errmock)
			},
			input: &user.CreateUserWithOAuthInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateUserWithOAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestInitializeUser(t *testing.T) {
	t.Parallel()
	m := &entity.Member{AccountID: ""}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.InitializeUserInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Get(ctx, "user-id", "account_id").Return(m, nil)
				mocks.db.Member.EXPECT().UpdateAccount(ctx, "user-id", "account-id", "username").Return(nil)
			},
			input: &user.InitializeUserInput{
				UserID:    "user-id",
				AccountID: "account-id",
				Username:  "username",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.InitializeUserInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Get(ctx, "user-id", "account_id").Return(nil, errmock)
			},
			input: &user.InitializeUserInput{
				UserID:    "user-id",
				AccountID: "account-id",
				Username:  "username",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				m := &entity.Member{AccountID: "account-id"}
				mocks.db.Member.EXPECT().Get(ctx, "user-id", "account_id").Return(m, nil)
			},
			input: &user.InitializeUserInput{
				UserID:    "user-id",
				AccountID: "account-id",
				Username:  "username",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to initilaze user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Get(ctx, "user-id", "account_id").Return(m, nil)
				mocks.db.Member.EXPECT().UpdateAccount(ctx, "user-id", "account-id", "username").Return(errmock)
			},
			input: &user.InitializeUserInput{
				UserID:    "user-id",
				AccountID: "account-id",
				Username:  "username",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.InitializeUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateUserEmail(t *testing.T) {
	t.Parallel()

	m := &entity.Member{
		ProviderType: entity.ProviderTypeEmail,
		Email:        "test-user@and-period.jp",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateUserEmailInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					OldEmail:    "test-user@and-period.jp",
					NewEmail:    "test-other@and-period.jp",
				}
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "provider_type", "email").Return(m, nil)
				mocks.userAuth.EXPECT().ChangeEmail(ctx, params).Return(nil)
			},
			input: &user.UpdateUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateUserEmailInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", errmock)
			},
			input: &user.UpdateUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "provider_type", "email").Return(nil, errmock)
			},
			input: &user.UpdateUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to unmatch provider type",
			setup: func(ctx context.Context, mocks *mocks) {
				m := &entity.Member{
					ProviderType: entity.ProviderTypeOAuth,
					Email:        "test-user@and-period.jp",
				}
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "provider_type", "email").Return(m, nil)
			},
			input: &user.UpdateUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to change email",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					OldEmail:    "test-user@and-period.jp",
					NewEmail:    "test-other@and-period.jp",
				}
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "provider_type", "email").Return(m, nil)
				mocks.userAuth.EXPECT().ChangeEmail(ctx, params).Return(errmock)
			},
			input: &user.UpdateUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateUserEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyUserEmail(t *testing.T) {
	t.Parallel()

	m := &entity.Member{
		UserID: "user-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.VerifyUserEmailInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					VerifyCode:  "123456",
				}
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "user_id").Return(m, nil)
				mocks.userAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test-user@and-period.jp", nil)
				mocks.db.Member.EXPECT().UpdateEmail(ctx, "user-id", "test-user@and-period.jp").Return(nil)
			},
			input: &user.VerifyUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.VerifyUserEmailInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", errmock)
			},
			input: &user.VerifyUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "user_id").Return(nil, errmock)
			},
			input: &user.VerifyUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to confirm change email",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					VerifyCode:  "123456",
				}
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "user_id").Return(m, nil)
				mocks.userAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("", errmock)
			},
			input: &user.VerifyUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to update email",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmChangeEmailParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					Username:    "cognito-id",
					VerifyCode:  "123456",
				}
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "user_id").Return(m, nil)
				mocks.userAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("test-user@and-period.jp", nil)
				mocks.db.Member.EXPECT().UpdateEmail(ctx, "user-id", "test-user@and-period.jp").Return(errmock)
			},
			input: &user.VerifyUserEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyUserEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateUserPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateUserPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangePasswordParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					OldPassword: "12345678",
					NewPassword: "12345678",
				}
				mocks.userAuth.EXPECT().ChangePassword(ctx, params).Return(nil)
			},
			input: &user.UpdateUserPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateUserPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid argument for password unmatch",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.UpdateUserPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "123456789",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to change password",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangePasswordParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					OldPassword: "12345678",
					NewPassword: "12345678",
				}
				mocks.userAuth.EXPECT().ChangePassword(ctx, params).Return(errmock)
			},
			input: &user.UpdateUserPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateUserPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestForgotUserPassword(t *testing.T) {
	t.Parallel()

	m := &entity.Member{CognitoID: "cognito-id"}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.ForgotUserPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(m, nil)
				mocks.userAuth.EXPECT().ForgotPassword(ctx, "cognito-id").Return(nil)
			},
			input: &user.ForgotUserPasswordInput{
				Email: "test-user@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ForgotUserPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(nil, errmock)
			},
			input: &user.ForgotUserPasswordInput{
				Email: "test-user@and-period.jp",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to forget password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(m, nil)
				mocks.userAuth.EXPECT().ForgotPassword(ctx, "cognito-id").Return(errmock)
			},
			input: &user.ForgotUserPasswordInput{
				Email: "test-user@and-period.jp",
			},
			expectErr: exception.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ForgotUserPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyUserPassword(t *testing.T) {
	t.Parallel()

	m := &entity.Member{CognitoID: "cognito-id"}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.VerifyUserPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmForgotPasswordParams{
					Username:    "cognito-id",
					VerifyCode:  "123456",
					NewPassword: "12345678",
				}
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(m, nil)
				mocks.userAuth.EXPECT().ConfirmForgotPassword(ctx, params).Return(nil)
			},
			input: &user.VerifyUserPasswordInput{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.VerifyUserPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.VerifyUserPasswordInput{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "12345678",
				PasswordConfirmation: "123456789",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(nil, errmock)
			},
			input: &user.VerifyUserPasswordInput{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to confirm forgot password",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmForgotPasswordParams{
					Username:    "cognito-id",
					VerifyCode:  "123456",
					NewPassword: "12345678",
				}
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(m, nil)
				mocks.userAuth.EXPECT().ConfirmForgotPassword(ctx, params).Return(errmock)
			},
			input: &user.VerifyUserPasswordInput{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "12345678",
				PasswordConfirmation: "12345678",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyUserPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	m := &entity.Member{
		UserID:       "user-id",
		CognitoID:    "cognito-id",
		ProviderType: entity.ProviderTypeEmail,
		Email:        "test-user@and-period.jp",
		PhoneNumber:  "+810000000000",
		CreatedAt:    now,
		UpdatedAt:    now,
		VerifiedAt:   now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.DeleteUserInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Get(ctx, "user-id").Return(m, nil)
				mocks.userAuth.EXPECT().DeleteUser(ctx, "cognito-id").Return(nil)
				mocks.db.Member.EXPECT().Delete(ctx, "user-id").Return(nil)
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
			name: "failed to delete cognito user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Get(ctx, "user-id").Return(m, errmock)
			},
			input: &user.DeleteUserInput{
				UserID: "user-id",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to delete cognito user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Get(ctx, "user-id").Return(m, nil)
				mocks.userAuth.EXPECT().DeleteUser(ctx, "cognito-id").Return(errmock)
			},
			input: &user.DeleteUserInput{
				UserID: "user-id",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to delete user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().Get(ctx, "user-id").Return(m, nil)
				mocks.userAuth.EXPECT().DeleteUser(ctx, "cognito-id").Return(nil)
				mocks.db.Member.EXPECT().Delete(ctx, "user-id").Return(errmock)
			},
			input: &user.DeleteUserInput{
				UserID: "user-id",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteUser(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
