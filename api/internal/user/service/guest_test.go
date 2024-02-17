package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUpsertGuest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpsertGuestInput
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.UpsertGuestInput{
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				Email:         "test@and-period.jp",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpsertGuestInput{},
			expectErr: exception.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.UpsertGuest(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
