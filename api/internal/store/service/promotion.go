package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) ListPromotions(
	ctx context.Context, in *store.ListPromotionsInput,
) (entity.Promotions, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	// TODO: 詳細の実装
	var (
		promotions entity.Promotions
		total      int64
	)
	return promotions, total, nil
}

func (s *service) GetPromotion(ctx context.Context, in *store.GetPromotionInput) (*entity.Promotion, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	// TODO: 詳細の実装
	promotion := &entity.Promotion{}
	return promotion, nil
}

func (s *service) CreatePromotion(ctx context.Context, in *store.CreatePromotionInput) (*entity.Promotion, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	// TODO: 詳細の実装
	promotion := &entity.Promotion{}
	return promotion, nil
}

func (s *service) UpdatePromotion(ctx context.Context, in *store.UpdatePromotionInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	// TODO: 詳細の実装
	return nil
}

func (s *service) DeletePromotion(ctx context.Context, in *store.DeletePromotionInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	// TODO: 詳細の実装
	return nil
}
