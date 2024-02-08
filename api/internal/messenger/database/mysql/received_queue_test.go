package mysql

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

func TestReceivedQueue(t *testing.T) {
	assert.NotNil(t, newReceivedQueue(nil))
}

func TestReceivedQueue_Get(t *testing.T) {
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

	q := testReceivedQueue("queue-id", entity.NotifyTypeEmail, now())
	err = db.DB.Create(&q).Error
	require.NoError(t, err)

	type args struct {
		queueID    string
		notifyType entity.NotifyType
	}
	type want struct {
		queue *entity.ReceivedQueue
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
				queueID:    "queue-id",
				notifyType: entity.NotifyTypeEmail,
			},
			want: want{
				queue: q,
				err:   nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				queueID:    "other-id",
				notifyType: entity.NotifyTypeEmail,
			},
			want: want{
				queue: nil,
				err:   database.ErrNotFound,
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

			db := &receivedQueue{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.queueID, tt.args.notifyType)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.queue, actual)
		})
	}
}

func TestReceivedQueue_Create(t *testing.T) {
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

	q := testReceivedQueue("queue-id", entity.NotifyTypeEmail, now())

	type args struct {
		queue *entity.ReceivedQueue
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
				queue: q,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err = db.DB.Create(&q).Error
				require.NoError(t, err)
			},
			args: args{
				queue: q,
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

			err := delete(ctx, receivedQueueTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &receivedQueue{db: db, now: now}
			err = db.MultiCreate(ctx, tt.args.queue)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestReceivedQueue_UpdateDone(t *testing.T) {
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

	q := testReceivedQueue("queue-id", entity.NotifyTypeEmail, now())

	type args struct {
		queueID    string
		notifyType entity.NotifyType
		done       bool
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
				err = db.DB.Create(&q).Error
				require.NoError(t, err)
			},
			args: args{
				queueID: "queue-id",
				done:    true,
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

			err := delete(ctx, receivedQueueTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &receivedQueue{db: db, now: now}
			err = db.UpdateDone(ctx, tt.args.queueID, tt.args.notifyType, tt.args.done)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testReceivedQueue(id string, typ entity.NotifyType, now time.Time) *entity.ReceivedQueue {
	q := &entity.ReceivedQueue{
		ID:         id,
		NotifyType: typ,
		EventType:  entity.EventTypeUnknown,
		UserType:   entity.UserTypeUser,
		UserIDs:    []string{"user-id"},
		Done:       false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_ = q.FillJSON()
	return q
}
