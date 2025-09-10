package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) topRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/top", h.authentication)

	r.GET("/orders", h.TopOrders)
}

// @Summary     注文集計取得
// @Description 指定期間の注文統計情報（注文数、売上、決済手段別集計、推移）を取得します。
// @Tags        Top
// @Router      /v1/top/orders [get]
// @Security    bearerauth
// @Param       startAt query integer false "集計開始日時（unixtime,未指定の場合は１週間前の時刻）" example("1640962800")
// @Param       endAt query integer false "集計終了日時（unixtime,未指定の場合は現在時刻）" example("1640962800")
// @Param       periodType query string false "集計期間（未指定の場合は日次）" example("day")
// @Param       shopId query string false "店舗ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.TopOrdersResponse
func (h *handler) TopOrders(ctx *gin.Context) {
	const defaultPeriodType = types.TopOrderPeriodTypeDay

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
	shopID := util.GetQuery(ctx, "shopId", "")
	if getAdminType(ctx).Response() == types.AdminTypeCoordinator {
		shopID = getShopID(ctx)
	}
	periodTypeStr := util.GetQuery(ctx, "periodType", string(defaultPeriodType))

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
			ShopID:       shopID,
			CreatedAtGte: startAt,
			CreatedAtLt:  endAt,
		}
		current, err = h.store.AggregateOrders(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		diff := endAt.Sub(startAt)
		in := &store.AggregateOrdersInput{
			ShopID:       shopID,
			CreatedAtGte: startAt.Add(-diff),
			CreatedAtLt:  endAt.Add(-diff),
		}
		previous, err = h.store.AggregateOrders(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.AggregateOrdersByPeriodInput{
			ShopID:       shopID,
			PeriodType:   periodType.StoreEntity(),
			CreatedAtGte: startAt,
			CreatedAtLt:  endAt,
		}
		orders, err = h.store.AggregateOrdersByPeriod(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &store.AggregateOrdersByPaymentMethodTypeInput{
			ShopID:       shopID,
			CreatedAtGte: startAt,
			CreatedAtLt:  endAt,
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

	res := &types.TopOrdersResponse{
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
