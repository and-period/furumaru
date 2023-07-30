package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcast(t *testing.T) {
	assert.NotNil(t, NewBroadcast(nil))
}

func TestBroadcast_GetByScheduleID(t *testing.T) {
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
	schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
	err = db.DB.Create(&schedule).Error

	b := testBroadcast("broadcast-id", "schedule-id", now())
	err = db.DB.Create(&b).Error

	type args struct {
		scheduleID string
	}
	type want struct {
		broadcast *entity.Broadcast
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
				scheduleID: "schedule-id",
			},
			want: want{
				broadcast: b,
				hasErr:    false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				scheduleID: "",
			},
			want: want{
				broadcast: nil,
				hasErr:    true,
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

			db := &broadcast{db: db, now: now}
			actual, err := db.GetByScheduleID(ctx, tt.args.scheduleID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.broadcast, actual)
		})
	}
}

func testBroadcast(broadcastID, scheduleID string, now time.Time) *entity.Broadcast {
	return &entity.Broadcast{
		ID:         broadcastID,
		ScheduleID: scheduleID,
		Status:     entity.BroadcastStatusIdle,
		InputURL:   "rtmp://127.0.0.1/1935/app/instance",
		OutputURL:  "http://example.com/index.m3u8",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
