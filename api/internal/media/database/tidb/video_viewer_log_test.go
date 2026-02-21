package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoViewerLog(t *testing.T) {
	assert.NotNil(t, NewVideoViewerLog(nil))
}

func TestVideoViewerLog_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(t.Context())
	require.NoError(t, err)

	video := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	err = db.DB.Create(&video).Error
	require.NoError(t, err)

	type args struct {
		log *entity.VideoViewerLog
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
				log: testVideoViewerLog("video-id", "session-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				log := testVideoViewerLog("video-id", "session-id", "user-id", now())
				err := db.DB.Create(&log).Error
				require.NoError(t, err)
			},
			args: args{
				log: testVideoViewerLog("video-id", "session-id", "user-id", now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, videoViewerLogTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &videoViewerLog{db: db, now: now}
			err = db.Create(ctx, tt.args.log)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestVideoViewerLog_GetTotal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(t.Context())
	require.NoError(t, err)

	video := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	err = db.DB.Create(&video).Error
	require.NoError(t, err)

	logs := make(entity.VideoViewerLogs, 4)
	logs[0] = testVideoViewerLog("video-id", "session-id01", "user-id01", now())
	logs[1] = testVideoViewerLog("video-id", "session-id01", "user-id01", now().Add(1*time.Minute))
	logs[2] = testVideoViewerLog("video-id", "session-id02", "user-id02", now())
	logs[3] = testVideoViewerLog("video-id", "session-id02", "user-id02", now().Add(1*time.Minute))
	logs[3].UserAgent = entity.ExcludeUserAgentLogs[0]
	for _, log := range logs {
		err = db.DB.Create(&log).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.GetVideoTotalViewersParams
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
				params: &database.GetVideoTotalViewersParams{
					VideoID:      "video-id",
					CreatedAtGte: now().Add(-time.Minute),
					CreatedAtLt:  now().Add(time.Hour),
				},
			},
			want: want{
				total: 2,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			tt.setup(ctx, t, db)

			db := &videoViewerLog{db: db, now: now}
			total, err := db.GetTotal(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, total)
		})
	}
}

func TestVideoViewerLog_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(t.Context())
	require.NoError(t, err)

	video := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	err = db.DB.Create(&video).Error
	require.NoError(t, err)

	logs := make(entity.VideoViewerLogs, 4)
	logs[0] = testVideoViewerLog("video-id", "session-id01", "user-id01", now())
	logs[1] = testVideoViewerLog("video-id", "session-id01", "user-id01", now().Add(1*time.Minute))
	logs[2] = testVideoViewerLog("video-id", "session-id02", "user-id02", now())
	logs[3] = testVideoViewerLog("video-id", "session-id02", "user-id02", now().Add(1*time.Minute))
	logs[3].UserAgent = entity.ExcludeUserAgentLogs[0]
	for _, log := range logs {
		err = db.DB.Create(&log).Error
		require.NoError(t, err)
	}

	type args struct {
		params *database.AggregateVideoViewerLogsParams
	}
	type want struct {
		logs entity.AggregatedVideoViewerLogs
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
				params: &database.AggregateVideoViewerLogsParams{
					VideoID:      "video-id",
					Interval:     entity.AggregateVideoViewerLogIntervalMinute,
					CreatedAtGte: now().Add(-time.Minute),
					CreatedAtLt:  now().Add(time.Hour),
				},
			},
			want: want{
				logs: entity.AggregatedVideoViewerLogs{
					{
						VideoID:    "video-id",
						ReportedAt: jst.Date(now().Year(), now().Month(), now().Day(), now().Hour(), now().Minute(), 0, 0),
						Total:      2,
					},
					{
						VideoID:    "video-id",
						ReportedAt: jst.Date(now().Year(), now().Month(), now().Day(), now().Hour(), now().Minute()+1, 0, 0),
						Total:      1,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			tt.setup(ctx, t, db)

			db := &videoViewerLog{db: db, now: now}
			logs, err := db.Aggregate(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.logs, logs)
		})
	}
}

func testVideoViewerLog(videoID, sessionID, userID string, now time.Time) *entity.VideoViewerLog {
	return &entity.VideoViewerLog{
		VideoID:   videoID,
		SessionID: sessionID,
		CreatedAt: now,
		UserID:    userID,
		UserAgent: "user-agent",
		ClientIP:  "127.0.0.1",
		UpdatedAt: now,
	}
}
