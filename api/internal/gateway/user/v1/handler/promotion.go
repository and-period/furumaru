package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) getEnabledPromotion(ctx context.Context, promotionID string) (*service.Promotion, error) {
	in := &store.GetPromotionInput{
		PromotionID: promotionID,
		OnlyEnabled: true,
	}
	promotion, err := h.store.GetPromotion(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewPromotion(promotion), nil
}
