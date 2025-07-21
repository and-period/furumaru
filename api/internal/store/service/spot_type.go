package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListSpotTypes(
	ctx context.Context,
	in *store.ListSpotTypesInput,
) (entity.SpotTypes, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListSpotTypesParams{
		Name:   in.Name,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		types entity.SpotTypes
		total int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		types, err = s.db.SpotType.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.SpotType.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return types, total, nil
}

func (s *service) MultiGetSpotTypes(
	ctx context.Context,
	in *store.MultiGetSpotTypesInput,
) (entity.SpotTypes, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	types, err := s.db.SpotType.MultiGet(ctx, in.SpotTypeIDs)
	return types, internalError(err)
}

func (s *service) GetSpotType(
	ctx context.Context,
	in *store.GetSpotTypeInput,
) (*entity.SpotType, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	typ, err := s.db.SpotType.Get(ctx, in.SpotTypeID)
	return typ, internalError(err)
}

func (s *service) CreateSpotType(
	ctx context.Context,
	in *store.CreateSpotTypeInput,
) (*entity.SpotType, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewSpotTypeParams{
		Name: in.Name,
	}
	typ := entity.NewSpotType(params)
	if err := s.db.SpotType.Create(ctx, typ); err != nil {
		return nil, internalError(err)
	}
	return typ, nil
}

func (s *service) UpdateSpotType(ctx context.Context, in *store.UpdateSpotTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateSpotTypeParams{
		Name: in.Name,
	}
	err := s.db.SpotType.Update(ctx, in.SpotTypeID, params)
	return internalError(err)
}

func (s *service) DeleteSpotType(ctx context.Context, in *store.DeleteSpotTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.SpotType.Delete(ctx, in.SpotTypeID)
	return internalError(err)
}
