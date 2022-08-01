package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotification(t *testing.T) {
	t.Parallel()

	adminID := &user.GetAdminInput{
		AdminID: "admin-id",
	}
	admin := &uentity.Admin{
		ID:        "admin-id",
		Firstname: "あんど",
		Lastname:  "ぴりおど",
	}
	now := jst.Now()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.CreateNotificationInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminID).Return(admin, nil)
				mocks.db.Notification.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, notification *entity.Notification) error {
						expect := &entity.Notification{
							ID:          notification.ID, // ignore
							CreatedBy:   "admin-id",
							CreatorName: "ぴりおど あんど",
							UpdatedBy:   "admin-id",
							Title:       "キャベツ祭り開催",
							Body:        "旬のキャベツを売り出します",
							Targets:     []entity.TargetType{1, 2},
							Public:      true,
							PublishedAt: now,
						}
						assert.Equal(t, expect, notification)
						return nil
					})
			},
			input: &messenger.CreateNotificationInput{
				CreatedBy:   "admin-id",
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Targets:     []entity.TargetType{1, 2},
				Public:      true,
				PublishedAt: now,
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
			name:  "invalid targets format",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &messenger.CreateNotificationInput{
				CreatedBy:   "admin-id",
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Targets:     []entity.TargetType{4},
				Public:      true,
				PublishedAt: now,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get admin",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminID).Return(nil, exception.ErrNotFound)
			},
			input: &messenger.CreateNotificationInput{
				CreatedBy:   "admin-id",
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Targets:     []entity.TargetType{1, 2},
				Public:      true,
				PublishedAt: now,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), adminID).Return(admin, nil)
				mocks.db.Notification.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &messenger.CreateNotificationInput{
				CreatedBy:   "admin-id",
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを売り出します",
				Targets:     []entity.TargetType{1, 2},
				Public:      true,
				PublishedAt: now,
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateNotification(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
