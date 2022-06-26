package worker

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAdmin(t *testing.T) {
	t.Parallel()

	admin := &uentity.Admin{
		Lastname:      "&.",
		Firstname:     "農園",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "のうえん",
		Email:         "test-admin@and-period.jp",
	}
	ps := []*mailer.Personalization{
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
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *messenger.WorkerPayload
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.GetAdminInput{AdminID: "admin-id"}
				mocks.user.EXPECT().GetAdmin(ctx, in).Return(admin, nil)
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, "register-admin", ps).Return(nil)
			},
			payload: &messenger.WorkerPayload{
				EventType: messenger.EventTypeRegisterAdmin,
				UserType:  messenger.UserTypeAdministrator,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID: "register-admin",
					Substitutions: map[string]string{
						"パスワード": "!Qaz2wsx",
					},
				},
			},
			expectErr: nil,
		},
		{
			name:  "user id of payload is empty",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &messenger.WorkerPayload{
				EventType: messenger.EventTypeRegisterAdmin,
				UserType:  messenger.UserTypeAdministrator,
				UserIDs:   []string{},
				Email:     nil,
			},
			expectErr: errInvalidPayloadFormat,
		},
		{
			name:  "email of payload is empty",
			setup: func(ctx context.Context, mocks *mocks) {},
			payload: &messenger.WorkerPayload{
				EventType: messenger.EventTypeRegisterAdmin,
				UserType:  messenger.UserTypeAdministrator,
				UserIDs:   []string{"admin-id"},
				Email:     nil,
			},
			expectErr: errInvalidPayloadFormat,
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.GetAdminInput{AdminID: "admin-id"}
				mocks.user.EXPECT().GetAdmin(ctx, in).Return(nil, errmock)
			},
			payload: &messenger.WorkerPayload{
				EventType: messenger.EventTypeRegisterAdmin,
				UserType:  messenger.UserTypeAdministrator,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID: "register-admin",
					Substitutions: map[string]string{
						"パスワード": "!Qaz2wsx",
					},
				},
			},
			expectErr: errmock,
		},
		{
			name: "failed to send mail",
			setup: func(ctx context.Context, mocks *mocks) {
				in := &user.GetAdminInput{AdminID: "admin-id"}
				mocks.user.EXPECT().GetAdmin(ctx, in).Return(admin, nil)
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, "register-admin", ps).Return(errmock)
			},
			payload: &messenger.WorkerPayload{
				EventType: messenger.EventTypeRegisterAdmin,
				UserType:  messenger.UserTypeAdministrator,
				UserIDs:   []string{"admin-id"},
				Email: &entity.MailConfig{
					EmailID: "register-admin",
					Substitutions: map[string]string{
						"パスワード": "!Qaz2wsx",
					},
				},
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.registerAdmin(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
