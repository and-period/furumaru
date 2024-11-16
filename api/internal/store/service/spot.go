package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) ListSpots(ctx context.Context, in *store.ListSpotsInput) (entity.Spots, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	// TODO: 詳細の実装
	return entity.Spots{}, 0, nil
}

func (s *service) ListSpotsByGeolocation(ctx context.Context, in *store.ListSpotsByGeolocationInput) (entity.Spots, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return entity.Spots{}, nil
}

func (s *service) GetSpot(ctx context.Context, in *store.GetSpotInput) (*entity.Spot, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.Spot{}, nil
}

func (s *service) CreateSpotByUser(ctx context.Context, in *store.CreateSpotByUserInput) (*entity.Spot, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.Spot{}, nil
}

func (s *service) CreateSpotByAdmin(ctx context.Context, in *store.CreateSpotByAdminInput) (*entity.Spot, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.Spot{}, nil
}

func (s *service) UpdateSpot(ctx context.Context, in *store.UpdateSpotInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}

func (s *service) DeleteSpot(ctx context.Context, in *store.DeleteSpotInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}

func (s *service) ApproveSpot(ctx context.Context, in *store.ApproveSpotInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}
