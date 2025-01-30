package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) topRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/top")

	r.GET("/orders", h.TopOrders)
}

func (h *handler) TopOrders(ctx *gin.Context) {
	const defaultPeriodType = service.TopOrderPeriodTypeDay

	defaultEndAt := h.now()
	defaultStartAt := defaultEndAt.AddDate(0, 0, -7)

	startAtUnix, err := util.GetQueryInt64(ctx, "startAt", defaultStartAt.Unix())
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	endAtUnix, err := util.GetQueryInt64(ctx, "endAt", defaultEndAt.Unix())
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	coordinatorID := util.GetQuery(ctx, "coordinatorId", "")
	if getAdminType(ctx) == service.AdminTypeCoordinator {
		coordinatorID = getAdminID(ctx)
	}
	periodTypeStr := util.GetQuery(ctx, "periodType", defaultPeriodType.Response())

	startAt := jst.ParseFromUnix(startAtUnix)
	endAt := jst.ParseFromUnix(endAtUnix)
	periodType := service.NewTopOrderPeriodTypeFromRequest(periodTypeStr)

	var (
		current  *sentity.AggregatedOrder
		previous *sentity.AggregatedOrder
		orders   sentity.AggregatedPeriodOrders
		payments sentity.AggregatedOrderPayments
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &store.AggregateOrdersInput{
			CoordinatorID: coordinatorID,
			CreatedAtGte:  startAt,
			CreatedAtLt:   endAt,
		}
		current, err = h.store.AggregateOrders(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		diff := endAt.Sub(startAt)
		in := &store.AggregateOrdersInput{
			CoordinatorID: coordinatorID,
			CreatedAtGte:  startAt.Add(-diff),
			CreatedAtLt:   endAt.Add(-diff),
		}
		previous, err = h.store.AggregateOrders(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.AggregateOrdersByPeriodInput{
			CoordinatorID: coordinatorID,
			PeriodType:    periodType.StoreEntity(),
			CreatedAtGte:  startAt,
			CreatedAtLt:   endAt,
		}
		orders, err = h.store.AggregateOrdersByPeriod(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.AggregateOrdersByPaymentMethodTypeInput{
			CoordinatorID: coordinatorID,
			CreatedAtGte:  startAt,
			CreatedAtLt:   endAt,
		}
		payments, err = h.store.AggregateOrdersByPaymentMethodType(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}
	if current == nil {
		current = &sentity.AggregatedOrder{}
	}
	if previous == nil {
		previous = &sentity.AggregatedOrder{}
	}

	res := &response.TopOrdersResponse{
		StartAt:     jst.Unix(startAt),
		EndAt:       jst.Unix(endAt),
		PeriodType:  periodType.Response(),
		Orders:      service.NewTopOrderValue(current.OrderCount, previous.OrderCount).Response(),
		Users:       service.NewTopOrderValue(current.UserCount, previous.UserCount).Response(),
		Sales:       service.NewTopOrderValue(current.SalesTotal, previous.SalesTotal).Response(),
		Payments:    service.NewTopOrderPayments(payments).Response(),
		SalesTrends: service.NewTopOrderSalesTrends(periodType, startAt, endAt, orders).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
