package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/geolocation"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListExperiences(ctx context.Context, in *store.ListExperiencesInput) (entity.Experiences, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListExperiencesParams{
		Name:           in.Name,
		HostPrefecture: in.PrefectureCode,
		ShopID:         in.ShopID,
		ProducerID:     in.ProducerID,
		OnlyPublished:  in.OnlyPublished,
		ExcludeDeleted: in.ExcludeDeleted,
		Limit:          int(in.Limit),
		Offset:         int(in.Offset),
	}
	if in.ExcludeFinished {
		params.EndAtGte = s.now()
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

func (s *service) ListExperiencesByGeolocation(
	ctx context.Context, in *store.ListExperiencesByGeolocationInput,
) (entity.Experiences, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.ListExperiencesByGeolocationParams{
		ShopID:         in.ShopID,
		ProducerID:     in.ProducerID,
		Longitude:      in.Longitude,
		Latitude:       in.Latitude,
		Radius:         in.Radius,
		OnlyPublished:  in.OnlyPublished,
		ExcludeDeleted: in.ExcludeDeleted,
	}
	if in.ExcludeFinished {
		params.EndAtGte = s.now()
	}
	experiences, err := s.db.Experience.ListByGeolocation(ctx, params)
	return experiences, internalError(err)
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
		_, err = s.db.Shop.Get(ectx, in.ShopID)
		return
	})
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
	if errors.Is(err, exception.ErrNotFound) || errors.Is(err, database.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid coordinator or producer: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	prefecture, err := codes.ToPrefectureJapanese(in.HostPrefectureCode)
	if err != nil {
		return nil, fmt.Errorf("api: invalid host prefecture code: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	locationIn := &geolocation.GetGeolocationInput{
		Address: &geolocation.Address{
			PostalCode:   in.HostPostalCode,
			Prefecture:   prefecture,
			City:         in.HostCity,
			AddressLine1: in.HostAddressLine1,
			AddressLine2: in.HostAddressLine2,
		},
	}
	location, err := s.geolocation.GetGeolocation(ctx, locationIn)
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewExperienceParams{
		ShopID:                in.ShopID,
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
		Duration:              in.Duration,
		Direction:             in.Direction,
		BusinessOpenTime:      in.BusinessOpenTime,
		BusinessCloseTime:     in.BusinessCloseTime,
		HostPostalCode:        in.HostPostalCode,
		HostPrefectureCode:    in.HostPrefectureCode,
		HostCity:              in.HostCity,
		HostAddressLine1:      in.HostAddressLine1,
		HostAddressLine2:      in.HostAddressLine2,
		HostLongitude:         location.Longitude,
		HostLatitude:          location.Latitude,
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
	prefecture, err := codes.ToPrefectureJapanese(in.HostPrefectureCode)
	if err != nil {
		return fmt.Errorf("api: invalid host prefecture code: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	locationIn := &geolocation.GetGeolocationInput{
		Address: &geolocation.Address{
			PostalCode:   in.HostPostalCode,
			Prefecture:   prefecture,
			City:         in.HostCity,
			AddressLine1: in.HostAddressLine1,
			AddressLine2: in.HostAddressLine2,
		},
	}
	location, err := s.geolocation.GetGeolocation(ctx, locationIn)
	if err != nil {
		return internalError(err)
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
		Duration:              in.Duration,
		Direction:             in.Direction,
		BusinessOpenTime:      in.BusinessOpenTime,
		BusinessCloseTime:     in.BusinessCloseTime,
		HostPostalCode:        in.HostPostalCode,
		HostPrefectureCode:    in.HostPrefectureCode,
		HostCity:              in.HostCity,
		HostAddressLine1:      in.HostAddressLine1,
		HostAddressLine2:      in.HostAddressLine2,
		HostLongitude:         location.Longitude,
		HostLatitude:          location.Latitude,
		StartAt:               in.StartAt,
		EndAt:                 in.EndAt,
	}
	err = s.db.Experience.Update(ctx, in.ExperienceID, params)
	return internalError(err)
}

func (s *service) DeleteExperience(ctx context.Context, in *store.DeleteExperienceInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	videosIn := &media.ListExperienceVideosInput{
		ExperienceID: in.ExperienceID,
	}
	videos, err := s.media.ListExperienceVideos(ctx, videosIn)
	if err != nil {
		return internalError(err)
	}
	if len(videos) > 0 {
		return fmt.Errorf("service: experience has videos: %w", exception.ErrFailedPrecondition)
	}
	err = s.db.Experience.Delete(ctx, in.ExperienceID)
	return internalError(err)
}
