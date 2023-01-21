package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
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

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	notifications := make(entity.Notifications, 3)
	notifications[0] = testNotification("notification-id01", false, now())
	notifications[1] = testNotification("notification-id02", true, now())
	notifications[2] = testNotification("notification-id03", true, now())
	err = db.DB.Create(&notifications).Error
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
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
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
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListNotificationsParams{
					Orders: []*ListNotificationsOrder{
						{Key: entity.NotificationOrderByPublishedAt, OrderByASC: false},
					},
				},
			},
			want: want{
				notifications: notifications,
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

			tt.setup(ctx, t, db)

			db := &notification{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.notifications, actual)
		})
	}
}

func TestNotificaiton_Count(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	notifications := make(entity.Notifications, 3)
	notifications[0] = testNotification("notification-id01", true, now())
	notifications[1] = testNotification("notification-id02", true, now())
	notifications[2] = testNotification("notification-id03", true, now())
	err = db.DB.Create(&notifications).Error
	require.NoError(t, err)

	type args struct {
		params *ListNotificationsParams
	}
	type want struct {
		total  int64
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListNotificationsParams{
					Since: now(),
					Until: now(),
				},
			},
			want: want{
				total:  3,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &notification{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestNotification_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	n := testNotification("notification-id", false, now())
	err = db.DB.Create(&n).Error
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
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
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
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
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

			tt.setup(ctx, t, db)

			db := &notification{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.notificationID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.notification, actual)
		})
	}
}

func TestNotification_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		notification *entity.Notification
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				notification: testNotification("notification-id", true, now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				n := testNotification("notification-id", true, now())
				err = db.DB.Create(&n).Error
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

			err := delete(ctx, notificationTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &notification{db: db, now: now}
			err = db.Create(ctx, tt.args.notification)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestNotification_Update(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args = struct {
		notificationID string
		params         *UpdateNotificationParams
	}
	type want = struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				notification := testNotification("notification-id", true, now())
				err = db.DB.Create(&notification).Error
				require.NoError(t, err)
			},
			args: args{
				notificationID: "notification-id",
				params: &UpdateNotificationParams{
					Title: "キャベツ祭り開催",
					Body:  "旬のキャベツが大安売り",
					Targets: []entity.TargetType{
						entity.PostTargetProducers,
						entity.PostTargetCoordinators,
					},
					PublishedAt: now(),
					UpdatedBy:   "admin-id",
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				notificationID: "notification-id",
				params:         &UpdateNotificationParams{},
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

			err := delete(ctx, notificationTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &notification{db: db, now: now}
			err = db.Update(ctx, tt.args.notificationID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestNotificaiton_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	type args struct {
		notificationID string
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				notification := testNotification("notification-id", true, now())
				err = db.DB.Create(&notification).Error
				require.NoError(t, err)
			},
			args: args{
				notificationID: "notification-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				notificationID: "notification-id",
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

			err := delete(ctx, notificationTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &notification{db: db, now: now}
			err = db.Delete(ctx, tt.args.notificationID)
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
