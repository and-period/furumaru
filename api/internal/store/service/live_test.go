package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	pivs "github.com/and-period/furumaru/api/pkg/ivs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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

func TestUpdateLivePublic(t *testing.T) {
	t.Parallel()

	live := &entity.Live{
		ID: "live-id",
	}

	channelIn := &pivs.CreateChannelParams{
		LatencyMode: "NORMAL",
		Name:        "channel-name",
		ChannelType: "BASIC",
	}

	channelOut := &ivs.CreateChannelOutput{
		Channel: &types.Channel{
			Arn:  aws.String("channel-arn"),
			Name: aws.String("channel-name"),
		},
		StreamKey: &types.StreamKey{
			Arn:        aws.String("streamKey-arn"),
			ChannelArn: aws.String("channel-name"),
		},
	}

	params := &database.UpdateLivePublicParams{
		Published:    true,
		Canceled:     false,
		ChannelArn:   "channel-arn",
		StreamKeyArn: "streamKey-arn",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateLivePublicInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.ivs.EXPECT().CreateChannel(gomock.Any(), channelIn).Return(channelOut, nil)
				mocks.db.Live.EXPECT().UpdatePublic(ctx, "live-id", params).Return(nil)
			},
			input: &store.UpdateLivePublicInput{
				LiveID:      "live-id",
				Published:   true,
				Canceled:    false,
				ChannelName: "channel-name",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateLivePublicInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(nil, assert.AnError)
			},
			input: &store.UpdateLivePublicInput{
				LiveID:      "live-id",
				Published:   true,
				Canceled:    false,
				ChannelName: "channel-name",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to createa channel",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.ivs.EXPECT().CreateChannel(gomock.Any(), channelIn).Return(nil, assert.AnError)
			},
			input: &store.UpdateLivePublicInput{
				LiveID:      "live-id",
				Published:   true,
				Canceled:    false,
				ChannelName: "channel-name",
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to update live",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Live.EXPECT().Get(ctx, "live-id").Return(live, nil)
				mocks.ivs.EXPECT().CreateChannel(gomock.Any(), channelIn).Return(channelOut, nil)
				mocks.db.Live.EXPECT().UpdatePublic(ctx, "live-id", params).Return(assert.AnError)
			},
			input: &store.UpdateLivePublicInput{
				LiveID:      "live-id",
				Published:   true,
				Canceled:    false,
				ChannelName: "channel-name",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateLivePublic(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
