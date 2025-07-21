package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListNotificaitons(t *testing.T) {
	t.Parallel()

	since := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	until := jst.Date(2022, 6, 30, 18, 30, 0, 0)
	params := &database.ListNotificationsParams{
		Limit:  30,
		Offset: 0,
		Since:  since,
		Until:  until,
		Orders: []*database.ListNotificationsOrder{
			{Key: database.ListNotificationsOrderByPublishedAt, OrderByASC: true},
		},
	}
	notifications := entity.Notifications{
		{
			ID:          "notification-id",
			Type:        entity.NotificationTypeSystem,
			Title:       "キャベツ祭り開催",
			Body:        "旬のキャベツを売り出します",
			PublishedAt: since,
			Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
			CreatedBy:   "admin-id",
			UpdatedBy:   "admin-id",
			CreatedAt:   since,
			UpdatedAt:   since,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *messenger.ListNotificationsInput
		expect      entity.Notifications
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().List(gomock.Any(), params).Return(notifications, nil)
				mocks.db.Notification.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &messenger.ListNotificationsInput{
				Limit:  30,
				Offset: 0,
				Since:  since,
				Until:  until,
				Orders: []*messenger.ListNotificationsOrder{
					{Key: messenger.ListNotificationsOrderByPublishedAt, OrderByASC: true},
				},
			},
			expect:      notifications,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &messenger.ListNotificationsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list notifications",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().
					List(gomock.Any(), params).
					Return(nil, assert.AnError)
				mocks.db.Notification.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &messenger.ListNotificationsInput{
				Limit:  30,
				Offset: 0,
				Since:  since,
				Until:  until,
				Orders: []*messenger.ListNotificationsOrder{
					{Key: messenger.ListNotificationsOrderByPublishedAt, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count notifications",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().List(gomock.Any(), params).Return(notifications, nil)
				mocks.db.Notification.EXPECT().
					Count(gomock.Any(), params).
					Return(int64(0), assert.AnError)
			},
			input: &messenger.ListNotificationsInput{
				Limit:  30,
				Offset: 0,
				Since:  since,
				Until:  until,
				Orders: []*messenger.ListNotificationsOrder{
					{Key: messenger.ListNotificationsOrderByPublishedAt, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, total, err := service.ListNotifications(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.ElementsMatch(t, tt.expect, actual)
				assert.Equal(t, tt.expectTotal, total)
			}),
		)
	}
}

func TestGetNotification(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	notification := &entity.Notification{
		ID:          "notification-id",
		Type:        entity.NotificationTypeSystem,
		Title:       "キャベツ祭り開催",
		Body:        "旬のキャベツを売り出します",
		PublishedAt: now,
		Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
		CreatedBy:   "admin-id",
		UpdatedBy:   "admin-id",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.GetNotificationInput
		expect    *entity.Notification
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
			},
			input: &messenger.GetNotificationInput{
				NotificationID: "notification-id",
			},
			expect:    notification,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.GetNotificationInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().
					Get(ctx, "notification-id").
					Return(nil, assert.AnError)
			},
			input: &messenger.GetNotificationInput{
				NotificationID: "notification-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.GetNotification(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}

func TestCreateNotification(t *testing.T) {
	t.Parallel()

	adminIn := &user.GetAdminInput{
		AdminID: "admin-id",
	}
	admin := &uentity.Admin{
		ID:        "admin-id",
		Firstname: "あんど",
		Lastname:  "ぴりおど",
	}
	promotionIn := &store.GetPromotionInput{
		PromotionID: "promotion-id",
	}
	promotion := &sentity.Promotion{
		ID: "promotion-id",
	}
	now := jst.Now()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.CreateNotificationInput
		expectErr error
	}{
		{
			name: "success system",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
				mocks.db.Notification.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, notification *entity.Notification) error {
						expect := &entity.Notification{
							ID:    notification.ID, // ignore
							Type:  entity.NotificationTypeSystem,
							Title: "キャベツ祭り開催",
							Body:  "旬のキャベツを売り出します",
							Note:  "",
							Targets: []entity.NotificationTarget{
								entity.NotificationTargetUsers,
							},
							PublishedAt: now.AddDate(0, 0, 1),
							PromotionID: "",
							CreatedBy:   "admin-id",
							UpdatedBy:   "admin-id",
						}
						assert.Equal(t, expect, notification)
						return nil
					})
				// 非同期関連
				mocks.db.Notification.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			},
			input: &messenger.CreateNotificationInput{
				Type:        entity.NotificationTypeSystem,
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Note:        "",
				Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				CreatedBy:   "admin-id",
				PromotionID: "",
			},
			expectErr: nil,
		},
		{
			name: "success promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
				mocks.store.EXPECT().GetPromotion(gomock.Any(), promotionIn).Return(promotion, nil)
				mocks.db.Notification.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, notification *entity.Notification) error {
						expect := &entity.Notification{
							ID:    notification.ID, // ignore
							Type:  entity.NotificationTypePromotion,
							Title: "",
							Body:  "旬のキャベツを売り出します",
							Note:  "",
							Targets: []entity.NotificationTarget{
								entity.NotificationTargetUsers,
							},
							PublishedAt: now.AddDate(0, 0, 1),
							PromotionID: "promotion-id",
							CreatedBy:   "admin-id",
							UpdatedBy:   "admin-id",
						}
						assert.Equal(t, expect, notification)
						return nil
					})
				// 非同期関連
				mocks.db.Notification.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			},
			input: &messenger.CreateNotificationInput{
				Type:        entity.NotificationTypePromotion,
				Title:       "",
				Body:        "旬のキャベツを売り出します",
				Note:        "",
				Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				CreatedBy:   "admin-id",
				PromotionID: "promotion-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.CreateNotificationInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "not found admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().
					GetAdmin(gomock.Any(), adminIn).
					Return(nil, exception.ErrNotFound)
			},
			input: &messenger.CreateNotificationInput{
				Type:        entity.NotificationTypeSystem,
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Note:        "",
				Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				CreatedBy:   "admin-id",
				PromotionID: "",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(nil, assert.AnError)
			},
			input: &messenger.CreateNotificationInput{
				Type:        entity.NotificationTypeSystem,
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Note:        "",
				Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				CreatedBy:   "admin-id",
				PromotionID: "",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "invalid domain validation",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
			},
			input: &messenger.CreateNotificationInput{
				Type:        entity.NotificationTypeSystem,
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Note:        "",
				Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, -1),
				CreatedBy:   "admin-id",
				PromotionID: "",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
				mocks.db.Notification.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &messenger.CreateNotificationInput{
				Type:        entity.NotificationTypeSystem,
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Note:        "",
				Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				CreatedBy:   "admin-id",
				PromotionID: "",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				_, err := service.CreateNotification(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}, withNow(now)),
		)
	}
}

func TestUpdateNotification(t *testing.T) {
	t.Parallel()

	adminIn := &user.GetAdminInput{
		AdminID: "admin-id",
	}
	admin := &uentity.Admin{
		ID: "admin-id",
	}
	now := jst.Now()

	params := &database.UpdateNotificationParams{
		Targets: []entity.NotificationTarget{
			entity.NotificationTargetProducers,
			entity.NotificationTargetCoordinators,
		},
		Title:       "キャベツ祭り開催",
		Body:        "旬のキャベツが大安売り",
		Note:        "",
		PublishedAt: now.AddDate(0, 0, 1),
		UpdatedBy:   "admin-id",
	}
	notification := func() *entity.Notification {
		return &entity.Notification{
			ID:          "notification-id",
			Type:        entity.NotificationTypeSystem,
			Title:       "キャベツ祭り開催",
			Body:        "旬のキャベツを売り出します",
			PublishedAt: now.AddDate(0, 0, 1),
			Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
			CreatedBy:   "admin-id",
			UpdatedBy:   "admin-id",
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.UpdateNotificationInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
				mocks.db.Notification.EXPECT().
					Get(ctx, "notification-id").
					Return(notification(), nil)
				mocks.db.Notification.EXPECT().Update(ctx, "notification-id", params).Return(nil)
				// 非同期関連
				mocks.db.Notification.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			},
			input: &messenger.UpdateNotificationInput{
				NotificationID: "notification-id",
				Title:          "キャベツ祭り開催",
				Body:           "旬のキャベツが大安売り",
				Note:           "",
				Targets: []entity.NotificationTarget{
					entity.NotificationTargetProducers,
					entity.NotificationTargetCoordinators,
				},
				PublishedAt: now.AddDate(0, 0, 1),
				UpdatedBy:   "admin-id",
			},
			expectErr: nil,
		},
		{
			name: "not found admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().
					GetAdmin(gomock.Any(), adminIn).
					Return(nil, exception.ErrNotFound)
			},
			input: &messenger.UpdateNotificationInput{
				NotificationID: "notification-id",
				Title:          "キャベツ祭り開催",
				Body:           "旬のキャベツが大安売り",
				Note:           "",
				Targets: []entity.NotificationTarget{
					entity.NotificationTargetProducers,
					entity.NotificationTargetCoordinators,
				},
				PublishedAt: now.AddDate(0, 0, 1),
				UpdatedBy:   "admin-id",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(nil, assert.AnError)
			},
			input: &messenger.UpdateNotificationInput{
				NotificationID: "notification-id",
				Title:          "キャベツ祭り開催",
				Body:           "旬のキャベツが大安売り",
				Note:           "",
				Targets: []entity.NotificationTarget{
					entity.NotificationTargetProducers,
					entity.NotificationTargetCoordinators,
				},
				PublishedAt: now.AddDate(0, 0, 1),
				UpdatedBy:   "admin-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.UpdateNotificationInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
				mocks.db.Notification.EXPECT().
					Get(ctx, "notification-id").
					Return(nil, assert.AnError)
			},
			input: &messenger.UpdateNotificationInput{
				NotificationID: "notification-id",
				Title:          "キャベツ祭り開催",
				Body:           "旬のキャベツが大安売り",
				Note:           "",
				Targets: []entity.NotificationTarget{
					entity.NotificationTargetProducers,
					entity.NotificationTargetCoordinators,
				},
				PublishedAt: now.AddDate(0, 0, 1),
				UpdatedBy:   "admin-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "already published",
			setup: func(ctx context.Context, mocks *mocks) {
				notification := notification()
				notification.PublishedAt = now.AddDate(0, 0, -1)
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
			},
			input: &messenger.UpdateNotificationInput{
				NotificationID: "notification-id",
				Title:          "キャベツ祭り開催",
				Body:           "旬のキャベツが大安売り",
				Note:           "",
				Targets: []entity.NotificationTarget{
					entity.NotificationTargetProducers,
					entity.NotificationTargetCoordinators,
				},
				PublishedAt: now.AddDate(0, 0, 1),
				UpdatedBy:   "admin-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "invalid domain validation",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
				mocks.db.Notification.EXPECT().
					Get(ctx, "notification-id").
					Return(notification(), nil)
			},
			input: &messenger.UpdateNotificationInput{
				NotificationID: "notification-id",
				Title:          "キャベツ祭り開催",
				Body:           "旬のキャベツが大安売り",
				Note:           "",
				Targets: []entity.NotificationTarget{
					entity.NotificationTargetProducers,
					entity.NotificationTargetCoordinators,
				},
				PublishedAt: now.AddDate(0, 0, -1),
				UpdatedBy:   "admin-id",
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminIn).Return(admin, nil)
				mocks.db.Notification.EXPECT().
					Get(ctx, "notification-id").
					Return(notification(), nil)
				mocks.db.Notification.EXPECT().
					Update(ctx, "notification-id", params).
					Return(assert.AnError)
			},
			input: &messenger.UpdateNotificationInput{
				NotificationID: "notification-id",
				Title:          "キャベツ祭り開催",
				Body:           "旬のキャベツが大安売り",
				Note:           "",
				Targets: []entity.NotificationTarget{
					entity.NotificationTargetProducers,
					entity.NotificationTargetCoordinators,
				},
				PublishedAt: now.AddDate(0, 0, 1),
				UpdatedBy:   "admin-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.UpdateNotification(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}, withNow(now)),
		)
	}
}

func TestDeleteNotification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.DeleteNotificationInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Delete(ctx, "notification-id").Return(nil)
			},
			input: &messenger.DeleteNotificationInput{
				NotificationID: "notification-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.DeleteNotificationInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Delete(ctx, "notification-id").Return(assert.AnError)
			},
			input: &messenger.DeleteNotificationInput{
				NotificationID: "notification-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				err := service.DeleteNotification(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}
