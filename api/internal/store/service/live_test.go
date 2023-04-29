package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	pivs "github.com/and-period/furumaru/api/pkg/ivs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMultiGetLives(t *testing.T) {
	t.Parallel()

	lives := entity.Lives{
		{
			ID:           "live-id01",
			ChannelArn:   "channel-arn",
			StreamKeyArn: "streamKey-arn",
		},
		{
			ID:           "live-id02",
			ChannelArn:   "channel-arn",
			StreamKeyArn: "streamKey-arn",
		},
	}
	channelIn := &pivs.GetChannelParams{
		Arn: "channel-arn",
	}

	streamIn := &pivs.GetStreamParams{
		ChannelArn: "channel-arn",
	}

	streamKeyIn := &pivs.GetStreamKeyParams{
		StreamKeyArn: "streamKey-arn",
	}

	channel := &types.Channel{
		Arn:            aws.String("channel-arn"),
		IngestEndpoint: aws.String("ingest-endpoint"),
		Name:           aws.String("配信チャンネル"),
		PlaybackUrl:    aws.String("playback-url"),
	}

	stream := &types.Stream{
		ChannelArn:  aws.String("channel-arn"),
		StreamId:    aws.String("stream-id"),
		ViewerCount: 100,
	}

	streamKey := &types.StreamKey{
		Arn:        aws.String("streamKey-arn"),
		ChannelArn: aws.String("channel-arn"),
		Value:      aws.String("streamKey-value"),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetLivesInput
		expect    entity.Lives
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().MultiGet(ctx, []string{"live-id01", "live-id02"}).Return(lives, nil)
				mocks.ivs.EXPECT().GetChannel(gomock.Any(), channelIn).Return(channel, nil).Times(2)
				mocks.ivs.EXPECT().GetStream(gomock.Any(), streamIn).Return(stream, nil).Times(2)
				mocks.ivs.EXPECT().GetStreamKey(gomock.Any(), streamKeyIn).Return(streamKey, nil).Times(2)
			},
			input: &store.MultiGetLivesInput{
				LiveIDs: []string{"live-id01", "live-id02"},
			},
			expect: entity.Lives{
				{
					ID:             "live-id01",
					ChannelArn:     "channel-arn",
					StreamKeyArn:   "streamKey-arn",
					ChannelName:    "配信チャンネル",
					IngestEndpoint: "ingest-endpoint",
					StreamKey:      "streamKey-value",
					StreamID:       "stream-id",
					PlaybackURL:    "playback-url",
					ViewerCount:    100,
				},
				{
					ID:             "live-id02",
					ChannelArn:     "channel-arn",
					StreamKeyArn:   "streamKey-arn",
					ChannelName:    "配信チャンネル",
					IngestEndpoint: "ingest-endpoint",
					StreamKey:      "streamKey-value",
					StreamID:       "stream-id",
					PlaybackURL:    "playback-url",
					ViewerCount:    100,
				},
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetLives(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestListByScheduleID(t *testing.T) {
	t.Parallel()

	lives := entity.Lives{
		{
			ID:           "live-id01",
			ScheduleID:   "schedule-id",
			ChannelArn:   "channel-arn",
			StreamKeyArn: "streamKey-arn",
		},
		{
			ID:           "live-id02",
			ScheduleID:   "schedule-id",
			ChannelArn:   "channel-arn",
			StreamKeyArn: "streamKey-arn",
		},
	}
	channelIn := &pivs.GetChannelParams{
		Arn: "channel-arn",
	}

	streamIn := &pivs.GetStreamParams{
		ChannelArn: "channel-arn",
	}

	streamKeyIn := &pivs.GetStreamKeyParams{
		StreamKeyArn: "streamKey-arn",
	}

	channel := &types.Channel{
		Arn:            aws.String("channel-arn"),
		IngestEndpoint: aws.String("ingest-endpoint"),
		Name:           aws.String("配信チャンネル"),
		PlaybackUrl:    aws.String("playback-url"),
	}

	stream := &types.Stream{
		ChannelArn:  aws.String("channel-arn"),
		StreamId:    aws.String("stream-id"),
		ViewerCount: 100,
	}

	streamKey := &types.StreamKey{
		Arn:        aws.String("streamKey-arn"),
		ChannelArn: aws.String("channel-arn"),
		Value:      aws.String("streamKey-value"),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListLivesByScheduleIDInput
		expect    entity.Lives
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().ListByScheduleID(ctx, "schedule-id").Return(lives, nil)
				mocks.ivs.EXPECT().GetChannel(gomock.Any(), channelIn).Return(channel, nil).Times(2)
				mocks.ivs.EXPECT().GetStream(gomock.Any(), streamIn).Return(stream, nil).Times(2)
				mocks.ivs.EXPECT().GetStreamKey(gomock.Any(), streamKeyIn).Return(streamKey, nil).Times(2)
			},
			input: &store.ListLivesByScheduleIDInput{
				ScheduleID: "schedule-id",
			},
			expect: entity.Lives{
				{
					ID:             "live-id01",
					ScheduleID:     "schedule-id",
					ChannelArn:     "channel-arn",
					StreamKeyArn:   "streamKey-arn",
					ChannelName:    "配信チャンネル",
					IngestEndpoint: "ingest-endpoint",
					StreamKey:      "streamKey-value",
					StreamID:       "stream-id",
					PlaybackURL:    "playback-url",
					ViewerCount:    100,
				},
				{
					ID:             "live-id02",
					ScheduleID:     "schedule-id",
					ChannelArn:     "channel-arn",
					StreamKeyArn:   "streamKey-arn",
					ChannelName:    "配信チャンネル",
					IngestEndpoint: "ingest-endpoint",
					StreamKey:      "streamKey-value",
					StreamID:       "stream-id",
					PlaybackURL:    "playback-url",
					ViewerCount:    100,
				},
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListLivesByScheduleID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetLive(t *testing.T) {
	t.Parallel()
	live := &entity.Live{
		ID:           "live-id",
		ChannelArn:   "channel-arn",
		StreamKeyArn: "streamKey-arn",
	}

	channelIn := &pivs.GetChannelParams{
		Arn: "channel-arn",
	}

	streamIn := &pivs.GetStreamParams{
		ChannelArn: "channel-arn",
	}

	streamKeyIn := &pivs.GetStreamKeyParams{
		StreamKeyArn: "streamKey-arn",
	}

	channel := &types.Channel{
		Arn:            aws.String("channel-arn"),
		IngestEndpoint: aws.String("ingest-endpoint"),
		Name:           aws.String("配信チャンネル"),
		PlaybackUrl:    aws.String("playback-url"),
	}

	stream := &types.Stream{
		ChannelArn:  aws.String("channel-arn"),
		StreamId:    aws.String("stream-id"),
		ViewerCount: 100,
	}

	streamKey := &types.StreamKey{
		Arn:        aws.String("streamKey-arn"),
		ChannelArn: aws.String("channel-arn"),
		Value:      aws.String("streamKey-value"),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetLiveInput
		expect    *entity.Live
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.ivs.EXPECT().GetChannel(gomock.Any(), channelIn).Return(channel, nil)
				mocks.ivs.EXPECT().GetStream(gomock.Any(), streamIn).Return(stream, nil)
				mocks.ivs.EXPECT().GetStreamKey(gomock.Any(), streamKeyIn).Return(streamKey, nil)
			},
			input: &store.GetLiveInput{
				LiveID: "live-id",
			},
			expect: &entity.Live{
				ID:             "live-id",
				ChannelArn:     "channel-arn",
				StreamKeyArn:   "streamKey-arn",
				ChannelName:    "配信チャンネル",
				IngestEndpoint: "ingest-endpoint",
				StreamKey:      "streamKey-value",
				StreamID:       "stream-id",
				PlaybackURL:    "playback-url",
				ViewerCount:    100,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetLiveInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(nil, assert.AnError)
			},
			input: &store.GetLiveInput{
				LiveID: "live-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to not found channel",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.ivs.EXPECT().GetChannel(gomock.Any(), channelIn).Return(nil, exception.ErrNotFound)
				mocks.ivs.EXPECT().GetStream(gomock.Any(), streamIn).Return(stream, nil)
				mocks.ivs.EXPECT().GetStreamKey(gomock.Any(), streamKeyIn).Return(streamKey, nil)
			},
			input: &store.GetLiveInput{
				LiveID: "live-id",
			},
			expect:    nil,
			expectErr: exception.ErrNotFound,
		},
		{
			name: "failed to not found stream",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.ivs.EXPECT().GetChannel(gomock.Any(), channelIn).Return(channel, nil)
				mocks.ivs.EXPECT().GetStream(gomock.Any(), streamIn).Return(nil, exception.ErrNotFound)
				mocks.ivs.EXPECT().GetStreamKey(gomock.Any(), streamKeyIn).Return(streamKey, nil)
			},
			input: &store.GetLiveInput{
				LiveID: "live-id",
			},
			expect:    nil,
			expectErr: exception.ErrNotFound,
		},
		{
			name: "failed to not found streamKey",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.ivs.EXPECT().GetChannel(gomock.Any(), channelIn).Return(channel, nil)
				mocks.ivs.EXPECT().GetStream(gomock.Any(), streamIn).Return(stream, nil)
				mocks.ivs.EXPECT().GetStreamKey(gomock.Any(), streamKeyIn).Return(nil, exception.ErrNotFound)
			},
			input: &store.GetLiveInput{
				LiveID: "live-id",
			},
			expect:    nil,
			expectErr: exception.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetLive(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
