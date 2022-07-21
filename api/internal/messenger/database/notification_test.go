package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotification(t *testing.T) {
	assert.NotNil(t, NewNotification(nil))
}

func TestNotification_List(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, notificationTable)
	notifications := make(entity.Notifications, 3)
	notifications[0] = testNotification("notification-id01", false, now())
	notifications[1] = testNotification("notification-id02", true, now())
	notifications[2] = testNotification("notification-id03", true, now())
	err = m.db.DB.Create(&notifications).Error
	require.NoError(t, err)

	type args struct {
		params *ListNotificationsParams
	}
	type want struct {
		notifications entity.Notifications
		hasErr        bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				params: &ListNotificationsParams{
					Limit:         20,
					Offset:        1,
					Since:         now().Add(-time.Hour),
					Until:         now(),
					OnlyPublished: true,
				},
			},
			want: want{
				notifications: notifications[2:],
				hasErr:        false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &notification{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreNotificationsField(actual, now())
			assert.ElementsMatch(t, tt.want.notifications, actual)
		})
	}
}

func TestNotification_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, notificationTable)
	n := testNotification("notification-id", false, now())
	err = m.db.DB.Create(&n).Error
	require.NoError(t, err)

	type args struct {
		notificationID string
	}
	type want struct {
		notification *entity.Notification
		hasErr       bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				notificationID: "notification-id",
			},
			want: want{
				notification: n,
				hasErr:       false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				notificationID: "other-id",
			},
			want: want{
				notification: nil,
				hasErr:       true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &notification{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.notificationID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreNotificationField(actual, now())
			assert.Equal(t, tt.want.notification, actual)
		})
	}
}

func TestNotification_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 6, 26, 19, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		notification *entity.Notification
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				notification: testNotification("notification-id", true, now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				n := testNotification("notification-id", true, now())
				err = m.db.DB.Create(&n).Error
				require.NoError(t, err)
			},
			args: args{
				notification: testNotification("notification-id", true, now()),
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := m.dbDelete(ctx, notificationTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &notification{db: m.db, now: now}
			err = db.Create(ctx, tt.args.notification)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testNotification(id string, public bool, now time.Time) *entity.Notification {
	n := &entity.Notification{
		ID:          id,
		Title:       "お知らせタイトル",
		Body:        "お知らせの内容です。",
		Targets:     []entity.TargetType{entity.PostTargetProducers},
		Public:      public,
		CreatorName: "&. スタッフ",
		CreatedBy:   "coordinator-id",
		UpdatedBy:   "coordinator-id",
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: now,
	}
	_ = n.FillJSON()
	return n
}

func fillIgnoreNotificationField(n *entity.Notification, now time.Time) {
	if n == nil {
		return
	}
	_ = n.FillJSON()
	n.PublishedAt = now
	n.CreatedAt = now
	n.UpdatedAt = now
}

func fillIgnoreNotificationsField(ns entity.Notifications, now time.Time) {
	for i := range ns {
		fillIgnoreNotificationField(ns[i], now)
	}
}
