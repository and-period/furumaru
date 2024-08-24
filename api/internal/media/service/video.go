package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

func (s *service) ListVideos(ctx context.Context, in *media.ListVideosInput) (entity.Videos, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	// TODO: 詳細の実装
	return entity.Videos{}, 0, nil
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
