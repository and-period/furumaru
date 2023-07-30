package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestGetBroadcastByScheduleID(t *testing.T) {
	t.Parallel()

	broadcast := &entity.Broadcast{
		ID:         "broadcast-id",
		ScheduleID: "schedule-id",
		Status:     entity.BroadcastStatusIdle,
		InputURL:   "rtmp://127.0.0.1:1935/app/instance",
		OutputURL:  "http://example.com/index.m3u8",
		CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetBroadcastByScheduleIDInput
		expect    *entity.Broadcast
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &store.GetBroadcastByScheduleIDInput{
				ScheduleID: "schedule-id",
			},
			expect:    broadcast,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetBroadcastByScheduleIDInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &store.GetBroadcastByScheduleIDInput{
				ScheduleID: "schedule-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetBroadcastByScheduleID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
