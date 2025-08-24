package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpsertGuest(t *testing.T) {
	t.Parallel()
	now := time.Now()
	guest := &entity.Guest{
		UserID:        "user-id",
		Lastname:      "&.",
		Firstname:     "ゲスト",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "げすと",
		Email:         "test@example.com",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	params := &database.UpdateGuestParams{
		Lastname:      "&.",
		Firstname:     "利用者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "りようしゃ",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpsertGuestInput
		expectErr error
	}{
		{
			name: "success to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Guest.EXPECT().GetByEmail(ctx, "test@example.com").Return(nil, database.ErrNotFound)
				mocks.db.Guest.EXPECT().Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, user *entity.User) error {
						expect := &entity.User{
							ID:         user.ID, // ignore
							Type:       entity.UserTypeGuest,
							Registered: false,
							Guest: entity.Guest{
								UserID:        user.ID, // ignore
								Lastname:      "&.",
								Firstname:     "利用者",
								LastnameKana:  "あんどどっと",
								FirstnameKana: "りようしゃ",
								Email:         "test@example.com",
							},
						}
						assert.Equal(t, expect, user)
						return nil
					})
			},
			input: &user.UpsertGuestInput{
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				Email:         "test@example.com",
			},
			expectErr: nil,
		},
		{
			name: "success to udpate",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Guest.EXPECT().GetByEmail(ctx, "test@example.com").Return(guest, nil)
				mocks.db.Guest.EXPECT().Update(ctx, "user-id", params).Return(nil)
			},
			input: &user.UpsertGuestInput{
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				Email:         "test@example.com",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpsertGuestInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "success to get by email",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Guest.EXPECT().GetByEmail(ctx, "test@example.com").Return(nil, assert.AnError)
			},
			input: &user.UpsertGuestInput{
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				Email:         "test@example.com",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "success to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Guest.EXPECT().GetByEmail(ctx, "test@example.com").Return(nil, database.ErrNotFound)
				mocks.db.Guest.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &user.UpsertGuestInput{
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				Email:         "test@example.com",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "success to update",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Guest.EXPECT().GetByEmail(ctx, "test@example.com").Return(guest, nil)
				mocks.db.Guest.EXPECT().Update(ctx, "user-id", params).Return(assert.AnError)
			},
			input: &user.UpsertGuestInput{
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				Email:         "test@example.com",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.UpsertGuest(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
