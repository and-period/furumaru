package database

import (
	"context"
	"fmt"
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
		lives    entity.Lives
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
				lives:    testLives("live-id", "schedule-id", "producer-id", now(), 3),
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
				lives:    testLives("live-id", "schedule-id", "producer-id", now(), 3),
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
			err = db.Create(ctx, tt.args.schedule, tt.args.lives)
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

func testLives(id, scheduleID, producerID string, now time.Time, length int) entity.Lives {
	lives := make(entity.Lives, length)

	for i := 0; i < length; i++ {
		lives[i] = &entity.Live{
			ID:          fmt.Sprintf("%s-%2d", id, i),
			ScheduleID:  scheduleID,
			Title:       "配信のタイトル",
			Description: "配信の説明",
			ProducerID:  producerID,
			StartAt:     now,
			EndAt:       now,
			Recommends:  []string{"product-id1", "product-id2"},
		}
	}

	return lives
}
