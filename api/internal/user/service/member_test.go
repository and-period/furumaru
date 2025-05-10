package service

import (
	"context"
	"testing"
	"time"

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
						ProviderType:  entity.UserAuthProviderTypeEmail,
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
						expectUser.Member.UserID, expectUser.CognitoID = u.ID, u.CognitoID
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
					ProviderType:  entity.UserAuthProviderTypeEmail,
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
					ProviderType:  entity.UserAuthProviderTypeEmail,
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

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyMember(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateMemberEmail(t *testing.T) {
	t.Parallel()

	m := &entity.Member{
		ProviderType: entity.UserAuthProviderTypeEmail,
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
					ProviderType: entity.UserAuthProviderTypeGoogle,
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

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.VerifyMemberPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateMemberUsername(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateMemberUsernameInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().UpdateUsername(ctx, "user-id", "username").Return(nil)
			},
			input: &user.UpdateMemberUsernameInput{
				UserID:   "user-id",
				Username: "username",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateMemberUsernameInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update username",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().UpdateUsername(ctx, "user-id", "username").Return(assert.AnError)
			},
			input: &user.UpdateMemberUsernameInput{
				UserID:   "user-id",
				Username: "username",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateMemberUsername(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateMemberAccountID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateMemberAccountIDInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().UpdateAccountID(ctx, "user-id", "account-id").Return(nil)
			},
			input: &user.UpdateMemberAccountIDInput{
				UserID:    "user-id",
				AccountID: "account-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateMemberAccountIDInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update account id",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().UpdateAccountID(ctx, "user-id", "account-id").Return(assert.AnError)
			},
			input: &user.UpdateMemberAccountIDInput{
				UserID:    "user-id",
				AccountID: "account-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateMemberAccountID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateMemberThumbnailURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateMemberThumbnailURLInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().UpdateThumbnailURL(ctx, "user-id", "http://example.com/thumbnail.png").Return(nil)
			},
			input: &user.UpdateMemberThumbnailURLInput{
				UserID:       "user-id",
				ThumbnailURL: "http://example.com/thumbnail.png",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateMemberThumbnailURLInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update thumbnail url",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Member.EXPECT().UpdateThumbnailURL(ctx, "user-id", "http://example.com/thumbnail.png").Return(assert.AnError)
			},
			input: &user.UpdateMemberThumbnailURLInput{
				UserID:       "user-id",
				ThumbnailURL: "http://example.com/thumbnail.png",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateMemberThumbnailURL(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestAuthMemberWithGoogle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.AuthMemberWithGoogleInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mockAuthMemberWithOAuth(mocks, t, entity.UserAuthProviderTypeGoogle)
			},
			input: &user.AuthMemberWithGoogleInput{
				AuthMemberDetailWithOAuth: user.AuthMemberDetailWithOAuth{
					SessionID:   "session-id",
					State:       "state",
					RedirectURI: "http://example.com/auth/external/callback",
				},
			},
			expect:    "http://example.com/auth/eternal/callback",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.AuthMemberWithGoogleInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			authURL, err := service.AuthMemberWithGoogle(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, authURL)
		}))
	}
}

func TestCreateMemberWithGoogle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateMemberWithGoogleInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mockCreateMemberWithOAuth(mocks, t, entity.UserAuthProviderTypeGoogle)
			},
			input: &user.CreateMemberWithGoogleInput{
				CreateMemberDetailWithOAuth: user.CreateMemberDetailWithOAuth{
					SessionID:     "session-id",
					Code:          "code",
					Nonce:         "nonce",
					RedirectURI:   "http://example.com/auth/external/callback",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					PhoneNumber:   "+810000000000",
				},
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateMemberWithGoogleInput{},
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateMemberWithGoogle(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestAuthMemberWithLINE(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.AuthMemberWithLINEInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mockAuthMemberWithOAuth(mocks, t, entity.UserAuthProviderTypeLINE)
			},
			input: &user.AuthMemberWithLINEInput{
				AuthMemberDetailWithOAuth: user.AuthMemberDetailWithOAuth{
					SessionID:   "session-id",
					State:       "state",
					RedirectURI: "http://example.com/auth/external/callback",
				},
			},
			expect:    "http://example.com/auth/eternal/callback",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.AuthMemberWithLINEInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			authURL, err := service.AuthMemberWithLINE(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, authURL)
		}))
	}
}

func TestCreateMemberWithLINE(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateMemberWithLINEInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mockCreateMemberWithOAuth(mocks, t, entity.UserAuthProviderTypeLINE)
			},
			input: &user.CreateMemberWithLINEInput{
				CreateMemberDetailWithOAuth: user.CreateMemberDetailWithOAuth{
					SessionID:     "session-id",
					Code:          "code",
					Nonce:         "nonce",
					RedirectURI:   "http://example.com/auth/external/callback",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					PhoneNumber:   "+810000000000",
				},
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateMemberWithLINEInput{},
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateMemberWithLINE(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func mockAuthMemberWithOAuth(m *mocks, t *testing.T, providerType entity.UserAuthProviderType) {
	params := &cognito.GenerateAuthURLParams{
		State:        "state",
		Nonce:        "nonce",
		ProviderType: providerType.ToCognito(),
		RedirectURI:  "http://example.com/auth/external/callback",
	}

	m.cache.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
	m.userAuth.EXPECT().
		GenerateAuthURL(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, p *cognito.GenerateAuthURLParams) (string, error) {
			p.Nonce = "nonce"
			assert.Equal(t, params, p)
			return "http://example.com/auth/eternal/callback", nil
		})
}

func TestAuthMemberWithOAuth(t *testing.T) {
	t.Parallel()

	params := &cognito.GenerateAuthURLParams{
		State:        "state",
		Nonce:        "nonce",
		ProviderType: cognito.ProviderTypeGoogle,
		RedirectURI:  "http://example.com/auth/google/callback",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *authMemberWithOAuthParams
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
				mocks.userAuth.EXPECT().
					GenerateAuthURL(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, p *cognito.GenerateAuthURLParams) (string, error) {
						p.Nonce = "nonce"
						assert.Equal(t, params, p)
						return "http://example.com/auth/google", nil
					})
			},
			input: &authMemberWithOAuthParams{
				payload: &user.AuthMemberDetailWithOAuth{
					State:       "state",
					RedirectURI: "http://example.com/auth/google/callback",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expect:    "http://example.com/auth/google",
			expectErr: nil,
		},
		{
			name: "failed to insert cache",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().Insert(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &authMemberWithOAuthParams{
				payload: &user.AuthMemberDetailWithOAuth{
					State:       "state",
					RedirectURI: "http://example.com/auth/google/callback",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to generate auth url",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
				mocks.userAuth.EXPECT().GenerateAuthURL(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &authMemberWithOAuthParams{
				payload: &user.AuthMemberDetailWithOAuth{
					State:       "state",
					RedirectURI: "http://example.com/auth/google/callback",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			authURL, err := service.authMemberWithOAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, authURL)
		}))
	}
}

func mockCreateMemberWithOAuth(m *mocks, t *testing.T, providerType entity.UserAuthProviderType) {
	tokenParams := &cognito.GetAccessTokenParams{
		Code:        "code",
		RedirectURI: "http://example.com/auth/external/callback",
	}
	token := &cognito.AuthResult{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
		ExpiresIn:    3600,
	}
	authUser := &cognito.AuthUser{
		Username:    "external_username",
		Email:       "test@example.com",
		PhoneNumber: "",
		Identities: []*cognito.AuthUserIdentity{
			{
				UserID:       "username",
				ProviderType: providerType.ToCognito(),
				Primary:      true,
				DateCreated:  0,
			},
		},
	}
	expectUser := &entity.User{
		Registered: true,
		Member: entity.Member{
			Username:      "username",
			CognitoID:     "external_username",
			AccountID:     "account-id",
			Lastname:      "&.",
			Firstname:     "利用者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "りようしゃ",
			ProviderType:  providerType,
			Email:         "test@example.com",
			PhoneNumber:   "+810000000000",
		},
	}

	m.cache.EXPECT().
		Get(gomock.Any(), &entity.UserAuthEvent{SessionID: "session-id"}).
		DoAndReturn(func(ctx context.Context, event *entity.UserAuthEvent) error {
			event.ProviderType = providerType
			event.Nonce = "nonce"
			return nil
		})
	m.userAuth.EXPECT().GetAccessToken(gomock.Any(), tokenParams).Return(token, nil)
	m.userAuth.EXPECT().GetUser(gomock.Any(), "access-token").Return(authUser, nil)
	m.db.Member.EXPECT().
		Create(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *entity.User, auth func(ctx context.Context) error) error {
			expectUser.ID = u.ID
			expectUser.Member.UserID = u.ID
			assert.Equal(t, expectUser, u)
			return auth(ctx)
		})
}

func TestCreateMemberWithOAuth(t *testing.T) {
	t.Parallel()

	tokenParams := &cognito.GetAccessTokenParams{
		Code:        "code",
		RedirectURI: "http://example.com/auth/google/callback",
	}
	token := &cognito.AuthResult{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
		ExpiresIn:    3600,
	}
	authUser := &cognito.AuthUser{
		Username:    "google_username",
		Email:       "test@example.com",
		PhoneNumber: "",
		Identities: []*cognito.AuthUserIdentity{
			{
				UserID:       "username",
				ProviderType: cognito.ProviderTypeGoogle,
				Primary:      true,
				DateCreated:  0,
			},
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *createMemberWithOAuthParams
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				expectUser := &entity.User{
					Registered: true,
					Member: entity.Member{
						Username:      "username",
						CognitoID:     "google_username",
						AccountID:     "account-id",
						Lastname:      "&.",
						Firstname:     "利用者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "りようしゃ",
						ProviderType:  entity.UserAuthProviderTypeGoogle,
						Email:         "test@example.com",
						PhoneNumber:   "+810000000000",
					},
				}
				mocks.cache.EXPECT().
					Get(ctx, &entity.UserAuthEvent{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.UserAuthEvent) error {
						event.ProviderType = entity.UserAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.userAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.userAuth.EXPECT().GetUser(ctx, "access-token").Return(authUser, nil)
				mocks.db.Member.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *entity.User, auth func(ctx context.Context) error) error {
						expectUser.ID = u.ID
						expectUser.Member.UserID = u.ID
						assert.Equal(t, expectUser, u)
						return auth(ctx)
					})
			},
			input: &createMemberWithOAuthParams{
				payload: &user.CreateMemberDetailWithOAuth{
					SessionID:     "session-id",
					Code:          "code",
					Nonce:         "nonce",
					RedirectURI:   "http://example.com/auth/google/callback",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					PhoneNumber:   "+810000000000",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expectErr: nil,
		},
		{
			name: "unmatch nonce",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.UserAuthEvent{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.UserAuthEvent) error {
						event.ProviderType = entity.UserAuthProviderTypeGoogle
						event.Nonce = "invalid-nonce"
						return nil
					})
			},
			input: &createMemberWithOAuthParams{
				payload: &user.CreateMemberDetailWithOAuth{
					SessionID:     "session-id",
					Code:          "code",
					Nonce:         "nonce",
					RedirectURI:   "http://example.com/auth/google/callback",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					PhoneNumber:   "+810000000000",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "unmatch provider type",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.UserAuthEvent{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.UserAuthEvent) error {
						event.ProviderType = entity.UserAuthProviderTypeLINE
						event.Nonce = "nonce"
						return nil
					})
			},
			input: &createMemberWithOAuthParams{
				payload: &user.CreateMemberDetailWithOAuth{
					SessionID:     "session-id",
					Code:          "code",
					Nonce:         "nonce",
					RedirectURI:   "http://example.com/auth/google/callback",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					PhoneNumber:   "+810000000000",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to get auth user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.UserAuthEvent{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.UserAuthEvent) error {
						event.ProviderType = entity.UserAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.userAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(nil, assert.AnError)
			},
			input: &createMemberWithOAuthParams{
				payload: &user.CreateMemberDetailWithOAuth{
					SessionID:     "session-id",
					Code:          "code",
					Nonce:         "nonce",
					RedirectURI:   "http://example.com/auth/google/callback",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					PhoneNumber:   "+810000000000",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get auth user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.UserAuthEvent{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.UserAuthEvent) error {
						event.ProviderType = entity.UserAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.userAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.userAuth.EXPECT().GetUser(ctx, "access-token").Return(nil, assert.AnError)
			},
			input: &createMemberWithOAuthParams{
				payload: &user.CreateMemberDetailWithOAuth{
					SessionID:     "session-id",
					Code:          "code",
					Nonce:         "nonce",
					RedirectURI:   "http://example.com/auth/google/callback",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					PhoneNumber:   "+810000000000",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create member",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.cache.EXPECT().
					Get(ctx, &entity.UserAuthEvent{SessionID: "session-id"}).
					DoAndReturn(func(ctx context.Context, event *entity.UserAuthEvent) error {
						event.ProviderType = entity.UserAuthProviderTypeGoogle
						event.Nonce = "nonce"
						return nil
					})
				mocks.userAuth.EXPECT().GetAccessToken(ctx, tokenParams).Return(token, nil)
				mocks.userAuth.EXPECT().GetUser(ctx, "access-token").Return(authUser, nil)
				mocks.db.Member.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &createMemberWithOAuthParams{
				payload: &user.CreateMemberDetailWithOAuth{
					SessionID:     "session-id",
					Code:          "code",
					Nonce:         "nonce",
					RedirectURI:   "http://example.com/auth/google/callback",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					PhoneNumber:   "+810000000000",
				},
				providerType: entity.UserAuthProviderTypeGoogle,
				redirectURI:  "http://example.com/auth/google/callback",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.createMemberWithOAuth(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
