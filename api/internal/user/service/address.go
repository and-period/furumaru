package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListAddresses(
	ctx context.Context,
	in *user.ListAddressesInput,
) (entity.Addresses, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListAddressesParams{
		UserID: in.UserID,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		addresses entity.Addresses
		total     int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		addresses, err = s.db.Address.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Address.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return addresses, total, nil
}

func (s *service) ListDefaultAddresses(
	ctx context.Context,
	in *user.ListDefaultAddressesInput,
) (entity.Addresses, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	addresses, err := s.db.Address.ListDefault(ctx, in.UserIDs)
	return addresses, internalError(err)
}

func (s *service) MultiGetAddresses(
	ctx context.Context,
	in *user.MultiGetAddressesInput,
) (entity.Addresses, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	addresses, err := s.db.Address.MultiGet(ctx, in.AddressIDs)
	return addresses, internalError(err)
}

func (s *service) MultiGetAddressesByRevision(
	ctx context.Context,
	in *user.MultiGetAddressesByRevisionInput,
) (entity.Addresses, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	addresses, err := s.db.Address.MultiGetByRevision(ctx, in.AddressRevisionIDs)
	return addresses, internalError(err)
}

func (s *service) GetAddress(
	ctx context.Context,
	in *user.GetAddressInput,
) (*entity.Address, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	address, err := s.db.Address.Get(ctx, in.AddressID)
	if err != nil {
		return nil, internalError(err)
	}
	if in.UserID != address.UserID {
		return nil, fmt.Errorf(
			"service: this address belongs to another user: %w",
			exception.ErrForbidden,
		)
	}
	return address, nil
}

func (s *service) GetDefaultAddress(
	ctx context.Context,
	in *user.GetDefaultAddressInput,
) (*entity.Address, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	address, err := s.db.Address.GetDefault(ctx, in.UserID)
	return address, internalError(err)
}

func (s *service) CreateAddress(
	ctx context.Context,
	in *user.CreateAddressInput,
) (*entity.Address, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewAddressParams{
		UserID:         in.UserID,
		IsDefault:      in.IsDefault,
		Lastname:       in.Lastname,
		Firstname:      in.Firstname,
		LastnameKana:   in.LastnameKana,
		FirstnameKana:  in.FirstnameKana,
		PostalCode:     in.PostalCode,
		PrefectureCode: in.PrefectureCode,
		City:           in.City,
		AddressLine1:   in.AddressLine1,
		AddressLine2:   in.AddressLine2,
		PhoneNumber:    in.PhoneNumber,
	}
	address, err := entity.NewAddress(params)
	if err != nil {
		return nil, fmt.Errorf(
			"service: failed to new address: %w: %s",
			exception.ErrInvalidArgument,
			err.Error(),
		)
	}
	if err := s.db.Address.Create(ctx, address); err != nil {
		return nil, internalError(err)
	}
	return address, nil
}

func (s *service) UpdateAddress(ctx context.Context, in *user.UpdateAddressInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if _, err := codes.ToPrefectureJapanese(in.PrefectureCode); err != nil {
		return fmt.Errorf(
			"service: invalid prefecture code: %w: %s",
			exception.ErrInvalidArgument,
			err.Error(),
		)
	}
	address, err := s.db.Address.Get(ctx, in.AddressID, "user_id")
	if err != nil {
		return internalError(err)
	}
	if in.UserID != address.UserID {
		return fmt.Errorf(
			"service: this address belongs to another user: %w",
			exception.ErrForbidden,
		)
	}
	params := &database.UpdateAddressParams{
		Lastname:       in.Lastname,
		Firstname:      in.Firstname,
		LastnameKana:   in.LastnameKana,
		FirstnameKana:  in.FirstnameKana,
		PostalCode:     in.PostalCode,
		PrefectureCode: in.PrefectureCode,
		City:           in.City,
		AddressLine1:   in.AddressLine1,
		AddressLine2:   in.AddressLine2,
		PhoneNumber:    in.PhoneNumber,
		IsDefault:      in.IsDefault,
	}
	err = s.db.Address.Update(ctx, in.AddressID, in.UserID, params)
	return internalError(err)
}

func (s *service) DeleteAddress(ctx context.Context, in *user.DeleteAddressInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Address.Delete(ctx, in.AddressID, in.UserID)
	return internalError(err)
}
