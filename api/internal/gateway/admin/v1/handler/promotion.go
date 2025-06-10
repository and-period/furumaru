package handler

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) promotionRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/promotions", h.authentication)

	r.GET("", h.ListPromotions)
	r.POST("", h.CreatePromotion)
	r.GET("/:promotionId", h.filterAccessPromotion, h.GetPromotion)
	r.PATCH("/:promotionId", h.filterAccessPromotion, h.UpdatePromotion)
	r.DELETE("/:promotionId", h.filterAccessPromotion, h.DeletePromotion)
}

func (h *handler) filterAccessPromotion(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			promotion, err := h.getPromotion(ctx, util.GetParam(ctx, "promotionId"))
			if err != nil {
				return false, err
			}
			if service.PromotionTargetType(promotion.TargetType) == service.PromotionTargetTypeAllShop {
				return true, nil
			}
			shop, err := h.getShop(ctx, getShopID(ctx))
			if err != nil {
				return false, err
			}
			return promotion.ShopID == shop.ID, nil
		},
		producer: func(ctx *gin.Context) (bool, error) {
			promotion, err := h.getPromotion(ctx, util.GetParam(ctx, "promotionId"))
			if err != nil {
				return false, err
			}
			if service.PromotionTargetType(promotion.TargetType) == service.PromotionTargetTypeAllShop {
				return true, nil
			}
			shop, err := h.getShop(ctx, promotion.ID)
			if err != nil {
				return false, err
			}
			return slices.Contains(shop.ProducerIDs, getAdminID(ctx)), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Next()
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
	shopID := util.GetQuery(ctx, "shopId", "")

	in := &store.ListPromotionsInput{
		ShopID: shopID,
		Title:  util.GetQuery(ctx, "title", ""),
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}
	if getAdminType(ctx) == service.AdminTypeCoordinator {
		withAllTarget, err := util.GetQueryBool(ctx, "withAllTarget", true)
		if err != nil {
			h.badRequest(ctx, err)
			return
		}
		shop, err := h.getShop(ctx, getShopID(ctx))
		if err != nil {
			h.httpError(ctx, err)
			return
		}

		in.ShopID = shop.ID
		in.WithAllTarget = withAllTarget
	}

	promotions, total, err := h.store.ListPromotions(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(promotions) == 0 {
		res := &response.PromotionsResponse{
			Promotions: []*response.Promotion{},
			Shops:      []*response.Shop{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		aggregates map[string]*sentity.AggregatedOrderPromotion
		shops      service.Shops
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		aggregates, err = h.aggregateOrdersByPromotion(ectx, promotions.IDs()...)
		return
	})
	eg.Go(func() (err error) {
		shops, err = h.multiGetShops(ectx, promotions.ShopIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.PromotionsResponse{
		Promotions: service.NewPromotions(promotions, aggregates).Response(),
		Shops:      shops.Response(),
		Total:      total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newPromotionOrders(ctx *gin.Context) ([]*store.ListPromotionsOrder, error) {
	categories := map[string]store.ListPromotionsOrderKey{
		"title":     store.ListPromotionsOrderByTitle,
		"public":    store.ListPromotionsOrderByPublic,
		"startAt":   store.ListPromotionsOrderByStartAt,
		"endAt":     store.ListPromotionsOrderByEndAt,
		"createdAt": store.ListPromotionsOrderByCreatedAt,
		"updatedAt": store.ListPromotionsOrderByUpdatedAt,
	}
	params := util.GetOrders(ctx)
	if len(params) == 0 {
		res := []*store.ListPromotionsOrder{
			{Key: store.ListPromotionsOrderByStartAt, OrderByASC: false},
			{Key: store.ListPromotionsOrderByPublic, OrderByASC: false},
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
	if promotion.ShopID == "" {
		ctx.JSON(http.StatusOK, res)
		return
	}

	shop, err := h.getShop(ctx, promotion.ShopID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res.Shop = shop.Response()
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreatePromotion(ctx *gin.Context) {
	req := &request.CreatePromotionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.CreatePromotionInput{
		AdminID:      getAdminID(ctx),
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
	if promotion.ShopID == "" {
		ctx.JSON(http.StatusOK, res)
		return
	}

	shop, err := h.getShop(ctx, promotion.ShopID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res.Shop = shop.Response()
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
		AdminID:      getAdminID(ctx),
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
