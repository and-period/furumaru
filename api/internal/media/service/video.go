package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListVideos(ctx context.Context, in *media.ListVideosInput) (entity.Videos, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListVideosParams{
		Name:          in.Name,
		CoordinatorID: in.CoordinatorID,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
	}
	var (
		videos entity.Videos
		total  int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		videos, err = s.db.Video.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Video.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return videos, total, nil
}

func (s *service) GetVideo(ctx context.Context, in *media.GetVideoInput) (*entity.Video, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.Video{}, nil
}

func (s *service) CreateVideo(ctx context.Context, in *media.CreateVideoInput) (*entity.Video, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.Video{}, nil
}

func (s *service) UpdateVideo(ctx context.Context, in *media.UpdateVideoInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}

func (s *service) DeleteVideo(ctx context.Context, in *media.DeleteVideoInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}
