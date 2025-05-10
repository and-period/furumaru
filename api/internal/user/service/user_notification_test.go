package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestMultiGetUserNotifications(t *testing.T) {
	t.Parallel()
	now := time.Now()
	notifications := entity.UserNotifications{
		{
			UserID:    "user-id",
			Disabled:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetUserNotificationsInput
		expect    entity.UserNotifications
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.UserNotification.EXPECT().MultiGet(ctx, []string{"user-id"}).Return(notifications, nil)
			},
			input: &user.MultiGetUserNotificationsInput{
				UserIDs: []string{"user-id"},
			},
			expect:    notifications,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetUserNotificationsInput{
				UserIDs: []string{""},
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get notifications",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.UserNotification.EXPECT().MultiGet(ctx, []string{"user-id"}).Return(nil, assert.AnError)
			},
			input: &user.MultiGetUserNotificationsInput{
				UserIDs: []string{"user-id"},
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetUserNotifications(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetUserNotification(t *testing.T) {
	t.Parallel()
	now := time.Now()
	notification := &entity.UserNotification{
		UserID:    "user-id",
		Disabled:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetUserNotificationInput
		expect    *entity.UserNotification
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.UserNotification.EXPECT().Get(ctx, "user-id").Return(notification, nil)
			},
			input: &user.GetUserNotificationInput{
				UserID: "user-id",
			},
			expect:    notification,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetUserNotificationInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.UserNotification.EXPECT().Get(ctx, "user-id").Return(nil, assert.AnError)
			},
			input: &user.GetUserNotificationInput{
				UserID: "user-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetUserNotification(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUpdateUserNotification(t *testing.T) {
	t.Parallel()
	now := time.Now()
	notification := &entity.UserNotification{
		UserID:    "user-id",
		Disabled:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateUserNotificationInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.UserNotification.EXPECT().Get(ctx, "user-id").Return(notification, nil)
				mocks.db.UserNotification.EXPECT().Upsert(ctx, notification).Return(nil)
			},
			input: &user.UpdateUserNotificationInput{
				UserID:  "user-id",
				Enabled: true,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateUserNotificationInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.UserNotification.EXPECT().Get(ctx, "user-id").Return(nil, assert.AnError)
			},
			input: &user.UpdateUserNotificationInput{
				UserID:  "user-id",
				Enabled: true,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to upsert user notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.UserNotification.EXPECT().Get(ctx, "user-id").Return(notification, nil)
				mocks.db.UserNotification.EXPECT().Upsert(ctx, notification).Return(assert.AnError)
			},
			input: &user.UpdateUserNotificationInput{
				UserID:  "user-id",
				Enabled: true,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateUserNotification(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
