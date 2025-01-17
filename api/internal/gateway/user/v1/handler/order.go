package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (h *handler) orderRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/orders", h.authentication)

	r.GET("", h.ListOrders)
	r.GET("/:orderId", h.GetOrder)
}

func (h *handler) ListOrders(ctx *gin.Context) {
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
	types, err := util.GetQueryInt32s(ctx, "type")
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	orderTypes := make([]sentity.OrderType, len(types))
	for i, t := range types {
		orderTypes[i] = sentity.OrderType(t)
	}
	orderStatuses := []sentity.OrderStatus{
		sentity.OrderStatusUnpaid,
		sentity.OrderStatusWaiting,
		sentity.OrderStatusPreparing,
		sentity.OrderStatusShipped,
		sentity.OrderStatusCompleted,
		sentity.OrderStatusCanceled,
		sentity.OrderStatusRefunded,
	}
	ordersIn := &store.ListOrdersInput{
		UserID:   h.getUserID(ctx),
		Limit:    limit,
		Offset:   offset,
		Types:    orderTypes,
		Statuses: orderStatuses,
	}
	orders, total, err := h.store.ListOrders(ctx, ordersIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(orders) == 0 {
		res := &response.OrdersResponse{
			Order:        []*response.Order{},
			Coordinators: []*response.Coordinator{},
			Promotions:   []*response.Promotion{},
			Products:     []*response.Product{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		addresses    service.Addresses
		coordinators service.Coordinators
		promotions   service.Promotions
		oproducts    service.Products
		cproducts    service.Products
		experiences  service.Experiences
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		addresses, err = h.multiGetAddressesByRevision(ectx, orders.AddressRevisionIDs())
		return
	})
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinatorsWithDeleted(ectx, orders.CoordinatorIDs())
		return
	})
	eg.Go(func() (err error) {
		promotions, err = h.multiGetPromotion(ectx, orders.PromotionIDs())
		return
	})
	eg.Go(func() (err error) {
		oproducts, err = h.multiGetProductsByRevision(ectx, orders.ProductRevisionIDs())
		if err != nil {
			return
		}
		cproducts, err = h.multiGetProducts(ectx, oproducts.IDs())
		return
	})
	eg.Go(func() (err error) {
		experiences, err = h.multiGetExperiencesByRevision(ectx, orders.ExperienceRevisionIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.OrdersResponse{
		Order:        service.NewOrders(orders, addresses.MapByRevision(), oproducts.MapByRevision(), experiences.MapByRevision()).Response(),
		Coordinators: coordinators.Response(),
		Promotions:   promotions.Response(),
		Products:     cproducts.Response(),
		Experiences:  experiences.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetOrder(ctx *gin.Context) {
	order, err := h.getOrder(ctx, h.getUserID(ctx), util.GetParam(ctx, "orderId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		coordinator *service.Coordinator
		promotion   *service.Promotion
		products    service.Products
		experience  *service.Experience
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinatorWithDeleted(ectx, order.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		if order.PromotionID == "" {
			return
		}
		promotion, err = h.getPromotion(ectx, order.PromotionID)
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProducts(ectx, order.ProductIDs())
		return
	})
	eg.Go(func() (err error) {
		if order.Experience.ExperienceID == "" {
			return
		}
		experience, err = h.getExperience(ectx, order.Experience.ExperienceID)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.OrderResponse{
		Order:       order.Response(),
		Coordinator: coordinator.Response(),
		Promotion:   promotion.Response(),
		Products:    products.Response(),
		Experience:  experience.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getOrder(ctx context.Context, userID, orderID string) (*service.Order, error) {
	in := &store.GetOrderInput{
		OrderID: orderID,
	}
	order, err := h.store.GetOrder(ctx, in)
	if err != nil {
		return nil, err
	}
	if userID != order.UserID {
		// 不正の疑いがあるため、リクエスト情報をログ出力しておく
		h.logger.Warn("UserId does not match order information", zap.String("userId", userID), zap.String("orderId", orderID))
		return nil, fmt.Errorf("%s: %w", errNotFoundOrder, exception.ErrNotFound)
	}
	var (
		addresses   service.Addresses
		products    service.Products
		experiences service.Experiences
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		addresses, err = h.multiGetAddressesByRevision(ectx, order.AddressRevisionIDs())
		return
	})
	eg.Go(func() (err error) {
		products, err = h.multiGetProductsByRevision(ectx, order.ProductRevisionIDs())
		return
	})
	eg.Go(func() (err error) {
		experiences, err = h.multiGetExperiencesByRevision(ectx, []int64{order.OrderExperience.ExperienceRevisionID})
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return service.NewOrder(order, addresses.MapByRevision(), products.MapByRevision(), experiences.MapByRevision()), nil
}
