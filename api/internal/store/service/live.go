package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/ivs"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func (s *service) GetLive(ctx context.Context, in *store.GetLiveInput) (*entity.Live, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}

	live, err := s.db.Live.Get(ctx, in.LiveID)
	if err != nil {
		return nil, exception.InternalError(err)
	}

	channelIn := &ivs.GetChannelParams{
		Arn: live.ChannelArn,
	}
	channel, err := s.ivs.GetChannel(ctx, channelIn)
	if err != nil {
		return nil, exception.InternalError(err)
	}

	streamIn := &ivs.GetStreamParams{
		ChannelArn: live.ChannelArn,
	}
	stream, err := s.ivs.GetStream(ctx, streamIn)
	if err != nil {
		return nil, exception.InternalError(err)
	}

	streamKeyIn := &ivs.GetStreamKeyParams{
		StreamKeyArn: live.StreamKeyArn,
	}
	streamKey, err := s.ivs.GetStreamKey(ctx, streamKeyIn)
	if err != nil {
		return nil, exception.InternalError(err)
	}

	fillIvsParams := &entity.FillLiveIvsParams{
		ChannelName:    aws.ToString(channel.Name),
		IngestEndpoint: aws.ToString(channel.IngestEndpoint),
		StreamKey:      aws.ToString(streamKey.Value),
		PlaybackURL:    aws.ToString(channel.PlaybackUrl),
		StreamID:       aws.ToString(stream.StreamId),
		ViewerCount:    aws.ToInt64(&stream.ViewerCount),
	}
	live.FillIVS(*fillIvsParams)
	return live, nil
}
