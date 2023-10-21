package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListBroadcasts(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	params := &database.ListBroadcastsParams{
		Limit:        30,
		Offset:       0,
		OnlyArchived: true,
		Orders: []*database.ListBroadcastsOrder{
			{Key: entity.BroadcastOrderByUpdatedAt, OrderByASC: true},
		},
	}
	broadcasts := entity.Broadcasts{
		{
			ID:         "broadcast-id",
			ScheduleID: "schedule-id",
			Status:     entity.BroadcastStatusIdle,
			InputURL:   "rtmp://127.0.0.1:1935/app/instance",
			OutputURL:  "http://example.com/index.m3u8",
			ArchiveURL: "http://example.com/archive.mp4",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *media.ListBroadcastsInput
		expect      entity.Broadcasts
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().List(gomock.Any(), params).Return(broadcasts, nil)
				mocks.db.Broadcast.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &media.ListBroadcastsInput{
				Limit:        30,
				Offset:       0,
				OnlyArchived: true,
				Orders: []*media.ListBroadcastsOrder{
					{Key: entity.BroadcastOrderByUpdatedAt, OrderByASC: true},
				},
			},
			expect:      broadcasts,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &media.ListBroadcastsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list broadcasts",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Broadcast.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &media.ListBroadcastsInput{
				Limit:        30,
				Offset:       0,
				OnlyArchived: true,
				Orders: []*media.ListBroadcastsOrder{
					{Key: entity.BroadcastOrderByUpdatedAt, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count broadcasts",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().List(gomock.Any(), params).Return(broadcasts, nil)
				mocks.db.Broadcast.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &media.ListBroadcastsInput{
				Limit:        30,
				Offset:       0,
				OnlyArchived: true,
				Orders: []*media.ListBroadcastsOrder{
					{Key: entity.BroadcastOrderByUpdatedAt, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListBroadcasts(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

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
		input     *media.GetBroadcastByScheduleIDInput
		expect    *entity.Broadcast
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.GetBroadcastByScheduleIDInput{
				ScheduleID: "schedule-id",
			},
			expect:    broadcast,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.GetBroadcastByScheduleIDInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.GetBroadcastByScheduleIDInput{
				ScheduleID: "schedule-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
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

func TestCreateBroadcast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.CreateBroadcastInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, broadcast *entity.Broadcast) error {
						expect := &entity.Broadcast{
							ID:         broadcast.ID, // ignore
							ScheduleID: "schedule-id",
							Type:       entity.BroadcastTypeNormal,
							Status:     entity.BroadcastStatusDisabled,
						}
						assert.Equal(t, expect, broadcast)
						return nil
					})
			},
			input: &media.CreateBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.CreateBroadcastInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &media.CreateBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateBroadcast(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
