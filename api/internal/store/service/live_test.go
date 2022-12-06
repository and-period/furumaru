package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/ivs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
	"github.com/stretchr/testify/assert"
)

func TestGetLive(t *testing.T) {
	t.Parallel()
	live := &entity.Live{
		ID:           "live-id",
		ChannelArn:   "channel-arn",
		StreamKeyArn: "streamKey-arn",
	}

	channelIn := &ivs.GetChannelParams{
		Arn: "channel-arn",
	}

	streamIn := &ivs.GetStreamParams{
		ChannelArn: "channel-arn",
	}

	streamKeyIn := &ivs.GetStreamKeyParams{
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
				mocks.ivs.EXPECT().GetChannel(ctx, channelIn).Return(channel, nil)
				mocks.ivs.EXPECT().GetStream(ctx, streamIn).Return(stream, nil)
				mocks.ivs.EXPECT().GetStreamKey(ctx, streamKeyIn).Return(streamKey, nil)
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
