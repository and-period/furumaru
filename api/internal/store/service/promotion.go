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

func (s *service) ListPromotions(
	ctx context.Context, in *store.ListPromotionsInput,
) (entity.Promotions, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	orders := make([]*database.ListPromotionsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListPromotionsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListPromotionsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
		Orders: orders,
	}
	var (
		promotions entity.Promotions
		total      int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		promotions, err = s.db.Promotion.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Promotion.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return promotions, total, nil
}

func (s *service) GetPromotion(ctx context.Context, in *store.GetPromotionInput) (*entity.Promotion, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	promotion, err := s.db.Promotion.Get(ctx, in.PromotionID)
	return promotion, exception.InternalError(err)
}

func (s *service) CreatePromotion(ctx context.Context, in *store.CreatePromotionInput) (*entity.Promotion, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &entity.NewPromotionParams{
		Title:        in.Title,
		Description:  in.Description,
		Public:       in.Public,
		PublishedAt:  in.PublishedAt,
		DiscountType: in.DiscountType,
		DiscountRate: in.DiscountRate,
		Code:         in.Code,
		CodeType:     in.CodeType,
		StartAt:      in.StartAt,
		EndAt:        in.EndAt,
	}
	promotion := entity.NewPromotion(params)
	if err := promotion.Validate(); err != nil {
		return nil, fmt.Errorf("api: validation error: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err := s.db.Promotion.Create(ctx, promotion); err != nil {
		return nil, exception.InternalError(err)
	}
	return promotion, nil
}

func (s *service) UpdatePromotion(ctx context.Context, in *store.UpdatePromotionInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdatePromotionParams{
		Title:        in.Title,
		Description:  in.Description,
		Public:       in.Public,
		PublishedAt:  in.PublishedAt,
		DiscountType: in.DiscountType,
		DiscountRate: in.DiscountRate,
		Code:         in.Code,
		CodeType:     in.CodeType,
		StartAt:      in.StartAt,
		EndAt:        in.EndAt,
	}
	err := s.db.Promotion.Update(ctx, in.PromotionID, params)
	return exception.InternalError(err)
}

func (s *service) DeletePromotion(ctx context.Context, in *store.DeletePromotionInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	// TODO: 詳細の実装
	return nil
}
