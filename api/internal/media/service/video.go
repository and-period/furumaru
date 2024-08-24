package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
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
	video, err := s.db.Video.Get(ctx, in.VideoID)
	return video, internalError(err)
}

func (s *service) CreateVideo(ctx context.Context, in *media.CreateVideoInput) (*entity.Video, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &user.GetCoordinatorInput{
			CoordinatorID: in.CoordinatorID,
		}
		_, err = s.user.GetCoordinator(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.MultiGetProductsInput{
			ProductIDs: in.ProductIDs,
		}
		products, err := s.store.MultiGetProducts(ectx, in)
		if err != nil {
			return err
		}
		if len(products) != len(in.ProductIDs) {
			return fmt.Errorf("service: product not found: %w", exception.ErrNotFound)
		}
		return nil
	})
	eg.Go(func() (err error) {
		in := &store.MultiGetExperiencesInput{
			ExperienceIDs: in.ExperienceIDs,
		}
		experiences, err := s.store.MultiGetExperiences(ectx, in)
		if err != nil {
			return err
		}
		if len(experiences) != len(in.ExperienceIDs) {
			return fmt.Errorf("service: experience not found: %w", exception.ErrNotFound)
		}
		return nil
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("service: invalid request: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}

	params := &entity.NewVideoParams{
		CoordinatorID: in.CoordinatorID,
		ProductIDs:    in.ProductIDs,
		ExperienceIDs: in.ExperienceIDs,
		Title:         in.Title,
		Description:   in.Description,
		ThumbnailURL:  in.ThumbnailURL,
		VideoURL:      in.VideoURL,
		Public:        in.Public,
		Limited:       in.Limited,
		PublishedAt:   in.PublishedAt,
	}
	video := entity.NewVideo(params)
	if err := s.db.Video.Create(ctx, video); err != nil {
		return nil, internalError(err)
	}
	return video, nil
}

func (s *service) UpdateVideo(ctx context.Context, in *media.UpdateVideoInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &store.MultiGetProductsInput{
			ProductIDs: in.ProductIDs,
		}
		products, err := s.store.MultiGetProducts(ectx, in)
		if err != nil {
			return err
		}
		if len(products) != len(in.ProductIDs) {
			return fmt.Errorf("service: product not found: %w", exception.ErrNotFound)
		}
		return nil
	})
	eg.Go(func() (err error) {
		in := &store.MultiGetExperiencesInput{
			ExperienceIDs: in.ExperienceIDs,
		}
		experiences, err := s.store.MultiGetExperiences(ectx, in)
		if err != nil {
			return err
		}
		if len(experiences) != len(in.ExperienceIDs) {
			return fmt.Errorf("service: experience not found: %w", exception.ErrNotFound)
		}
		return nil
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("service: invalid request: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}

	params := &database.UpdateVideoParams{
		Title:         in.Title,
		Description:   in.Description,
		ProductIDs:    in.ProductIDs,
		ExperienceIDs: in.ExperienceIDs,
		ThumbnailURL:  in.ThumbnailURL,
		VideoURL:      in.VideoURL,
		Public:        in.Public,
		Limited:       in.Limited,
		PublishedAt:   in.PublishedAt,
	}
	err = s.db.Video.Update(ctx, in.VideoID, params)
	return internalError(err)
}

func (s *service) DeleteVideo(ctx context.Context, in *media.DeleteVideoInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Video.Delete(ctx, in.VideoID)
	return internalError(err)
}
