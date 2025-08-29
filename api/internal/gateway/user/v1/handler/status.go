package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        Status
// @tag.description ステータス関連
func (h *handler) statusRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/statuses")

	r.GET("/payments", h.ListPaymentStatuses)
}

// @Summary     支払手段一覧取得
// @Description 利用可能な支払手段の一覧を取得します。
// @Tags        Status
// @Router      /statuses/payments [get]
// @Produce     json
// @Success     200 {object} response.PaymentSystemsResponse
func (h *handler) ListPaymentStatuses(ctx *gin.Context) {
	methodTypes := []entity.PaymentMethodType{
		entity.PaymentMethodTypeCash,
		entity.PaymentMethodTypeCreditCard,
		entity.PaymentMethodTypeKonbini,
		entity.PaymentMethodTypeBankTransfer,
		entity.PaymentMethodTypePayPay,
		entity.PaymentMethodTypeLinePay,
		entity.PaymentMethodTypeMerpay,
		entity.PaymentMethodTypeRakutenPay,
		entity.PaymentMethodTypeAUPay,
		entity.PaymentMethodTypePaidy,
		entity.PaymentMethodTypePayEasy,
	}
	in := &store.MultiGetPaymentSystemsInput{
		MethodTypes: methodTypes,
	}
	systems, err := h.store.MultiGetPaymentSystems(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.PaymentSystemsResponse{
		Systems: service.NewPaymentSystems(systems).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getPaymentSystem(ctx context.Context, methodType service.PaymentMethodType) (*service.PaymentSystem, error) {
	in := &store.GetPaymentSystemInput{
		MethodType: methodType.StoreEntity(),
	}
	system, err := h.store.GetPaymentSystem(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewPaymentSystem(system), nil
}
