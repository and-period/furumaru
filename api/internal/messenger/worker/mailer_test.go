package worker

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/stretchr/testify/assert"
)

func TestMultiSendMail(t *testing.T) {
	t.Parallel()

	in := &user.MultiGetAdminsInput{
		AdminIDs: []string{"admin-id"},
	}
	admins := uentity.Admins{
		{
			Lastname:  "&.",
			Firstname: "スタッフ",
			Email:     "test-user@and-period.jp",
		},
	}
	personalizations := []*mailer.Personalization{
		{
			Name:    "&. スタッフ",
			Address: "test-user@and-period.jp",
			Type:    mailer.AddressTypeTo,
			Substitutions: map[string]interface{}{
				"key": "value",
				"氏名":  "&. スタッフ",
			},
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.WorkerPayload
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdmins(ctx, in).Return(admins, nil)
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, entity.EmailIDAdminRegister, personalizations).Return(nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expectErr: nil,
		},
		{
			name: "failed to new personalizations",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdmins(ctx, in).Return(nil, errmock)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expectErr: errmock,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.multiSendMail(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestSendMail(t *testing.T) {
	t.Parallel()

	personalizations := []*mailer.Personalization{
		{
			Name:    "&. 農園",
			Address: "test-admin@and-period.jp",
			Type:    mailer.AddressTypeTo,
			Substitutions: map[string]interface{}{
				"氏名":     "&. 農園",
				"パスワード":  "!Qaz2wsx",
				"サイトURL": "https://admin.and-period.jp/signin",
			},
		},
	}

	tests := []struct {
		name             string
		setup            func(ctx context.Context, mocks *mocks)
		emailID          string
		personalizations []*mailer.Personalization
		expectErr        error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, "email-id", personalizations).Return(nil)
			},
			emailID:          "email-id",
			personalizations: personalizations,
			expectErr:        nil,
		},
		{
			name:             "personalizations is empty",
			setup:            func(ctx context.Context, mocks *mocks) {},
			emailID:          "email-id",
			personalizations: nil,
			expectErr:        nil,
		},
		{
			name: "failed to send info mail",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, "email-id", personalizations).Return(errmock)
			},
			emailID:          "email-id",
			personalizations: personalizations,
			expectErr:        exception.ErrUnknown,
		},
		{
			name: "failed to send info mail with retry",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.mailer.EXPECT().
					MultiSendFromInfo(ctx, "email-id", personalizations).
					Return(mailer.ErrUnavailable).Times(2)
			},
			emailID:          "email-id",
			personalizations: personalizations,
			expectErr:        exception.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.sendMail(ctx, tt.emailID, tt.personalizations...)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestPersonalizations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.WorkerPayload
		expect    []*mailer.Personalization
		expectErr error
	}{
		{
			name: "success admins",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetAdminsInput{AdminIDs: []string{"admin-id"}}
				admins := uentity.Admins{{
					Lastname:  "&.",
					Firstname: "スタッフ",
					Email:     "test-user@and-period.jp",
				}}
				mocks.user.EXPECT().MultiGetAdmins(ctx, in).Return(admins, nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expect: []*mailer.Personalization{
				{
					Name:    "&. スタッフ",
					Address: "test-user@and-period.jp",
					Type:    mailer.AddressTypeTo,
					Substitutions: map[string]interface{}{
						"key": "value",
						"氏名":  "&. スタッフ",
					},
				},
			},
			expectErr: nil,
		},
		{
			name: "success administrators",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetAdministratorsInput{AdministratorIDs: []string{"admin-id"}}
				administrators := uentity.Administrators{{
					Lastname:  "&.",
					Firstname: "スタッフ",
					Email:     "test-user@and-period.jp",
				}}
				mocks.user.EXPECT().MultiGetAdministrators(ctx, in).Return(administrators, nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeAdministrator,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expect: []*mailer.Personalization{
				{
					Name:    "&. スタッフ",
					Address: "test-user@and-period.jp",
					Type:    mailer.AddressTypeTo,
					Substitutions: map[string]interface{}{
						"key": "value",
						"氏名":  "&. スタッフ",
					},
				},
			},
			expectErr: nil,
		},
		{
			name: "success coordinators",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetCoordinatorsInput{CoordinatorIDs: []string{"admin-id"}}
				coordinators := uentity.Coordinators{{
					Lastname:  "&.",
					Firstname: "スタッフ",
					Email:     "test-user@and-period.jp",
				}}
				mocks.user.EXPECT().MultiGetCoordinators(ctx, in).Return(coordinators, nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeCoordinator,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expect: []*mailer.Personalization{
				{
					Name:    "&. スタッフ",
					Address: "test-user@and-period.jp",
					Type:    mailer.AddressTypeTo,
					Substitutions: map[string]interface{}{
						"key": "value",
						"氏名":  "&. スタッフ",
					},
				},
			},
			expectErr: nil,
		},
		{
			name: "success producers",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetProducersInput{ProducerIDs: []string{"admin-id"}}
				producers := uentity.Producers{{
					Lastname:  "&.",
					Firstname: "スタッフ",
					Email:     "test-user@and-period.jp",
				}}
				mocks.user.EXPECT().MultiGetProducers(ctx, in).Return(producers, nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeProducer,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expect: []*mailer.Personalization{
				{
					Name:    "&. スタッフ",
					Address: "test-user@and-period.jp",
					Type:    mailer.AddressTypeTo,
					Substitutions: map[string]interface{}{
						"key": "value",
						"氏名":  "&. スタッフ",
					},
				},
			},
			expectErr: nil,
		},
		{
			name: "success users",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.MultiGetUsersInput{UserIDs: []string{"user-id"}}
				users := uentity.Users{
					{Username: "&. スタッフ", Email: "test-user@and-period.jp"},
					{Username: "&. スタッフ", Email: ""},
				}
				mocks.user.EXPECT().MultiGetUsers(ctx, in).Return(users, nil)
			},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expect: []*mailer.Personalization{
				{
					Name:    "&. スタッフ",
					Address: "test-user@and-period.jp",
					Type:    mailer.AddressTypeTo,
					Substitutions: map[string]interface{}{
						"key": "value",
						"氏名":  "&. スタッフ",
					},
				},
			},
			expectErr: nil,
		},
		{
			name:  "success guest",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeGuest,
				UserIDs:   []string{},
				Guest: &entity.Guest{
					Name:  "&. スタッフ",
					Email: "test-user@and-period.jp",
				},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expect: []*mailer.Personalization{
				{
					Name:    "&. スタッフ",
					Address: "test-user@and-period.jp",
					Type:    mailer.AddressTypeTo,
					Substitutions: map[string]interface{}{
						"key": "value",
						"氏名":  "&. スタッフ",
					},
				},
			},
			expectErr: nil,
		},
		{
			name:  "failed to invalid user type",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &entity.WorkerPayload{
				EventType: entity.EventTypeAdminRegister,
				UserType:  entity.UserTypeNone,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expectErr: errUnknownUserType,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			actual, err := worker.newPersonalizations(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestFetchAdmins(t *testing.T) {
	t.Parallel()

	in := &user.MultiGetAdminsInput{
		AdminIDs: []string{"admin-id"},
	}
	admins := uentity.Admins{
		{
			ID:            "admin-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			CreatedAt:     jst.Date(2022, 7, 10, 18, 30, 0, 0),
			UpdatedAt:     jst.Date(2022, 7, 10, 18, 30, 0, 0),
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		adminIDs  []string
		execute   func(t *testing.T) func(name, email string)
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdmins(ctx, in).Return(admins, nil)
			},
			adminIDs: []string{"admin-id"},
			execute: func(t *testing.T) func(name, email string) {
				execute := func(name, email string) {
					assert.Equal(t, "&. スタッフ", name)
					assert.Equal(t, "test-admin@and-period.jp", email)
				}
				return execute
			},
			expectErr: nil,
		},
		{
			name: "failed to get admins",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdmins(ctx, in).Return(nil, errmock)
			},
			adminIDs: []string{"admin-id"},
			execute: func(t *testing.T) func(name, email string) {
				return nil
			},
			expectErr: errmock,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.fetchAdmins(ctx, tt.adminIDs, tt.execute(t))
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestFetchAdministrators(t *testing.T) {
	t.Parallel()

	in := &user.MultiGetAdministratorsInput{
		AdministratorIDs: []string{"admin-id"},
	}
	administrators := uentity.Administrators{
		{
			ID:            "administrator-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			Email:         "test-admin@and-period.jp",
			PhoneNumber:   "+819012345678",
			CreatedAt:     jst.Date(2022, 7, 10, 18, 30, 0, 0),
			UpdatedAt:     jst.Date(2022, 7, 10, 18, 30, 0, 0),
		},
	}

	tests := []struct {
		name             string
		setup            func(ctx context.Context, mocks *mocks)
		administratorIDs []string
		execute          func(t *testing.T) func(name, email string)
		expectErr        error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdministrators(ctx, in).Return(administrators, nil)
			},
			administratorIDs: []string{"admin-id"},
			execute: func(t *testing.T) func(name, email string) {
				execute := func(name, email string) {
					assert.Equal(t, "&. スタッフ", name)
					assert.Equal(t, "test-admin@and-period.jp", email)
				}
				return execute
			},
			expectErr: nil,
		},
		{
			name: "failed to get administrators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetAdministrators(ctx, in).Return(nil, errmock)
			},
			administratorIDs: []string{"admin-id"},
			execute: func(t *testing.T) func(name, email string) {
				return nil
			},
			expectErr: errmock,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.fetchAdministrators(ctx, tt.administratorIDs, tt.execute(t))
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestFetchCoordinators(t *testing.T) {
	t.Parallel()

	in := &user.MultiGetCoordinatorsInput{
		CoordinatorIDs: []string{"admin-id"},
	}
	coordinators := uentity.Coordinators{
		{
			ID:               "coordinator-id",
			Lastname:         "&.",
			Firstname:        "スタッフ",
			LastnameKana:     "あんどぴりおど",
			FirstnameKana:    "すたっふ",
			StoreName:        "&.農園",
			ThumbnailURL:     "https://and-period.jp/thumbnail.png",
			HeaderURL:        "https://and-period.jp/header.png",
			TwitterAccount:   "twitter-account",
			InstagramAccount: "instagram-account",
			FacebookAccount:  "facebook-account",
			Email:            "test-admin@and-period.jp",
			PhoneNumber:      "+819012345678",
			PostalCode:       "1000014",
			Prefecture:       "東京都",
			City:             "千代田区",
			AddressLine1:     "永田町1-7-1",
			AddressLine2:     "",
			CreatedAt:        jst.Date(2022, 7, 10, 18, 30, 0, 0),
			UpdatedAt:        jst.Date(2022, 7, 10, 18, 30, 0, 0),
		},
	}

	tests := []struct {
		name           string
		setup          func(ctx context.Context, mocks *mocks)
		coordinatorIDs []string
		execute        func(t *testing.T) func(name, email string)
		expectErr      error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetCoordinators(ctx, in).Return(coordinators, nil)
			},
			coordinatorIDs: []string{"admin-id"},
			execute: func(t *testing.T) func(name, email string) {
				execute := func(name, email string) {
					assert.Equal(t, "&. スタッフ", name)
					assert.Equal(t, "test-admin@and-period.jp", email)
				}
				return execute
			},
			expectErr: nil,
		},
		{
			name: "failed to get coordinators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetCoordinators(ctx, in).Return(nil, errmock)
			},
			coordinatorIDs: []string{"admin-id"},
			execute: func(t *testing.T) func(name, email string) {
				return nil
			},
			expectErr: errmock,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.fetchCoordinators(ctx, tt.coordinatorIDs, tt.execute(t))
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestFetchProducers(t *testing.T) {
	t.Parallel()

	in := &user.MultiGetProducersInput{
		ProducerIDs: []string{"admin-id"},
	}
	producers := uentity.Producers{
		{
			ID:            "admin-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "すたっふ",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			Email:         "test-admin@and-period.jp",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			AddressLine1:  "永田町1-7-1",
			AddressLine2:  "",
			CreatedAt:     jst.Date(2022, 7, 10, 18, 30, 0, 0),
			UpdatedAt:     jst.Date(2022, 7, 10, 18, 30, 0, 0),
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		producerIDs []string
		execute     func(t *testing.T) func(name, email string)
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetProducers(ctx, in).Return(producers, nil)
			},
			producerIDs: []string{"admin-id"},
			execute: func(t *testing.T) func(name, email string) {
				execute := func(name, email string) {
					assert.Equal(t, "&. スタッフ", name)
					assert.Equal(t, "test-admin@and-period.jp", email)
				}
				return execute
			},
			expectErr: nil,
		},
		{
			name: "failed to get producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetProducers(ctx, in).Return(nil, errmock)
			},
			producerIDs: []string{"admin-id"},
			execute: func(t *testing.T) func(name, email string) {
				return nil
			},
			expectErr: errmock,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.fetchProducers(ctx, tt.producerIDs, tt.execute(t))
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestFetchUsers(t *testing.T) {
	t.Parallel()

	in := &user.MultiGetUsersInput{
		UserIDs: []string{"user-id"},
	}
	users := uentity.Users{
		{
			ID:           "user-id",
			AccountID:    "account-id",
			CognitoID:    "cognito-id",
			Username:     "テストユーザー",
			ProviderType: uentity.ProviderTypeEmail,
			Email:        "test-user@and-period.jp",
			PhoneNumber:  "+810000000000",
			ThumbnailURL: "https://and-period.jp/thumbnail.png",
			CreatedAt:    jst.Date(2022, 7, 10, 18, 30, 0, 0),
			UpdatedAt:    jst.Date(2022, 7, 10, 18, 30, 0, 0),
			VerifiedAt:   jst.Date(2022, 7, 10, 18, 30, 0, 0),
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		userIDs   []string
		execute   func(t *testing.T) func(name, email string)
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetUsers(ctx, in).Return(users, nil)
			},
			userIDs: []string{"user-id"},
			execute: func(t *testing.T) func(name, email string) {
				execute := func(name, email string) {
					assert.Equal(t, "テストユーザー", name)
					assert.Equal(t, "test-user@and-period.jp", email)
				}
				return execute
			},
			expectErr: nil,
		},
		{
			name: "failed to get users",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().MultiGetUsers(ctx, in).Return(nil, errmock)
			},
			userIDs: []string{"user-id"},
			execute: func(t *testing.T) func(name, email string) {
				return nil
			},
			expectErr: errmock,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.fetchUsers(ctx, tt.userIDs, tt.execute(t))
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestFetchGuest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		guest     *entity.Guest
		execute   func(t *testing.T) func(name, email string)
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			guest: &entity.Guest{
				Name:  "テストユーザー",
				Email: "test-user@and-period.jp",
			},
			execute: func(t *testing.T) func(name, email string) {
				execute := func(name, email string) {
					assert.Equal(t, "テストユーザー", name)
					assert.Equal(t, "test-user@and-period.jp", email)
				}
				return execute
			},
			expectErr: nil,
		},
		{
			name:  "guest is empty",
			setup: func(ctx context.Context, mocks *mocks) {},
			guest: nil,
			execute: func(t *testing.T) func(name, email string) {
				return nil
			},
			expectErr: errGuestRequired,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.fetchGuest(ctx, tt.guest, tt.execute(t))
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
