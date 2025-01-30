package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
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

	internal := make(internalNotifications, 3)
	internal[0] = testNotification("notification-id01", now())
	internal[1] = testNotification("notification-id02", now())
	internal[2] = testNotification("notification-id03", now())
	err = db.DB.Table(notificationTable).Create(&internal).Error
	require.NoError(t, err)
	notifications, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		params *database.ListNotificationsParams
	}
	type want struct {
		notifications entity.Notifications
		err           error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListNotificationsParams{
					Limit:  20,
					Offset: 1,
					Since:  now().Add(-time.Hour),
					Until:  now().AddDate(0, 1, 0),
				},
			},
			want: want{
				notifications: notifications[1:],
				err:           nil,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListNotificationsParams{
					Orders: []*database.ListNotificationsOrder{
						{Key: database.ListNotificationsOrderByPublishedAt, OrderByASC: false},
					},
				},
			},
			want: want{
				notifications: notifications,
				err:           nil,
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
			assert.ErrorIs(t, err, tt.want.err)
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

	internal := make(internalNotifications, 3)
	internal[0] = testNotification("notification-id01", now())
	internal[1] = testNotification("notification-id02", now())
	internal[2] = testNotification("notification-id03", now())
	err = db.DB.Table(notificationTable).Create(&internal).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListNotificationsParams
	}
	type want struct {
		total int64
		err   error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListNotificationsParams{
					Since: now(),
					Until: now().AddDate(0, 1, 0),
				},
			},
			want: want{
				total: 3,
				err:   nil,
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
			assert.ErrorIs(t, err, tt.want.err)
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

	internal := testNotification("notification-id", now())
	err = db.DB.Table(notificationTable).Create(&internal).Error
	require.NoError(t, err)
	n, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		notificationID string
	}
	type want struct {
		notification *entity.Notification
		err          error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				notificationID: "notification-id",
			},
			want: want{
				notification: n,
				err:          nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				notificationID: "other-id",
			},
			want: want{
				notification: nil,
				err:          database.ErrNotFound,
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
			assert.ErrorIs(t, err, tt.want.err)
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

	internal := testNotification("notification-id", now())
	n, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		notification *entity.Notification
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				notification: n,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				n := testNotification("notification-id", now())
				err = db.DB.Table(notificationTable).Create(&n).Error
				require.NoError(t, err)
			},
			args: args{
				notification: n,
			},
			want: want{
				err: database.ErrAlreadyExists,
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
			assert.ErrorIs(t, err, tt.want.err)
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
		params         *database.UpdateNotificationParams
	}
	type want = struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				notification := testNotification("notification-id", now().AddDate(0, 0, 1))
				err = db.DB.Table(notificationTable).Create(&notification).Error
				require.NoError(t, err)
			},
			args: args{
				notificationID: "notification-id",
				params: &database.UpdateNotificationParams{
					Title: "キャベツ祭り開催",
					Body:  "旬のキャベツが大安売り",
					Targets: []entity.NotificationTarget{
						entity.NotificationTargetProducers,
						entity.NotificationTargetCoordinators,
					},
					PublishedAt: now().AddDate(0, 0, 1),
					UpdatedBy:   "admin-id",
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				notificationID: "notification-id",
				params:         &database.UpdateNotificationParams{},
			},
			want: want{
				err: database.ErrNotFound,
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
			assert.ErrorIs(t, err, tt.want.err)
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
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				notification := testNotification("notification-id", now())
				err = db.DB.Table(notificationTable).Create(&notification).Error
				require.NoError(t, err)
			},
			args: args{
				notificationID: "notification-id",
			},
			want: want{
				err: nil,
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
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testNotification(id string, now time.Time) *internalNotification {
	notification := &entity.Notification{
		ID:          id,
		Status:      entity.NotificationStatusWaiting,
		Title:       "お知らせタイトル",
		Body:        "お知らせの内容です。",
		Note:        "備考です",
		Targets:     []entity.NotificationTarget{entity.NotificationTargetUsers},
		CreatedBy:   "coordinator-id",
		UpdatedBy:   "coordinator-id",
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: now.AddDate(0, 0, 1),
	}
	internal, _ := newInternalNotification(notification)
	return internal
}
