package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
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

func TestBroadcastViewerLog_Aggregate(t *testing.T) {
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

	logs := make(entity.BroadcastViewerLogs, 4)
	logs[0] = testBroadcastViewerLog("broadcast-id", "session-id01", "user-id01", now())
	logs[1] = testBroadcastViewerLog("broadcast-id", "session-id01", "user-id01", now().Add(1*time.Minute))
	logs[2] = testBroadcastViewerLog("broadcast-id", "session-id02", "user-id02", now())
	logs[3] = testBroadcastViewerLog("broadcast-id", "session-id02", "user-id02", now().Add(1*time.Minute))
	logs[3].UserAgent = entity.ExcludeUserAgentLogs[0]
	for _, log := range logs {
		err = db.DB.Create(&log).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.AggregateBroadcastViewerLogsParams
	}
	type want struct {
		logs entity.AggregatedBroadcastViewerLogs
		err  error
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
				params: &database.AggregateBroadcastViewerLogsParams{
					BroadcastID:  "broadcast-id",
					Interval:     entity.AggregateBroadcastViewerLogIntervalMinute,
					CreatedAtGte: now().Add(-time.Minute),
					CreatedAtLt:  now().Add(time.Hour),
				},
			},
			want: want{
				logs: entity.AggregatedBroadcastViewerLogs{
					{
						BroadcastID: "broadcast-id",
						ReportedAt:  jst.Date(now().Year(), now().Month(), now().Day(), now().Hour(), now().Minute(), 0, 0),
						Total:       2,
					},
					{
						BroadcastID: "broadcast-id",
						ReportedAt:  jst.Date(now().Year(), now().Month(), now().Day(), now().Hour(), now().Minute()+1, 0, 0),
						Total:       1,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &broadcastViewerLog{db: db, now: now}
			logs, err := db.Aggregate(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.logs, logs)
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
