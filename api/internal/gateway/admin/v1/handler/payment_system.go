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

// @tag.name        PaymentSystem
// @tag.description 決済システム関連
func (h *handler) paymentSystemRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/payment-systems", h.authentication)

	r.GET("", h.ListPaymentSystems)
	r.PATCH("/:methodType", h.UpdatePaymentSystem)
}

// @Summary     決済システム一覧取得
// @Description 決済手段毎のシステム状態一覧を取得します。
// @Tags        PaymentSystem
// @Router      /v1/payment-systems [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} response.PaymentSystemsResponse
func (h *handler) ListPaymentSystems(ctx *gin.Context) {
	methodTypes := []entity.PaymentMethodType{
		entity.PaymentMethodTypeCreditCard,
		entity.PaymentMethodTypeBankTransfer,
		entity.PaymentMethodTypePayPay,
		entity.PaymentMethodTypeMerpay,
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

// @Summary     決済システム更新
// @Description 指定された決済手段のシステム状態を更新します。
// @Tags        PaymentSystem
// @Router      /v1/payment-systems/{methodType} [patch]
// @Security    bearerauth
// @Param       methodType path integer true "決済手段タイプ" example(1)
// @Accept      json
// @Param       request body request.UpdatePaymentSystemRequest true "決済システム情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) UpdatePaymentSystem(ctx *gin.Context) {
	req := &request.UpdatePaymentSystemRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	methodType, err := util.GetParamInt32(ctx, "methodType")
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
