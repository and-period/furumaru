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

	schedules := make(entity.Schedules, 3)
	schedules[0] = testSchedule(entity.ScheduleTypeNotification, "schedule-id01", now().Add(-time.Hour))
	schedules[1] = testSchedule(entity.ScheduleTypeNotification, "schedule-id02", now())
	schedules[2] = testSchedule(entity.ScheduleTypeNotification, "schedule-id03", now())
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
					Types:    []entity.ScheduleType{entity.ScheduleTypeNotification},
					Statuses: []entity.ScheduleStatus{entity.ScheduleStatusProcessing},
					Since:    now(),
					Until:    now().Add(time.Hour),
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

func TestSchedule_UpsertProcessing(t *testing.T) {
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
		schedule *entity.Schedule
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
			name:  "success create",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				schedule: testSchedule(entity.ScheduleTypeNotification, "schedule-id", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "success update",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				schedule := testSchedule(entity.ScheduleTypeNotification, "schedule-id", now().Add(-15*time.Minute))
				err := db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				schedule: testSchedule(entity.ScheduleTypeNotification, "schedule-id", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "not executable",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				schedule := testSchedule(entity.ScheduleTypeNotification, "schedule-id", now())
				err := db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				schedule: testSchedule(entity.ScheduleTypeNotification, "schedule-id", now()),
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
			err = db.UpsertProcessing(ctx, tt.args.schedule)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestSchedule_UpdateDone(t *testing.T) {
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
		messageType entity.ScheduleType
		messageID   string
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
				schedule := testSchedule(entity.ScheduleTypeNotification, "schedule-id", now())
				schedule.Status = entity.ScheduleStatusProcessing
				err := db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				messageType: entity.ScheduleTypeNotification,
				messageID:   "schedule-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				messageType: entity.ScheduleTypeNotification,
				messageID:   "schedule-id",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "already done",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				schedule := testSchedule(entity.ScheduleTypeNotification, "schedule-id", now())
				schedule.Status = entity.ScheduleStatusDone
				err := db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				messageType: entity.ScheduleTypeNotification,
				messageID:   "schedule-id",
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
			err = db.UpdateDone(ctx, tt.args.messageType, tt.args.messageID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestSchedule_UpdateCancel(t *testing.T) {
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
		messageType entity.ScheduleType
		messageID   string
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
				schedule := testSchedule(entity.ScheduleTypeNotification, "schedule-id", now().Add(-15*time.Minute))
				schedule.Status = entity.ScheduleStatusProcessing
				schedule.Count = 2
				err := db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				messageType: entity.ScheduleTypeNotification,
				messageID:   "schedule-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				messageType: entity.ScheduleTypeNotification,
				messageID:   "schedule-id",
			},
			want: want{
				hasErr: true,
			},
		},
		{
			name: "should not cancel",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				schedule := testSchedule(entity.ScheduleTypeNotification, "schedule-id", now())
				schedule.Status = entity.ScheduleStatusDone
				err := db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				messageType: entity.ScheduleTypeNotification,
				messageID:   "schedule-id",
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
			err = db.UpdateCancel(ctx, tt.args.messageType, tt.args.messageID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testSchedule(typ entity.ScheduleType, id string, now time.Time) *entity.Schedule {
	return &entity.Schedule{
		MessageType: typ,
		MessageID:   id,
		Status:      entity.ScheduleStatusProcessing,
		Count:       1,
		SentAt:      now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
