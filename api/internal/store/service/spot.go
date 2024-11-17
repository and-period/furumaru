package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListSpots(ctx context.Context, in *store.ListSpotsInput) (entity.Spots, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListSpotsParams{
		Name:            in.Name,
		UserID:          in.UserID,
		ExcludeApproved: in.ExcludeApproved,
		ExcludeDisabled: in.ExcludeDisabled,
		Limit:           int(in.Limit),
		Offset:          int(in.Offset),
	}
	var (
		spots entity.Spots
		total int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		spots, err = s.db.Spot.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Spot.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return spots, total, nil
}

func (s *service) ListSpotsByGeolocation(ctx context.Context, in *store.ListSpotsByGeolocationInput) (entity.Spots, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.ListSpotsByGeolocationParams{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
		Radius:    in.Radius,
	}
	spots, err := s.db.Spot.ListByGeolocation(ctx, params)
	return spots, internalError(err)
}

func (s *service) GetSpot(ctx context.Context, in *store.GetSpotInput) (*entity.Spot, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	spot, err := s.db.Spot.Get(ctx, in.SpotID)
	return spot, internalError(err)
}

func (s *service) CreateSpotByUser(ctx context.Context, in *store.CreateSpotByUserInput) (*entity.Spot, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.SpotParams{
		UserID:       in.UserID,
		Name:         in.Name,
		Description:  in.Description,
		ThumbnailURL: in.ThumbnailURL,
		Longitude:    in.Longitude,
		Latitude:     in.Latitude,
	}
	spot, err := entity.NewSpotByUser(params)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create spot. err=%s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err := s.db.Spot.Create(ctx, spot); err != nil {
		return nil, internalError(err)
	}
	return spot, nil
}

func (s *service) CreateSpotByAdmin(ctx context.Context, in *store.CreateSpotByAdminInput) (*entity.Spot, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.SpotParams{
		UserID:       in.AdminID,
		Name:         in.Name,
		Description:  in.Description,
		ThumbnailURL: in.ThumbnailURL,
		Longitude:    in.Longitude,
		Latitude:     in.Latitude,
	}
	spot, err := entity.NewSpotByAdmin(params)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create spot. err=%s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err := s.db.Spot.Create(ctx, spot); err != nil {
		return nil, internalError(err)
	}
	return spot, nil
}

func (s *service) UpdateSpot(ctx context.Context, in *store.UpdateSpotInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateSpotParams{
		Name:         in.Name,
		Description:  in.Description,
		ThumbnailURL: in.ThumbnailURL,
		Longitude:    in.Longitude,
		Latitude:     in.Latitude,
	}
	err := s.db.Spot.Update(ctx, in.SpotID, params)
	return internalError(err)
}

func (s *service) DeleteSpot(ctx context.Context, in *store.DeleteSpotInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Spot.Delete(ctx, in.SpotID)
	return internalError(err)
}

func (s *service) ApproveSpot(ctx context.Context, in *store.ApproveSpotInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.ApproveSpotParams{
		Approved:        in.Approved,
		ApprovedAdminID: in.AdminID,
	}
	err := s.db.Spot.Approve(ctx, in.SpotID, params)
	return internalError(err)
}
