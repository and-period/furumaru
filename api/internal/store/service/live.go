package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/ivs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
	"golang.org/x/sync/errgroup"
)

func (s *service) MultiGetLives(ctx context.Context, in *store.MultiGetLivesInput) (entity.Lives, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}

	lives, err := s.db.Live.MultiGet(ctx, in.LiveIDs)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	for i := range lives {
		err = s.getIVSDetails(ctx, lives[i])
		if err != nil {
			return nil, exception.InternalError(err)
		}
	}
	return lives, exception.InternalError(err)
}

func (s *service) ListLivesByScheduleID(ctx context.Context, in *store.ListLivesByScheduleIDInput) (entity.Lives, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}

	lives, err := s.db.Live.ListByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	for i := range lives {
		err = s.getIVSDetails(ctx, lives[i])
		if err != nil {
			return nil, exception.InternalError(err)
		}
	}
	return lives, exception.InternalError(err)
}

func (s *service) GetLive(ctx context.Context, in *store.GetLiveInput) (*entity.Live, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}

	live, err := s.db.Live.Get(ctx, in.LiveID)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	err = s.getIVSDetails(ctx, live)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	return live, exception.InternalError(err)
}

func (s *service) getIVSDetails(ctx context.Context, live *entity.Live) (err error) {
	var (
		channel   *types.Channel
		stream    *types.Stream
		streamKey *types.StreamKey
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &ivs.GetChannelParams{
			Arn: live.ChannelArn,
		}
		channel, err = s.ivs.GetChannel(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &ivs.GetStreamParams{
			ChannelArn: live.ChannelArn,
		}
		stream, err = s.ivs.GetStream(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &ivs.GetStreamKeyParams{
			StreamKeyArn: live.StreamKeyArn,
		}
		streamKey, err = s.ivs.GetStreamKey(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		return exception.InternalError(err)
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
	return nil
}
