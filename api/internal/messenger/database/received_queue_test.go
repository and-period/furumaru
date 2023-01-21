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

func TestReceivedQueue(t *testing.T) {
	assert.NotNil(t, NewReceivedQueue(nil))
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

	q := testReceivedQueue("queue-id", now())
	err = db.DB.Create(&q).Error
	require.NoError(t, err)

	type args struct {
		queueID string
	}
	type want struct {
		queue  *entity.ReceivedQueue
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
				queueID: "queue-id",
			},
			want: want{
				queue:  q,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				queueID: "other-id",
			},
			want: want{
				queue:  nil,
				hasErr: true,
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
			actual, err := db.Get(ctx, tt.args.queueID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
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

	q := testReceivedQueue("queue-id", now())

	type args struct {
		queue *entity.ReceivedQueue
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
				queue: q,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				err = db.DB.Create(&q).Error
				require.NoError(t, err)
			},
			args: args{
				queue: q,
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

			err := delete(ctx, receivedQueueTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &receivedQueue{db: db, now: now}
			err = db.Create(ctx, tt.args.queue)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestReceivedQueue_Update(t *testing.T) {
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

	q := testReceivedQueue("queue-id", now())

	type args struct {
		queueID string
		done    bool
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
				err = db.DB.Create(&q).Error
				require.NoError(t, err)
			},
			args: args{
				queueID: "queue-id",
				done:    true,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				queueID: "queue-id",
				done:    true,
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

			err := delete(ctx, receivedQueueTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &receivedQueue{db: db, now: now}
			err = db.UpdateDone(ctx, tt.args.queueID, tt.args.done)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testReceivedQueue(id string, now time.Time) *entity.ReceivedQueue {
	q := &entity.ReceivedQueue{
		ID:        id,
		EventType: entity.EventTypeUnknown,
		UserType:  entity.UserTypeUser,
		UserIDs:   []string{"user-id"},
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = q.FillJSON()
	return q
}
