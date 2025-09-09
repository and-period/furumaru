package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

// @tag.name        Promotion
// @tag.description プロモーション関連
func (h *handler) promotionRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/promotions")

	r.GET("/:code", h.GetPromotion)
}

// @Summary     プロモーション詳細取得
// @Description プロモーションコードから割引情報を取得します。
// @Tags        Promotion
// @Router      /promotions/{code} [get]
// @Param       code path string true "プロモーションコード"
// @Param       coordinatorId query string false "コーディネーターID"
// @Produce     json
// @Success     200 {object} types.PromotionResponse
// @Failure     403 {object} util.ErrorResponse "プロモーションが利用できません"
// @Failure     404 {object} util.ErrorResponse "プロモーションが見つかりません"
func (h *handler) GetPromotion(ctx *gin.Context) {
	in := &store.GetPromotionByCodeInput{
		PromotionCode: util.GetParam(ctx, "code"),
	}
	promotion, err := h.store.GetPromotionByCode(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if !promotion.IsEnabled(promotion.ShopID) {
		h.forbidden(ctx, errors.New("handler: this promotion is disabled"))
		return
	}
	res := &types.PromotionResponse{
		Promotion: service.NewPromotion(promotion).Response(),
	}
	if promotion.ShopID == "" {
		ctx.JSON(http.StatusOK, res)
		return
	}
	shop, err := h.getShop(ctx, promotion.ShopID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if shop.CoordinatorID != util.GetQuery(ctx, "coordinatorId", "") {
		h.forbidden(ctx, errors.New("handler: this promotion can only be used at certain shop"))
		return
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
