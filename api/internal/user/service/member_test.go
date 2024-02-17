package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateMember(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateMemberInput
		expectErr error
	}{
		{
			name: "success create user",
			setup: func(ctx context.Context, mocks *mocks) {
				expectUser := &entity.User{
					Registered: true,
					Member: entity.Member{
						Username:      "username",
						AccountID:     "account-id",
						Lastname:      "&.",
						Firstname:     "利用者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "りようしゃ",
						ProviderType:  entity.ProviderTypeEmail,
						Email:         "test@and-period.jp",
						PhoneNumber:   "+819012345678",
					},
				}
				expectSignUp := &cognito.SignUpParams{
					Email:       "test@and-period.jp",
					PhoneNumber: "+819012345678",
					Password:    "Passw0rd",
				}
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test@and-period.jp").Return(nil, database.ErrNotFound)
				mocks.userAuth.EXPECT().
					SignUp(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *cognito.SignUpParams) error {
						expectSignUp.Username = params.Username
						assert.Equal(t, expectSignUp, params)
						return nil
					})
				mocks.db.Member.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *entity.User, auth func(ctx context.Context) error) error {
						expectUser.ID = u.ID
						expectUser.Member.UserID, expectUser.Member.CognitoID = u.ID, u.CognitoID
						assert.Equal(t, expectUser, u)
						return auth(ctx)
					})
			},
			input: &user.CreateMemberInput{
				Username:             "username",
				AccountID:            "account-id",
				Lastname:             "&.",
				Firstname:            "利用者",
				LastnameKana:         "あんどどっと",
				FirstnameKana:        "りようしゃ",
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: nil,
		},
		{
			name: "success resend confirmation code",
			setup: func(ctx context.Context, mocks *mocks) {
				member := &entity.Member{
					UserID:        "user-id",
					CognitoID:     "cognito-id",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					ProviderType:  entity.ProviderTypeEmail,
					Email:         "test@and-period.jp",
					PhoneNumber:   "+819012345678",
					VerifiedAt:    time.Time{},
				}
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test@and-period.jp").Return(member, nil)
				mocks.userAuth.EXPECT().ResendSignUpCode(ctx, "cognito-id").Return(nil)
			},
			input: &user.CreateMemberInput{
				Username:             "username",
				AccountID:            "account-id",
				Lastname:             "&.",
				Firstname:            "利用者",
				LastnameKana:         "あんどどっと",
				FirstnameKana:        "りようしゃ",
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateMemberInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "failed to unmatch password",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.CreateMemberInput{
				Username:             "username",
				AccountID:            "account-id",
				Lastname:             "&.",
				Firstname:            "利用者",
				LastnameKana:         "あんどどっと",
				FirstnameKana:        "りようしゃ",
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "Passw0rd",
				PasswordConfirmation: "11111111",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get member",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test@and-period.jp").Return(nil, assert.AnError)
			},
			input: &user.CreateMemberInput{
				Username:             "username",
				AccountID:            "account-id",
				Lastname:             "&.",
				Firstname:            "利用者",
				LastnameKana:         "あんどどっと",
				FirstnameKana:        "りようしゃ",
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to resend confirmation code",
			setup: func(ctx context.Context, mocks *mocks) {
				member := &entity.Member{
					UserID:        "user-id",
					CognitoID:     "cognito-id",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					ProviderType:  entity.ProviderTypeEmail,
					Email:         "test@and-period.jp",
					PhoneNumber:   "+819012345678",
					VerifiedAt:    time.Time{},
				}
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test@and-period.jp").Return(member, nil)
				mocks.userAuth.EXPECT().ResendSignUpCode(ctx, "cognito-id").Return(assert.AnError)
			},
			input: &user.CreateMemberInput{
				Username:             "username",
				AccountID:            "account-id",
				Lastname:             "&.",
				Firstname:            "利用者",
				LastnameKana:         "あんどどっと",
				FirstnameKana:        "りようしゃ",
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test@and-period.jp").Return(nil, database.ErrNotFound)
				mocks.db.Member.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateMemberInput{
				Username:             "username",
				AccountID:            "account-id",
				Lastname:             "&.",
				Firstname:            "利用者",
				LastnameKana:         "あんどどっと",
				FirstnameKana:        "りようしゃ",
				Email:                "test@and-period.jp",
				PhoneNumber:          "+819012345678",
				Password:             "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateMember(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyMember(t *testing.T) {
	t.Parallel()

	u := &entity.User{
		ID: "user-id",
		Member: entity.Member{
			UserID:    "user-id",
			CognitoID: "cognito-id",
		},
		Registered: true,
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.VerifyMemberInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(u, nil)
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "cognito-id", "123456").Return(nil)
				mocks.db.Member.EXPECT().UpdateVerified(ctx, "user-id").Return(nil)
			},
			input: &user.VerifyMemberInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: nil,
		},
		{
			name: "success resend signup code",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(u, nil)
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "cognito-id", "123456").Return(cognito.ErrCodeExpired)
				mocks.userAuth.EXPECT().ResendSignUpCode(ctx, "cognito-id").Return(nil)
			},
			input: &user.VerifyMemberInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.VerifyMemberInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(nil, assert.AnError)
			},
			input: &user.VerifyMemberInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to confirm sign up",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(u, nil)
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "cognito-id", "123456").Return(assert.AnError)
			},
			input: &user.VerifyMemberInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to resend signup code",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(u, nil)
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "cognito-id", "123456").Return(cognito.ErrCodeExpired)
				mocks.userAuth.EXPECT().ResendSignUpCode(ctx, "cognito-id").Return(assert.AnError)
			},
			input: &user.VerifyMemberInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to update verified",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.User.EXPECT().Get(ctx, "user-id").Return(u, nil)
				mocks.userAuth.EXPECT().ConfirmSignUp(ctx, "cognito-id", "123456").Return(nil)
				mocks.db.Member.EXPECT().UpdateVerified(ctx, "user-id").Return(assert.AnError)
			},
			input: &user.VerifyMemberInput{
				UserID:     "user-id",
				VerifyCode: "123456",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyMember(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestCreateMemberWithOAuth(t *testing.T) {
	t.Parallel()

	auth := &cognito.AuthUser{
		Username:    "cognito-id",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+810000000000",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateMemberWithOAuthInput
		expect    *entity.User
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectUser := &entity.User{
					Registered: true,
					Member: entity.Member{
						Username:      "username",
						AccountID:     "account-id",
						Lastname:      "&.",
						Firstname:     "利用者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "りようしゃ",
						ProviderType:  entity.ProviderTypeOAuth,
						Email:         "test-user@and-period.jp",
						PhoneNumber:   "+810000000000",
					},
				}
				mocks.userAuth.EXPECT().GetUser(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(auth, nil)
				mocks.db.Member.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *entity.User, auth func(ctx context.Context) error) error {
						expectUser.ID = u.ID
						expectUser.Member.UserID, expectUser.Member.CognitoID = u.ID, u.CognitoID
						assert.Equal(t, expectUser, u)
						return auth(ctx)
					})
			},
			input: &user.CreateMemberWithOAuthInput{
				AccessToken:   "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Username:      "username",
				AccountID:     "account-id",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				PhoneNumber:   "+810000000000",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateMemberWithOAuthInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUser(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(nil, assert.AnError)
			},
			input: &user.CreateMemberWithOAuthInput{
				AccessToken:   "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Username:      "username",
				AccountID:     "account-id",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				PhoneNumber:   "+810000000000",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUser(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return(auth, nil)
				mocks.db.Member.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateMemberWithOAuthInput{
				AccessToken:   "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Username:      "username",
				AccountID:     "account-id",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				PhoneNumber:   "+810000000000",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateMemberWithOAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateMemberEmail(t *testing.T) {
	t.Parallel()

	m := &entity.Member{
		ProviderType: entity.ProviderTypeEmail,
		Email:        "test-user@and-period.jp",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateMemberEmailInput
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
			input: &user.UpdateMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateMemberEmailInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", assert.AnError)
			},
			input: &user.UpdateMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "provider_type", "email").Return(nil, assert.AnError)
			},
			input: &user.UpdateMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrInternal,
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
			input: &user.UpdateMemberEmailInput{
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
				mocks.userAuth.EXPECT().ChangeEmail(ctx, params).Return(assert.AnError)
			},
			input: &user.UpdateMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				Email:       "test-other@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateMemberEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyMemberEmail(t *testing.T) {
	t.Parallel()

	m := &entity.Member{
		UserID: "user-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.VerifyMemberEmailInput
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
			input: &user.VerifyMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.VerifyMemberEmailInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("", assert.AnError)
			},
			input: &user.VerifyMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get by cognito id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.userAuth.EXPECT().GetUsername(ctx, "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ").Return("cognito-id", nil)
				mocks.db.Member.EXPECT().GetByCognitoID(ctx, "cognito-id", "user_id").Return(nil, assert.AnError)
			},
			input: &user.VerifyMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrInternal,
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
				mocks.userAuth.EXPECT().ConfirmChangeEmail(ctx, params).Return("", assert.AnError)
			},
			input: &user.VerifyMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrInternal,
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
				mocks.db.Member.EXPECT().UpdateEmail(ctx, "user-id", "test-user@and-period.jp").Return(assert.AnError)
			},
			input: &user.VerifyMemberEmailInput{
				AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				VerifyCode:  "123456",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyMemberEmail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateMemberPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateMemberPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ChangePasswordParams{
					AccessToken: "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
					OldPassword: "12345678",
					NewPassword: "Passw0rd",
				}
				mocks.userAuth.EXPECT().ChangePassword(ctx, params).Return(nil)
			},
			input: &user.UpdateMemberPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateMemberPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid argument for password unmatch",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.UpdateMemberPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "Passw0rd",
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
					NewPassword: "Passw0rd",
				}
				mocks.userAuth.EXPECT().ChangePassword(ctx, params).Return(assert.AnError)
			},
			input: &user.UpdateMemberPasswordInput{
				AccessToken:          "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ",
				OldPassword:          "12345678",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateMemberPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestForgotMemberPassword(t *testing.T) {
	t.Parallel()

	m := &entity.Member{CognitoID: "cognito-id"}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.ForgotMemberPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(m, nil)
				mocks.userAuth.EXPECT().ForgotPassword(ctx, "cognito-id").Return(nil)
			},
			input: &user.ForgotMemberPasswordInput{
				Email: "test-user@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ForgotMemberPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(nil, assert.AnError)
			},
			input: &user.ForgotMemberPasswordInput{
				Email: "test-user@and-period.jp",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to forget password",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(m, nil)
				mocks.userAuth.EXPECT().ForgotPassword(ctx, "cognito-id").Return(assert.AnError)
			},
			input: &user.ForgotMemberPasswordInput{
				Email: "test-user@and-period.jp",
			},
			expectErr: exception.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ForgotMemberPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestVerifyMemberPassword(t *testing.T) {
	t.Parallel()

	m := &entity.Member{CognitoID: "cognito-id"}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.VerifyMemberPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmForgotPasswordParams{
					Username:    "cognito-id",
					VerifyCode:  "123456",
					NewPassword: "Passw0rd",
				}
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(m, nil)
				mocks.userAuth.EXPECT().ConfirmForgotPassword(ctx, params).Return(nil)
			},
			input: &user.VerifyMemberPasswordInput{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.VerifyMemberPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.VerifyMemberPasswordInput{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "123456789",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get by email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(nil, assert.AnError)
			},
			input: &user.VerifyMemberPasswordInput{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to confirm forgot password",
			setup: func(ctx context.Context, mocks *mocks) {
				params := &cognito.ConfirmForgotPasswordParams{
					Username:    "cognito-id",
					VerifyCode:  "123456",
					NewPassword: "Passw0rd",
				}
				mocks.db.Member.EXPECT().GetByEmail(ctx, "test-user@and-period.jp", "cognito_id").Return(m, nil)
				mocks.userAuth.EXPECT().ConfirmForgotPassword(ctx, params).Return(assert.AnError)
			},
			input: &user.VerifyMemberPasswordInput{
				Email:                "test-user@and-period.jp",
				VerifyCode:           "123456",
				NewPassword:          "Passw0rd",
				PasswordConfirmation: "Passw0rd",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyMemberPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateMemberThumbnails(t *testing.T) {
	t.Parallel()

	thumbnails := common.Images{
		{
			Size: common.ImageSizeSmall,
			URL:  "https://and-period.jp/thumbnail_240.png",
		},
		{
			Size: common.ImageSizeMedium,
			URL:  "https://and-period.jp/thumbnail_675.png",
		},
		{
			Size: common.ImageSizeLarge,
			URL:  "https://and-period.jp/thumbnail_900.png",
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateMemberThumbnailsInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().UpdateThumbnails(ctx, "user-id", thumbnails).Return(nil)
			},
			input: &user.UpdateMemberThumbnailsInput{
				UserID:     "user-id",
				Thumbnails: thumbnails,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateMemberThumbnailsInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update thumbnails",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().UpdateThumbnails(ctx, "user-id", thumbnails).Return(assert.AnError)
			},
			input: &user.UpdateMemberThumbnailsInput{
				UserID:     "user-id",
				Thumbnails: thumbnails,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateMemberThumbnails(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
