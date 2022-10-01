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

func TestLive(t *testing.T) {
	assert.NotNil(t, NewLive(nil))
}

func TestLive_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 8, 1, 0, 0, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, liveTable)
	schedule := testSchedule("schedule-id", now())
	err = m.db.DB.Create(&schedule).Error

	type args struct {
		live *entity.Live
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
				live: testLive("live-id", "schedule-id", "producer-id", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				live := testLive("live-id", "schedule-id", "producer-id", now())
				err = m.db.DB.Create(&live).Error
				require.NoError(t, err)
			},
			args: args{
				live: testLive("live-id", "schedule-id", "producer-id", now()),
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

			err := m.dbDelete(ctx, liveTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &live{db: m.db, now: now}
			err = db.Create(ctx, tt.args.live)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testLive(id, scheduleID, producerID string, now time.Time) *entity.Live {
	l := &entity.Live{
		ID:          id,
		ScheduleID:  scheduleID,
		Title:       "配信のタイトル",
		Description: "配信の説明",
		ProducerID:  producerID,
		StartAt:     now,
		EndAt:       now,
		Recommends:  []string{"product-id1", "product-id2"},
	}
	_ = l.FillJSON()
	return l
}
