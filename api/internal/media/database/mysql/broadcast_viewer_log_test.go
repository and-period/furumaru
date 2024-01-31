package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcastViewerLog(t *testing.T) {
	assert.NotNil(t, newBroadcastViewerLog(nil))
}

func TestBroadcastViewerLog_Create(t *testing.T) {
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

	broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
	err = db.DB.Create(&broadcast).Error
	require.NoError(t, err)

	type args struct {
		log *entity.BroadcastViewerLog
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
				log: testBroadcastViewerLog("broadcast-id", "session-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				log := testBroadcastViewerLog("broadcast-id", "session-id", "user-id", now())
				err := db.DB.Create(&log).Error
				require.NoError(t, err)
			},
			args: args{
				log: testBroadcastViewerLog("broadcast-id", "session-id", "user-id", now()),
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

			err := delete(ctx, broadcastViewerLogTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &broadcastViewerLog{db: db, now: now}
			err = db.Create(ctx, tt.args.log)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testBroadcastViewerLog(broadcastID, sessionID, userID string, now time.Time) *entity.BroadcastViewerLog {
	return &entity.BroadcastViewerLog{
		BroadcastID: broadcastID,
		SessionID:   sessionID,
		CreatedAt:   now,
		UserID:      userID,
		UserAgent:   "user-agent",
		ClientIP:    "127.0.0.1",
		UpdatedAt:   now,
	}
}
