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
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/geolocation"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListSpots(ctx context.Context, in *store.ListSpotsInput) (entity.Spots, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListSpotsParams{
		Name:            in.Name,
		SpotTypeIDs:     in.TypeIDs,
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
		SpotTypeIDs:     in.TypeIDs,
		Longitude:       in.Longitude,
		Latitude:        in.Latitude,
		Radius:          in.Radius,
		ExcludeDisabled: in.ExcludeDisabled,
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
	_, err := s.db.SpotType.Get(ctx, in.TypeID)
	if errors.Is(err, database.ErrNotFound) {
		return nil, fmt.Errorf("service: this spot type is not found: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	addressIn := &geolocation.GetAddressInput{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
	address, err := s.geolocation.GetAddress(ctx, addressIn)
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.SpotParams{
		UserID:       in.UserID,
		SpotTypeID:   in.TypeID,
		Name:         in.Name,
		Description:  in.Description,
		ThumbnailURL: in.ThumbnailURL,
		Longitude:    in.Longitude,
		Latitude:     in.Latitude,
		PostalCode:   address.PostalCode,
		Prefecture:   address.Prefecture,
		City:         address.City,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
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
	_, err := s.db.SpotType.Get(ctx, in.TypeID)
	if errors.Is(err, database.ErrNotFound) {
		return nil, fmt.Errorf("service: this spot type is not found: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	adminIn := &user.GetAdminInput{
		AdminID: in.AdminID,
	}
	admin, err := s.user.GetAdmin(ctx, adminIn)
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("service: this admin is not found: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	var userType entity.SpotUserType
	switch admin.Type {
	case uentity.AdminTypeCoordinator:
		userType = entity.SpotUserTypeCoordinator
	case uentity.AdminTypeProducer:
		userType = entity.SpotUserTypeProducer
	default:
		return nil, fmt.Errorf("service: unsupported admin type: %w", exception.ErrFailedPrecondition)
	}
	addressIn := &geolocation.GetAddressInput{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
	address, err := s.geolocation.GetAddress(ctx, addressIn)
	if err != nil {
		return nil, internalError(err)
	}

	params := &entity.SpotParams{
		UserType:     userType,
		UserID:       in.AdminID,
		SpotTypeID:   in.TypeID,
		Name:         in.Name,
		Description:  in.Description,
		ThumbnailURL: in.ThumbnailURL,
		Longitude:    in.Longitude,
		Latitude:     in.Latitude,
		PostalCode:   address.PostalCode,
		Prefecture:   address.Prefecture,
		City:         address.City,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
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
	_, err := s.db.SpotType.Get(ctx, in.TypeID)
	if errors.Is(err, database.ErrNotFound) {
		return fmt.Errorf("service: this spot type is not found: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}
	addressIn := &geolocation.GetAddressInput{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
	address, err := s.geolocation.GetAddress(ctx, addressIn)
	if err != nil {
		return internalError(err)
	}
	prefectureCode, err := codes.ToPrefectureValue(address.Prefecture)
	if err != nil {
		return fmt.Errorf("service: invalid prefecture: %w", exception.ErrInvalidArgument)
	}
	params := &database.UpdateSpotParams{
		SpotTypeID:     in.TypeID,
		Name:           in.Name,
		Description:    in.Description,
		ThumbnailURL:   in.ThumbnailURL,
		Longitude:      in.Longitude,
		Latitude:       in.Latitude,
		PostalCode:     address.PostalCode,
		PrefectureCode: prefectureCode,
		City:           address.City,
		AddressLine1:   address.AddressLine1,
		AddressLine2:   address.AddressLine2,
	}
	err = s.db.Spot.Update(ctx, in.SpotID, params)
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
