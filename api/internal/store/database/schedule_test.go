package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchedule(t *testing.T) {
	assert.NotNil(t, NewSchedule(nil))
}

func TestSchedule_List(t *testing.T) {
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

	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	schedules := make(entity.Schedules, 2)
	schedules[0] = testSchedule("schedule-id01", "coordinator-id", "shipping-id", now())
	schedules[1] = testSchedule("schedule-id02", "coordinator-id", "shipping-id", now())
	err = db.DB.Create(&schedules).Error
	require.NoError(t, err)

	type args struct {
		params *ListSchedulesParams
	}
	type want struct {
		schedules entity.Schedules
		hasErr    bool
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
				params: &ListSchedulesParams{
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				schedules: schedules[1:],
				hasErr:    false,
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

			db := &schedule{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.schedules, actual)
		})
	}
}

func TestSchedule_Count(t *testing.T) {
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

	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	schedules := make(entity.Schedules, 2)
	schedules[0] = testSchedule("schedule-id01", "coordinator-id", "shipping-id", now())
	schedules[1] = testSchedule("schedule-id02", "coordinator-id", "shipping-id", now())
	err = db.DB.Create(&schedules).Error
	require.NoError(t, err)

	type args struct {
		params *ListSchedulesParams
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
				params: &ListSchedulesParams{
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				total:  2,
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

			db := &schedule{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestSchedule_Get(t *testing.T) {
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

	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	s := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
	err = db.DB.Create(&s).Error
	require.NoError(t, err)

	type args struct {
		scheduleID string
	}
	type want struct {
		schedule *entity.Schedule
		hasErr   bool
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
				scheduleID: "schedule-id",
			},
			want: want{
				schedule: s,
				hasErr:   false,
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

			db := &schedule{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.scheduleID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.schedule, actual)
		})
	}
}

func TestSchedule_Create(t *testing.T) {
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

	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	s := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
	b := testBroadcast("broadcast-id", "schedule-id", now())

	type args struct {
		schedule  *entity.Schedule
		broadcast *entity.Broadcast
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
				schedule:  s,
				broadcast: b,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				schedule: s,
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

			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.Create(ctx, tt.args.schedule, tt.args.broadcast)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestSchedule_Update(t *testing.T) {
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

	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)

	type args struct {
		scheduleID string
		params     *UpdateScheduleParams
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
				schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				scheduleID: "schedule-id",
				params: &UpdateScheduleParams{
					ShippingID:      "shipping-id",
					Title:           "開催スケジュール",
					Description:     "開催スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					Public:          true,
					StartAt:         now().AddDate(0, -1, 0),
					EndAt:           now().AddDate(0, 1, 0),
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				scheduleID: "schedule-id",
				params:     &UpdateScheduleParams{},
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

			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.Update(ctx, tt.args.scheduleID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestSchedule_UpdateThumbnails(t *testing.T) {
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

	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)

	type args struct {
		scheduleID string
		thumbnails common.Images
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
				schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				scheduleID: "schedule-id",
				thumbnails: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/thumbnail_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/thumbnail_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/thumbnail_900.png",
					},
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				scheduleID: "schedule-id",
				thumbnails: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/thumbnail_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/thumbnail_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/thumbnail_900.png",
					},
				},
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to empty thumbnail url",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
				schedule.ThumbnailURL = ""
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				scheduleID: "schedule-id",
				thumbnails: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "https://and-period.jp/thumbnail_240.png",
					},
					{
						Size: common.ImageSizeMedium,
						URL:  "https://and-period.jp/thumbnail_675.png",
					},
					{
						Size: common.ImageSizeLarge,
						URL:  "https://and-period.jp/thumbnail_900.png",
					},
				},
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

			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.UpdateThumbnails(ctx, tt.args.scheduleID, tt.args.thumbnails)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testSchedule(id, coordinatorID, shippingID string, now time.Time) *entity.Schedule {
	schedule := &entity.Schedule{
		ID:              id,
		CoordinatorID:   coordinatorID,
		ShippingID:      shippingID,
		Status:          entity.ScheduleStatusLive,
		Title:           "旬の夏野菜配信",
		Description:     "旬の夏野菜特集",
		ThumbnailURL:    "https://and-period.jp/thumbnail.png",
		ImageURL:        "https://and-period.jp/image.png",
		OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
		Public:          true,
		Approved:        true,
		ApprovedAdminID: "admin-id",
		StartAt:         now.AddDate(0, -1, 0),
		EndAt:           now.AddDate(0, 1, 0),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	schedule.Fill(now)
	return schedule
}
