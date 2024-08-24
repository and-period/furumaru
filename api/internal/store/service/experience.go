package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListExperiences(ctx context.Context, in *store.ListExperiencesInput) (entity.Experiences, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListExperiencesParams{
		Name:          in.Name,
		CoordinatorID: in.CoordinatorID,
		ProducerID:    in.ProducerID,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
	}
	var (
		experiences entity.Experiences
		total       int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		experiences, err = s.db.Experience.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Experience.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return experiences, total, nil
}

func (s *service) MultiGetExperiences(ctx context.Context, in *store.MultiGetExperiencesInput) (entity.Experiences, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	experiences, err := s.db.Experience.MultiGet(ctx, in.ExperienceIDs)
	return experiences, internalError(err)
}

func (s *service) MultiGetExperiencesByRevision(
	ctx context.Context, in *store.MultiGetExperiencesByRevisionInput,
) (entity.Experiences, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	experiences, err := s.db.Experience.MultiGetByRevision(ctx, in.ExperienceRevisionIDs)
	return experiences, internalError(err)
}

func (s *service) GetExperience(ctx context.Context, in *store.GetExperienceInput) (*entity.Experience, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	experience, err := s.db.Experience.Get(ctx, in.ExperienceID)
	return experience, internalError(err)
}

func (s *service) CreateExperience(ctx context.Context, in *store.CreateExperienceInput) (*entity.Experience, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	media := make(entity.MultiExperienceMedia, len(in.Media))
	for i := range in.Media {
		media[i] = entity.NewExperienceMedia(in.Media[i].URL, in.Media[i].IsThumbnail)
	}
	if err := media.Validate(); err != nil {
		return nil, fmt.Errorf("api: invalid media format: %s: %w", err.Error(), exception.ErrInvalidArgument)
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
		in := &user.GetProducerInput{
			ProducerID: in.ProducerID,
		}
		_, err = s.user.GetProducer(ectx, in)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid coordinator or producer: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewExperienceParams{
		CoordinatorID:         in.CoordinatorID,
		ProducerID:            in.ProducerID,
		TypeID:                in.TypeID,
		Title:                 in.Title,
		Description:           in.Description,
		Public:                in.Public,
		SoldOut:               in.SoldOut,
		Media:                 media,
		RecommendedPoints:     in.RecommendedPoints,
		PromotionVideoURL:     in.PromotionVideoURL,
		HostPrefectureCode:    in.HostPrefectureCode,
		HostCity:              in.HostCity,
		StartAt:               in.StartAt,
		EndAt:                 in.EndAt,
		PriceAdult:            in.PriceAdult,
		PriceJuniorHighSchool: in.PriceJuniorHighSchool,
		PriceElementarySchool: in.PriceElementarySchool,
		PricePreschool:        in.PricePreschool,
		PriceSenior:           in.PriceSenior,
	}
	experience, err := entity.NewExperience(params)
	if err != nil {
		return nil, fmt.Errorf("api: invalid experience: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err := experience.Validate(); err != nil {
		return nil, fmt.Errorf("api: invalid experience: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err := s.db.Experience.Create(ctx, experience); err != nil {
		return nil, internalError(err)
	}
	return experience, nil
}

func (s *service) UpdateExperience(ctx context.Context, in *store.UpdateExperienceInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if _, err := codes.ToPrefectureJapanese(in.HostPrefectureCode); err != nil {
		return fmt.Errorf("api: invalid host prefecture code: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	media := make(entity.MultiExperienceMedia, len(in.Media))
	for i := range in.Media {
		media[i] = entity.NewExperienceMedia(in.Media[i].URL, in.Media[i].IsThumbnail)
	}
	if err := media.Validate(); err != nil {
		return fmt.Errorf("api: invalid media format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &database.UpdateExperienceParams{
		TypeID:                in.TypeID,
		Title:                 in.Title,
		Description:           in.Description,
		Public:                in.Public,
		SoldOut:               in.SoldOut,
		Media:                 media,
		PriceAdult:            in.PriceAdult,
		PriceJuniorHighSchool: in.PriceJuniorHighSchool,
		PriceElementarySchool: in.PriceElementarySchool,
		PricePreschool:        in.PricePreschool,
		PriceSenior:           in.PriceSenior,
		RecommendedPoints:     in.RecommendedPoints,
		PromotionVideoURL:     in.PromotionVideoURL,
		HostPrefectureCode:    in.HostPrefectureCode,
		HostCity:              in.HostCity,
		StartAt:               in.StartAt,
		EndAt:                 in.EndAt,
	}
	err := s.db.Experience.Update(ctx, in.ExperienceID, params)
	return internalError(err)
}

func (s *service) DeleteExperience(ctx context.Context, in *store.DeleteExperienceInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Experience.Delete(ctx, in.ExperienceID)
	return internalError(err)
}
