package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/stretchr/testify/assert"
)

func TestNotifyRegisterAdmin(t *testing.T) {
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
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyRegisterAdminInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, "register-admin", personalizations).Return(nil)
			},
			input: &messenger.NotifyRegisterAdminInput{
				Name:     "&. 農園",
				Email:    "test-admin@and-period.jp",
				Password: "!Qaz2wsx",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyRegisterAdminInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send info mail",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, "register-admin", personalizations).Return(errmock)
			},
			input: &messenger.NotifyRegisterAdminInput{
				Name:     "&. 農園",
				Email:    "test-admin@and-period.jp",
				Password: "!Qaz2wsx",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyRegisterAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
