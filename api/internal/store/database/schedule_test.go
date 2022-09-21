package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchedule(t *testing.T) {
	assert.NotNil(t, NewSchedule(nil))
}

func TestSchedule_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, scheduleTable)

	type args struct {
		schedule *entity.Schedule
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				schedule: testSchedule("schedule-id", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				schedule := testSchedule("schedule-id", now())
				err = m.db.DB.Create(&schedule).Error
				require.NoError(t, err)
			},
			args: args{
				schedule: testSchedule("schedule-id", now()),
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

			err := m.dbDelete(ctx, scheduleTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &schedule{db: m.db, now: now}
			err = db.Create(ctx, tt.args.schedule)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testSchedule(id string, now time.Time) *entity.Schedule {
	return &entity.Schedule{
		ID:           id,
		Title:        "旬の夏野菜配信",
		Description:  "旬の夏野菜特集",
		ThumbnailURL: "https://and-period.jp/thumbnail01.png",
		StartAt:      now,
		EndAt:        now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
