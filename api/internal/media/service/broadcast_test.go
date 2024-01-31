package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListBroadcasts(t *testing.T) {
	t.Parallel()

	now := jst.Date(2023, 10, 20, 18, 30, 0, 0)
	params := &database.ListBroadcastsParams{
		ScheduleIDs:   []string{"schedule-id"},
		CoordinatorID: "coordinator-id",
		Limit:         30,
		Offset:        0,
		OnlyArchived:  true,
		Orders: []*database.ListBroadcastsOrder{
			{Key: entity.BroadcastOrderByUpdatedAt, OrderByASC: true},
		},
	}
	broadcasts := entity.Broadcasts{
		{
			ID:            "broadcast-id",
			ScheduleID:    "schedule-id",
			CoordinatorID: "coordinator-id",
			Status:        entity.BroadcastStatusIdle,
			InputURL:      "rtmp://127.0.0.1:1935/app/instance",
			OutputURL:     "http://example.com/index.m3u8",
			ArchiveURL:    "http://example.com/archive.mp4",
			CreatedAt:     now,
			UpdatedAt:     now,
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
				ScheduleIDs:   []string{"schedule-id"},
				CoordinatorID: "coordinator-id",
				Limit:         30,
				Offset:        0,
				OnlyArchived:  true,
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
				ScheduleIDs:   []string{"schedule-id"},
				CoordinatorID: "coordinator-id",
				Limit:         30,
				Offset:        0,
				OnlyArchived:  true,
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
				ScheduleIDs:   []string{"schedule-id"},
				CoordinatorID: "coordinator-id",
				Limit:         30,
				Offset:        0,
				OnlyArchived:  true,
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
		ID:            "broadcast-id",
		ScheduleID:    "schedule-id",
		CoordinatorID: "coordinator-id",
		Status:        entity.BroadcastStatusIdle,
		InputURL:      "rtmp://127.0.0.1:1935/app/instance",
		OutputURL:     "http://example.com/index.m3u8",
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
							ID:            broadcast.ID, // ignore
							ScheduleID:    "schedule-id",
							CoordinatorID: "coordinator-id",
							Type:          entity.BroadcastTypeNormal,
							Status:        entity.BroadcastStatusDisabled,
						}
						assert.Equal(t, expect, broadcast)
						return nil
					})
			},
			input: &media.CreateBroadcastInput{
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
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
				ScheduleID:    "schedule-id",
				CoordinatorID: "coordinator-id",
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

func TestUpdateBroadcastArchive(t *testing.T) {
	t.Parallel()
	broadcast := &entity.Broadcast{
		ID:         "broadcast-id",
		ScheduleID: "schdule-id",
		Status:     entity.BroadcastStatusDisabled,
	}
	params := &database.UpdateBroadcastParams{
		UploadBroadcastArchiveParams: &database.UploadBroadcastArchiveParams{
			ArchiveURL:   "http://example.com/archive.mp4",
			ArchiveFixed: true,
		},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UpdateBroadcastArchiveInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.Broadcast.EXPECT().Update(ctx, "broadcast-id", params).Return(nil)
			},
			input: &media.UpdateBroadcastArchiveInput{
				ScheduleID: "schedule-id",
				ArchiveURL: "http://example.com/archive.mp4",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UpdateBroadcastArchiveInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.UpdateBroadcastArchiveInput{
				ScheduleID: "schedule-id",
				ArchiveURL: "http://example.com/archive.mp4",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "broadcast is enabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusActive}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.UpdateBroadcastArchiveInput{
				ScheduleID: "schedule-id",
				ArchiveURL: "http://example.com/archive.mp4",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to update broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.db.Broadcast.EXPECT().Update(ctx, "broadcast-id", params).Return(assert.AnError)
			},
			input: &media.UpdateBroadcastArchiveInput{
				ScheduleID: "schedule-id",
				ArchiveURL: "http://example.com/archive.mp4",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateBroadcastArchive(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestPauseBroadcast(t *testing.T) {
	t.Parallel()
	now := time.Now()
	broadcast := &entity.Broadcast{
		Status:                 entity.BroadcastStatusActive,
		MediaLiveChannelID:     "12345678",
		MediaLiveRTMPInputName: "rtmp",
	}
	params := &medialive.CreateScheduleParams{
		ChannelID: "12345678",
		Settings: []*medialive.ScheduleSetting{{
			Name:       fmt.Sprintf("%s immediate-pause", jst.Format(now, time.DateTime)),
			ActionType: medialive.ScheduleActionTypePauseState,
			StartType:  medialive.ScheduleStartTypeImmediate,
			Reference:  string(medialive.PipelineIDPipeline0),
		}},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.PauseBroadcastInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.media.EXPECT().CreateSchedule(ctx, params).Return(nil)
			},
			input: &media.PauseBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.PauseBroadcastInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.PauseBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.PauseBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to activate static image",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.media.EXPECT().CreateSchedule(ctx, params).Return(assert.AnError)
			},
			input: &media.PauseBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.PauseBroadcast(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestUnpauseBroadcast(t *testing.T) {
	t.Parallel()
	now := time.Now()
	broadcast := &entity.Broadcast{
		Status:                 entity.BroadcastStatusActive,
		MediaLiveChannelID:     "12345678",
		MediaLiveRTMPInputName: "rtmp",
	}
	params := &medialive.CreateScheduleParams{
		ChannelID: "12345678",
		Settings: []*medialive.ScheduleSetting{{
			Name:       fmt.Sprintf("%s immediate-unpause", jst.Format(now, time.DateTime)),
			ActionType: medialive.ScheduleActionTypeUnpauseState,
			StartType:  medialive.ScheduleStartTypeImmediate,
		}},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.UnpauseBroadcastInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.media.EXPECT().CreateSchedule(ctx, params).Return(nil)
			},
			input: &media.UnpauseBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.UnpauseBroadcastInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.UnpauseBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.UnpauseBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to activate static image",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.media.EXPECT().CreateSchedule(ctx, params).Return(assert.AnError)
			},
			input: &media.UnpauseBroadcastInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UnpauseBroadcast(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestActivateBroadcastRTMP(t *testing.T) {
	t.Parallel()
	now := time.Now()
	broadcast := &entity.Broadcast{
		Status:                 entity.BroadcastStatusActive,
		MediaLiveChannelID:     "12345678",
		MediaLiveRTMPInputName: "rtmp",
	}
	params := &medialive.CreateScheduleParams{
		ChannelID: "12345678",
		Settings: []*medialive.ScheduleSetting{{
			Name:       fmt.Sprintf("%s immediate-input-rtmp", jst.Format(now, time.DateTime)),
			ActionType: medialive.ScheduleActionTypeInputSwitch,
			StartType:  medialive.ScheduleStartTypeImmediate,
			Reference:  broadcast.MediaLiveRTMPInputName,
		}},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.ActivateBroadcastRTMPInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.media.EXPECT().CreateSchedule(ctx, params).Return(nil)
			},
			input: &media.ActivateBroadcastRTMPInput{
				ScheduleID: "schedule-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.ActivateBroadcastRTMPInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.ActivateBroadcastRTMPInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.ActivateBroadcastRTMPInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to activate static image",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.media.EXPECT().CreateSchedule(ctx, params).Return(assert.AnError)
			},
			input: &media.ActivateBroadcastRTMPInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ActivateBroadcastRTMP(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestActivateBroadcastMP4(t *testing.T) {
	t.Parallel()
	now := time.Now()
	broadcast := &entity.Broadcast{
		Status:                entity.BroadcastStatusActive,
		MediaLiveChannelID:    "12345678",
		MediaLiveMP4InputName: "mp4",
	}
	params := &medialive.CreateScheduleParams{
		ChannelID: "12345678",
		Settings: []*medialive.ScheduleSetting{{
			Name:       fmt.Sprintf("%s immediate-input-mp4", jst.Format(now, time.DateTime)),
			ActionType: medialive.ScheduleActionTypeInputSwitch,
			StartType:  medialive.ScheduleStartTypeImmediate,
			Reference:  broadcast.MediaLiveMP4InputName,
			Source:     "s3://example.mp4",
		}},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.ActivateBroadcastMP4Input
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.tmp.EXPECT().ReplaceURLToS3URI("http://example.com/example.mp4").Return("s3://example.mp4", nil)
				mocks.media.EXPECT().CreateSchedule(ctx, params).Return(nil)
			},
			input: &media.ActivateBroadcastMP4Input{
				ScheduleID: "schedule-id",
				InputURL:   "http://example.com/example.mp4",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.ActivateBroadcastMP4Input{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.ActivateBroadcastMP4Input{
				ScheduleID: "schedule-id",
				InputURL:   "http://example.com/example.mp4",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.ActivateBroadcastMP4Input{
				ScheduleID: "schedule-id",
				InputURL:   "http://example.com/example.mp4",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to replace s3 uri",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.tmp.EXPECT().ReplaceURLToS3URI("http://example.com/example.mp4").Return("", assert.AnError)
			},
			input: &media.ActivateBroadcastMP4Input{
				ScheduleID: "schedule-id",
				InputURL:   "http://example.com/example.mp4",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.tmp.EXPECT().ReplaceURLToS3URI("http://example.com/example.mp4").Return("s3://example.mp4", nil)
				mocks.media.EXPECT().CreateSchedule(ctx, params).Return(assert.AnError)
			},
			input: &media.ActivateBroadcastMP4Input{
				ScheduleID: "schedule-id",
				InputURL:   "http://example.com/example.mp4",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ActivateBroadcastMP4(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestActivateBroadcastStaticImage(t *testing.T) {
	t.Parallel()
	broadcast := &entity.Broadcast{
		Status:             entity.BroadcastStatusActive,
		MediaLiveChannelID: "12345678",
	}
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: "schedule-id",
	}
	schedule := &sentity.Schedule{
		ID:       "schedule-id",
		ImageURL: "http://example.com/image.png",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.ActivateBroadcastStaticImageInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
				mocks.storage.EXPECT().ReplaceURLToS3URI("http://example.com/image.png").Return("s3://image.png", nil)
				mocks.media.EXPECT().ActivateStaticImage(ctx, "12345678", "s3://image.png").Return(nil)
			},
			input: &media.ActivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.ActivateBroadcastStaticImageInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.ActivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.ActivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to get schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(nil, assert.AnError)
			},
			input: &media.ActivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to replace s3 uri",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
				mocks.storage.EXPECT().ReplaceURLToS3URI("http://example.com/image.png").Return("", assert.AnError)
			},
			input: &media.ActivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to activate static image",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
				mocks.storage.EXPECT().ReplaceURLToS3URI("http://example.com/image.png").Return("s3://image.png", nil)
				mocks.media.EXPECT().ActivateStaticImage(ctx, "12345678", "s3://image.png").Return(assert.AnError)
			},
			input: &media.ActivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ActivateBroadcastStaticImage(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeactivateBroadcastStaticImage(t *testing.T) {
	t.Parallel()
	broadcast := &entity.Broadcast{
		Status:             entity.BroadcastStatusActive,
		MediaLiveChannelID: "12345678",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *media.DeactivateBroadcastStaticImageInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.media.EXPECT().DeactivateStaticImage(ctx, "12345678").Return(nil)
			},
			input: &media.DeactivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &media.DeactivateBroadcastStaticImageInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get broadcast",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(nil, assert.AnError)
			},
			input: &media.DeactivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "broadcast is disabled",
			setup: func(ctx context.Context, mocks *mocks) {
				broadcast := &entity.Broadcast{Status: entity.BroadcastStatusDisabled}
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
			},
			input: &media.DeactivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to deactivate static image",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Broadcast.EXPECT().GetByScheduleID(ctx, "schedule-id").Return(broadcast, nil)
				mocks.media.EXPECT().DeactivateStaticImage(ctx, "12345678").Return(assert.AnError)
			},
			input: &media.DeactivateBroadcastStaticImageInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeactivateBroadcastStaticImage(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
