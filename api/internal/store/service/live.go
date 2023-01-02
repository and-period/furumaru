package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	pivs "github.com/and-period/furumaru/api/pkg/ivs"
	"github.com/aws/aws-sdk-go-v2/aws"
	ivs "github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
	"golang.org/x/sync/errgroup"
)

func (s *service) GetLive(ctx context.Context, in *store.GetLiveInput) (*entity.Live, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}

	live, err := s.db.Live.Get(ctx, in.LiveID)
	if err != nil {
		return nil, exception.InternalError(err)
	}

	var (
		channel   *types.Channel
		stream    *types.Stream
		streamKey *types.StreamKey
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &pivs.GetChannelParams{
			Arn: live.ChannelArn,
		}
		channel, err = s.ivs.GetChannel(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &pivs.GetStreamParams{
			ChannelArn: live.ChannelArn,
		}
		stream, err = s.ivs.GetStream(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &pivs.GetStreamKeyParams{
			StreamKeyArn: live.StreamKeyArn,
		}
		streamKey, err = s.ivs.GetStreamKey(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
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

func (s *service) UpdateLivePublic(ctx context.Context, in *store.UpdateLivePublicInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	live, err := s.db.Live.Get(ctx, in.LiveID)
	if err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateLivePublicParams{
		Published: in.Published,
		Canceled:  in.Canceled,
	}
	ivs := func(ctx context.Context) (*ivs.CreateChannelOutput, error) {
		if live.ChannelArn != "" {
			return nil, exception.ErrAlreadyExists
		}
		in := &pivs.CreateChannelParams{
			LatencyMode: "NORMAL",
			Name:        in.ChannelName,
			ChannelType: "BASIC",
		}
		cout, err := s.ivs.CreateChannel(ctx, in)
		return cout, err
	}
	err = s.db.Live.UpdateLivePublic(ctx, in.LiveID, params, ivs)
	return exception.InternalError(err)
}
