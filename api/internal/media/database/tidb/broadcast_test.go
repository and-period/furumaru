package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBroadcast(t *testing.T) {
	assert.NotNil(t, NewBroadcast(nil))
}

func TestBroadcast_List(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	internal := make(internalBroadcasts, 2)
	internal[0] = testBroadcast("broadcast-id01", "schedule-id01", "coordinator-id", now().AddDate(0, 1, 0))
	internal[1] = testBroadcast("broadcast-id02", "schedule-id02", "coordinator-id", now())
	err = db.DB.Table(broadcastTable).Create(&internal).Error
	require.NoError(t, err)
	broadcasts, err := internal.entities()
	require.NoError(t, err)

	type args struct {
		params *database.ListBroadcastsParams
	}
	type want struct {
		broadcasts entity.Broadcasts
		err        error
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
				params: &database.ListBroadcastsParams{
					Limit:  20,
					Offset: 1,
				},
			},
			want: want{
				broadcasts: broadcasts[1:],
				err:        nil,
			},
		},
		{
			name:  "success only archived",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListBroadcastsParams{
					OnlyArchived: true,
					Orders: []*database.ListBroadcastsOrder{
						{
							Key:        database.ListBroadcastsOrderByUpdatedAt,
							OrderByASC: true,
						},
					},
				},
			},
			want: want{
				broadcasts: entity.Broadcasts{},
				err:        nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &broadcast{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.broadcasts, actual)
		})
	}
}

func TestBroadcast_Count(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	internal := make(internalBroadcasts, 2)
	internal[0] = testBroadcast("broadcast-id01", "schedule-id01", "coordinator-id", now().AddDate(0, 1, 0))
	internal[1] = testBroadcast("broadcast-id02", "schedule-id02", "coordinator-id", now())
	err = db.DB.Table(broadcastTable).Create(&internal).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListBroadcastsParams
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
				params: &database.ListBroadcastsParams{
					Limit:  20,
					Offset: 1,
				},
			},
			want: want{
				total: 2,
				err:   nil,
			},
		},
		{
			name:  "success only archived",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListBroadcastsParams{
					OnlyArchived: true,
					Orders: []*database.ListBroadcastsOrder{
						{
							Key:        database.ListBroadcastsOrderByUpdatedAt,
							OrderByASC: true,
						},
					},
				},
			},
			want: want{
				total: 0,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &broadcast{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestBroadcast_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	internal := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
	err = db.DB.Table(broadcastTable).Create(&internal).Error
	require.NoError(t, err)
	b, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		broadcastID string
	}
	type want struct {
		broadcast *entity.Broadcast
		err       error
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
				broadcastID: "broadcast-id",
			},
			want: want{
				broadcast: b,
				err:       nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				broadcastID: "",
			},
			want: want{
				broadcast: nil,
				err:       database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &broadcast{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.broadcastID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.broadcast, actual)
		})
	}
}

func TestBroadcast_GetByScheduleID(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	internal := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
	err = db.DB.Table(broadcastTable).Create(&internal).Error
	require.NoError(t, err)
	b, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		scheduleID string
	}
	type want struct {
		broadcast *entity.Broadcast
		err       error
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
				scheduleID: "schedule-id",
			},
			want: want{
				broadcast: b,
				err:       nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				scheduleID: "",
			},
			want: want{
				broadcast: nil,
				err:       database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &broadcast{db: db, now: now}
			actual, err := db.GetByScheduleID(ctx, tt.args.scheduleID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.broadcast, actual)
		})
	}
}

func TestBroadcast_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	internal := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
	b, err := internal.entity()
	require.NoError(t, err)

	type args struct {
		broadcast *entity.Broadcast
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
				broadcast: b,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
				err := db.DB.Table(broadcastTable).Create(&broadcast).Error
				require.NoError(t, err)
			},
			args: args{
				broadcast: b,
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, broadcastTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &broadcast{db: db, now: now}
			err = db.Create(ctx, tt.args.broadcast)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestBroadcast_Update(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
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
		broadcastID string
		params      *database.UpdateBroadcastParams
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
			name: "success active",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
				err = db.DB.Table(broadcastTable).Create(&broadcast).Error
				require.NoError(t, err)
			},
			args: args{
				broadcastID: "broadcast-id",
				params: &database.UpdateBroadcastParams{
					Status: entity.BroadcastStatusActive,
					InitializeBroadcastParams: &database.InitializeBroadcastParams{
						InputURL:                  "rtmp://127.0.0.1:1935/live/a",
						OutputURL:                 "http://example.com/index.m3u8",
						CloudFrontDistributionArn: "aws/arn",
						MediaLiveChannelArn:       "aws/arn",
						MediaLiveChannelID:        "channel-id",
						MediaLiveRTMPInputArn:     "aws/arn",
						MediaLiveRTMPInputName:    "rtmp-input-name",
						MediaLiveMP4InputArn:      "aws/arn",
						MediaLiveMP4InputName:     "mp4-input-name",
						MediaStoreContainerArn:    "aws/arn",
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success upload archive",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
				err = db.DB.Table(broadcastTable).Create(&broadcast).Error
				require.NoError(t, err)
			},
			args: args{
				broadcastID: "broadcast-id",
				params: &database.UpdateBroadcastParams{
					Status: entity.BroadcastStatusActive,
					UploadBroadcastArchiveParams: &database.UploadBroadcastArchiveParams{
						ArchiveURL:   "http://example.com/master.mp4",
						ArchiveFixed: true,
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success update archive",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
				err = db.DB.Table(broadcastTable).Create(&broadcast).Error
				require.NoError(t, err)
			},
			args: args{
				broadcastID: "broadcast-id",
				params: &database.UpdateBroadcastParams{
					Status: entity.BroadcastStatusActive,
					UpdateBroadcastArchiveParams: &database.UpdateBroadcastArchiveParams{
						ArchiveMetadata: &entity.BroadcastArchiveMetadata{
							Subtitles: map[string]string{
								"jpn": "http://example.com/translate-jpn.vtt",
								"eng": "http://example.com/translate-eng.vtt",
							},
						},
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success youtube",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
				err = db.DB.Table(broadcastTable).Create(&broadcast).Error
				require.NoError(t, err)
			},
			args: args{
				broadcastID: "broadcast-id",
				params: &database.UpdateBroadcastParams{
					Status: entity.BroadcastStatusActive,
					UpsertYoutubeBroadcastParams: &database.UpsertYoutubeBroadcastParams{
						YoutubeAccount:     "test@example.com",
						YoutubeBroadcastID: "broadcast-id",
						YoutubeStreamID:    "stream-id",
						YoutubeStreamURL:   "rtmp://a.rtmp.youtube.com/live2",
						YoutubeStreamKey:   "stream-key",
						YoutubeBackupURL:   "rtmp://b.rtmp.youtube.com/live2?backup=1",
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success disable",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
				err = db.DB.Table(broadcastTable).Create(&broadcast).Error
				require.NoError(t, err)
			},
			args: args{
				broadcastID: "broadcast-id",
				params: &database.UpdateBroadcastParams{
					Status: entity.BroadcastStatusDisabled,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success archive",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
				err = db.DB.Table(broadcastTable).Create(&broadcast).Error
				require.NoError(t, err)
			},
			args: args{
				broadcastID: "broadcast-id",
				params: &database.UpdateBroadcastParams{
					Status: entity.BroadcastStatusDisabled,
					UploadBroadcastArchiveParams: &database.UploadBroadcastArchiveParams{
						ArchiveURL: "http://example.com/master.mp4",
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success other",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
				err = db.DB.Table(broadcastTable).Create(&broadcast).Error
				require.NoError(t, err)
			},
			args: args{
				broadcastID: "broadcast-id",
				params: &database.UpdateBroadcastParams{
					Status: entity.BroadcastStatusWaiting,
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, broadcastTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &broadcast{db: db, now: now}
			err = db.Update(ctx, tt.args.broadcastID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testBroadcast(broadcastID, scheduleID, coordinatorID string, now time.Time) *internalBroadcast {
	broadcast := &entity.Broadcast{
		ID:              broadcastID,
		ScheduleID:      scheduleID,
		CoordinatorID:   coordinatorID,
		Type:            entity.BroadcastTypeNormal,
		Status:          entity.BroadcastStatusIdle,
		InputURL:        "rtmp://127.0.0.1/1935/app/instance",
		OutputURL:       "http://example.com/index.m3u8",
		ArchiveURL:      "",
		ArchiveMetadata: &entity.BroadcastArchiveMetadata{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	internal, _ := newInternalBroadcast(broadcast)
	return internal
}
