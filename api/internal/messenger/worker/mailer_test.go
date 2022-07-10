package worker

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/stretchr/testify/assert"
)

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
