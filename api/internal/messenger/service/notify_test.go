package service

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/internal/messenger/entity"
	"github.com/and-period/marche/api/pkg/mailer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMessengerService_SendInfoMail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		setup            func(ctx context.Context, t *testing.T, mocks *mocks)
		message          *entity.MailConfig
		personalizations []*mailer.Personalization
		expect           bool
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, "email-id", gomock.Any()).Return(nil)
			},
			message: &entity.MailConfig{
				EmailID:       "email-id",
				Substitutions: map[string]string{},
			},
			personalizations: []*mailer.Personalization{
				{
					Name:          "送信先名",
					Address:       "test@and-period.jp",
					Type:          mailer.AddressTypeTo,
					Substitutions: map[string]interface{}{},
				},
			},
			expect: false,
		},
		{
			name: "failed to send mail",
			setup: func(ctx context.Context, t *testing.T, mocks *mocks) {
				mocks.mailer.EXPECT().MultiSendFromInfo(ctx, "email-id", gomock.Any()).Return(mailer.ErrTimeout).Times(2)
			},
			message: &entity.MailConfig{
				EmailID:       "email-id",
				Substitutions: map[string]string{},
			},
			personalizations: []*mailer.Personalization{
				{
					Name:          "送信先名",
					Address:       "test@and-period.jp",
					Type:          mailer.AddressTypeTo,
					Substitutions: map[string]interface{}{},
				},
			},
			expect: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mocks := newMocks(ctrl)
			srv := &messengerService{mailer: mocks.mailer, maxRetries: 1}
			tt.setup(ctx, t, mocks)
			err := srv.sendInfoMail(ctx, tt.message, tt.personalizations...)
			assert.Equal(t, tt.expect, err != nil, err)
		})
	}
}
