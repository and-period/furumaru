package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) multiGetPromotion(ctx context.Context, promotionIDs []string) (service.Promotions, error) {
	if len(promotionIDs) == 0 {
		return service.Promotions{}, nil
	}
	in := &store.MultiGetPromotionsInput{
		PromotionIDs: promotionIDs,
	}
	promotions, err := h.store.MultiGetPromotions(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewPromotions(promotions), nil
}

func (h *handler) getPromotion(ctx context.Context, promotionID string) (*service.Promotion, error) {
	in := &store.GetPromotionInput{
		PromotionID: promotionID,
	}
	promotion, err := h.store.GetPromotion(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewPromotion(promotion), nil
}

func (h *handler) getEnabledPromotion(ctx context.Context, code string) (*service.Promotion, error) {
	in := &store.GetPromotionByCodeInput{
		PromotionCode: code,
		OnlyEnabled:   true,
	}
	promotion, err := h.store.GetPromotionByCode(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewPromotion(promotion), nil
}
