package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestCreateBrodcastViewerLog(t *testing.T) {
	t.Parallel()

	broadcast := &entity.Broadcast{
		ID:            "broadcast-id",
		ScheduleID:    "schedule-id",
		CoordinatorID: "coordinator-id",
		Status:        entity.BroadcastStatusIdle,
		InputURL:      "rtmp://127.0.0.1:1935/app/instance",
		OutputURL:     "http://example.com/index.m3u8",
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	log := &entity.BroadcastViewerLog{
		BroadcastID: "broadcast-id",
		SessionID:   "session-id",
		UserID:      "user-id",
		UserAgent:   "user-agent",
		ClientIP:    "127.0.0.1",
	}

	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.CreateBroadcastViewerLogInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.BroadcastViewerLog.EXPECT().Create(ctx, log).Return(nil)
			},
			input: &media.CreateBroadcastViewerLogInput{
				ScheduleID: "schedule-id",
				SessionID:  "session-id",
				UserID:     "user-id",
				UserAgent:  "user-agent",
				ClientIP:   "127.0.0.1",
			},
			expect: nil,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.CreateBroadcastViewerLogInput{
				ScheduleID: "schedule-id",
				SessionID:  "session-id",
				UserID:     "user-id",
				UserAgent:  "user-agent",
				ClientIP:   "127.0.0.1",
			},
			expect: exception.ErrInternal,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.CreateBroadcastViewerLogInput{
				ScheduleID: "schedule-id",
				SessionID:  "session-id",
				UserID:     "user-id",
				UserAgent:  "user-agent",
				ClientIP:   "127.0.0.1",
			},
			expect: nil,
		},
		{
			name: "failed to create broadacast viewer log",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.BroadcastViewerLog.EXPECT().Create(ctx, log).Return(assert.AnError)
			},
			input: &media.CreateBroadcastViewerLogInput{
				ScheduleID: "schedule-id",
				SessionID:  "session-id",
				UserID:     "user-id",
				UserAgent:  "user-agent",
				ClientIP:   "127.0.0.1",
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.CreateBroadcastViewerLog(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}
