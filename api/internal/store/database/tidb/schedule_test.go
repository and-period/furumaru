package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchedule(t *testing.T) {
	assert.NotNil(t, NewSchedule(nil))
}

func TestSchedule_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shops := make(internalShops, 2)
	shops[0] = testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	shops[1] = testShop("shop-id02", "coordinator-id02", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shops).Error
	require.NoError(t, err)

	schedules := make(entity.Schedules, 2)
	schedules[0] = testSchedule("schedule-id01", "shop-id01", "coordinator-id", now())
	schedules[1] = testSchedule(
		"schedule-id02",
		"shop-id02",
		"coordinator-id",
		now().Add(time.Hour),
	)
	err = db.DB.Create(&schedules).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListSchedulesParams
	}
	type want struct {
		schedules entity.Schedules
		hasErr    bool
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
				params: &database.ListSchedulesParams{
					Limit:  2,
					Offset: 0,
				},
			},
			want: want{
				schedules: schedules,
				hasErr:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.schedules, actual)
		})
	}
}

func TestSchedule_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shops := make(internalShops, 2)
	shops[0] = testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	shops[1] = testShop("shop-id02", "coordinator-id02", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shops).Error
	require.NoError(t, err)

	schedules := make(entity.Schedules, 2)
	schedules[0] = testSchedule("schedule-id01", "shop-id01", "coordinator-id01", now())
	schedules[1] = testSchedule("schedule-id02", "shop-id02", "coordinator-id02", now())
	err = db.DB.Create(&schedules).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListSchedulesParams
	}
	type want struct {
		total  int64
		hasErr bool
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
				params: &database.ListSchedulesParams{
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestSchedule_MultiGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shops := make(internalShops, 2)
	shops[0] = testShop("shop-id01", "coordinator-id01", []string{}, []string{}, now())
	shops[1] = testShop("shop-id02", "coordinator-id02", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shops).Error
	require.NoError(t, err)

	schedules := make(entity.Schedules, 2)
	schedules[0] = testSchedule("schedule-id01", "shop-id01", "coordinator-id01", now())
	schedules[1] = testSchedule("schedule-id02", "shop-id02", "coordinator-id02", now())
	err = db.DB.Create(&schedules).Error
	require.NoError(t, err)

	type args struct {
		scheduleIDs []string
	}
	type want struct {
		schedules entity.Schedules
		hasErr    bool
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
				scheduleIDs: []string{"schedule-id01", "schedule-id02"},
			},
			want: want{
				schedules: schedules,
				hasErr:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.scheduleIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.schedules, actual)
		})
	}
}

func TestSchedule_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	s := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
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
				schedule: s,
				hasErr:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.scheduleID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.schedule, actual)
		})
	}
}

func TestSchedule_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	s := testSchedule("schedule-id", "shop-id", "coordinator-id", now())

	type args struct {
		schedule *entity.Schedule
	}
	type want struct {
		hasErr bool
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
				schedule: s,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
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
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.Create(ctx, tt.args.schedule)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestSchedule_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	type args struct {
		scheduleID string
		params     *database.UpdateScheduleParams
	}
	type want struct {
		hasErr bool
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
				schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
				schedule.StartAt = now().AddDate(0, 1, 0)
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				scheduleID: "schedule-id",
				params: &database.UpdateScheduleParams{
					Title:           "開催スケジュール",
					Description:     "開催スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					StartAt:         now().AddDate(0, 1, 0),
					EndAt:           now().AddDate(0, 2, 0),
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "failed to not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				scheduleID: "schedule-id",
				params: &database.UpdateScheduleParams{
					Title:           "開催スケジュール",
					Description:     "開催スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					StartAt:         now().AddDate(0, 1, 0),
					EndAt:           now().AddDate(0, 2, 0),
				},
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "failed to update",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				scheduleID: "schedule-id",
				params: &database.UpdateScheduleParams{
					Title:           "開催スケジュール",
					Description:     "開催スケジュールの詳細です。",
					ThumbnailURL:    "https://and-period.jp/thumbnail.png",
					ImageURL:        "https://and-period.jp/image.png",
					OpeningVideoURL: "https://and-period.jp/opening-video.mp4",
					StartAt:         now().AddDate(0, 1, 0),
					EndAt:           now().AddDate(0, 2, 0),
				},
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.Update(ctx, tt.args.scheduleID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestSchedule_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	type args struct {
		scheduleID string
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
				schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				scheduleID: "schedule-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.Delete(ctx, tt.args.scheduleID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestSchedule_Approve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	type args struct {
		scheduleID string
		params     *database.ApproveScheduleParams
	}
	type want struct {
		hasErr bool
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
				schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
				schedule.StartAt = now().AddDate(0, 1, 0)
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				scheduleID: "schedule-id",
				params: &database.ApproveScheduleParams{
					Approved:        true,
					ApprovedAdminID: "admin-id",
				},
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.Approve(ctx, tt.args.scheduleID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestSchedule_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	type args struct {
		scheduleID string
		public     bool
	}
	type want struct {
		hasErr bool
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
				schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
				schedule.StartAt = now().AddDate(0, 1, 0)
				err = db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				scheduleID: "schedule-id",
				public:     true,
			},
			want: want{
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, scheduleTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &schedule{db: db, now: now}
			err = db.Publish(ctx, tt.args.scheduleID, tt.args.public)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testSchedule(id, shopID, coordinatorID string, now time.Time) *entity.Schedule {
	schedule := &entity.Schedule{
		ID:              id,
		ShopID:          shopID,
		CoordinatorID:   coordinatorID,
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
