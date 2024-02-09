package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *handler) GetPromotion(ctx *gin.Context) {
	promotion, err := h.getPromotion(ctx, util.GetParam(ctx, "code"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if !promotion.Enabled() {
		h.forbidden(ctx, errors.New("handler: this promotion is disabled"))
		return
	}
	res := &response.PromotionResponse{
		Promotion: promotion.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

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

//nolint:unused
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
