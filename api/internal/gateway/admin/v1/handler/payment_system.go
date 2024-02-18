package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) paymentSystemRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/payment-systems", h.authentication)

	r.GET("", h.ListPaymentSystems)
	r.PATCH("/:methodType", h.UpdatePaymentSystem)
}

func (h *handler) ListPaymentSystems(ctx *gin.Context) {
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

func (h *handler) UpdatePaymentSystem(ctx *gin.Context) {
	req := &request.UpdatePaymentSystemRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	methodType, err := util.GetParamInt64(ctx, "methodType")
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &store.UpdatePaymentStatusInput{
		MethodType: service.PaymentMethodType(methodType).StoreEntity(),
		Status:     service.PaymentSystemStatus(req.Status).StoreEntity(),
	}
	if err := h.store.UpdatePaymentSystem(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
