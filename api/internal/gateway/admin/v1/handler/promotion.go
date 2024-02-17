package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
)

func (h *handler) promotionRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/promotions", h.authentication)

	r.GET("", h.ListPromotions)
	r.POST("", h.CreatePromotion)
	r.GET("/:promotionId", h.GetPromotion)
	r.PATCH("/:promotionId", h.UpdatePromotion)
	r.DELETE("/:promotionId", h.DeletePromotion)
}

func (h *handler) ListPromotions(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	orders, err := h.newPromotionOrders(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListPromotionsInput{
		Title:  util.GetQuery(ctx, "title", ""),
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}
	promotions, total, err := h.store.ListPromotions(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	aggregates, err := h.aggregateOrdersByPromotion(ctx, promotions.IDs()...)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.PromotionsResponse{
		Promotions: service.NewPromotions(promotions, aggregates).Response(),
		Total:      total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newPromotionOrders(ctx *gin.Context) ([]*store.ListPromotionsOrder, error) {
	categories := map[string]sentity.PromotionOrderBy{
		"title":       sentity.PromotionOrderByTitle,
		"public":      sentity.PromotionOrderByPublic,
		"publishedAt": sentity.PromotionOrderByPublishedAt,
		"startAt":     sentity.PromotionOrderByStartAt,
		"endAt":       sentity.PromotionOrderByEndAt,
		"createdAt":   sentity.PromotionOrderByCreatedAt,
		"updatedAt":   sentity.PromotionOrderByUpdatedAt,
	}
	params := util.GetOrders(ctx)
	if len(params) == 0 {
		res := []*store.ListPromotionsOrder{
			{Key: sentity.PromotionOrderByPublishedAt, OrderByASC: false},
			{Key: sentity.PromotionOrderByPublic, OrderByASC: false},
		}
		return res, nil
	}
	res := make([]*store.ListPromotionsOrder, len(params))
	for i, p := range params {
		key, ok := categories[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderKey)
		}
		res[i] = &store.ListPromotionsOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) GetPromotion(ctx *gin.Context) {
	promotion, err := h.getPromotion(ctx, util.GetParam(ctx, "promotionId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.PromotionResponse{
		Promotion: promotion.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreatePromotion(ctx *gin.Context) {
	req := &request.CreatePromotionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.CreatePromotionInput{
		Title:        req.Title,
		Description:  req.Description,
		Public:       req.Public,
		DiscountType: service.DiscountType(req.DiscountType).StoreEntity(),
		DiscountRate: req.DiscountRate,
		Code:         req.Code,
		CodeType:     sentity.PromotionCodeTypeAlways, // 回数無制限固定
		StartAt:      jst.ParseFromUnix(req.StartAt),
		EndAt:        jst.ParseFromUnix(req.EndAt),
	}
	promotion, err := h.store.CreatePromotion(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.PromotionResponse{
		// 初回は集計結果が存在しないためnilで渡す
		Promotion: service.NewPromotion(promotion, nil).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdatePromotion(ctx *gin.Context) {
	req := &request.UpdatePromotionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.UpdatePromotionInput{
		PromotionID:  util.GetParam(ctx, "promotionId"),
		Title:        req.Title,
		Description:  req.Description,
		Public:       req.Public,
		DiscountType: service.DiscountType(req.DiscountType).StoreEntity(),
		DiscountRate: req.DiscountRate,
		Code:         req.Code,
		CodeType:     sentity.PromotionCodeTypeAlways, // 回数無制限固定
		StartAt:      jst.ParseFromUnix(req.StartAt),
		EndAt:        jst.ParseFromUnix(req.EndAt),
	}
	if err := h.store.UpdatePromotion(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeletePromotion(ctx *gin.Context) {
	in := &store.DeletePromotionInput{
		PromotionID: util.GetParam(ctx, "promotionId"),
	}
	if err := h.store.DeletePromotion(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) multiGetPromotions(ctx context.Context, promotionIDs []string) (service.Promotions, error) {
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
	aggregates, err := h.aggregateOrdersByPromotion(ctx, promotionIDs...)
	if err != nil {
		return nil, err
	}
	return service.NewPromotions(promotions, aggregates), nil
}

func (h *handler) getPromotion(ctx context.Context, promotionID string) (*service.Promotion, error) {
	in := &store.GetPromotionInput{
		PromotionID: promotionID,
	}
	promotion, err := h.store.GetPromotion(ctx, in)
	if err != nil {
		return nil, err
	}
	aggregates, err := h.aggregateOrdersByPromotion(ctx, promotionID)
	if err != nil {
		return nil, err
	}
	return service.NewPromotion(promotion, aggregates[promotionID]), nil
}

func (h *handler) aggregateOrdersByPromotion(
	ctx context.Context,
	promotionIDs ...string,
) (map[string]*sentity.AggregatedOrderPromotion, error) {
	if len(promotionIDs) == 0 {
		return map[string]*sentity.AggregatedOrderPromotion{}, nil
	}
	in := &store.AggregateOrdersByPromotionInput{
		PromotionIDs: promotionIDs,
	}
	aggregates, err := h.store.AggregateOrdersByPromotion(ctx, in)
	if err != nil {
		return nil, err
	}
	return aggregates.Map(), nil
}
