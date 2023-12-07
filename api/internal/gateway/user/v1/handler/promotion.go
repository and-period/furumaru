package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
)

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
