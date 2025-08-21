package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/gin-gonic/gin"
)

// @tag.name        Checkout
// @tag.description チェックアウト関連
func (h *handler) checkoutRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/checkouts", h.authentication)

	r.POST("", h.Checkout)
	r.GET("/:transactionId", h.GetCheckoutState)
}

// @Summary     購入する
// @Description 商品を購入します。
// @Tags        Checkout
// @Router      /facilities/{facilityId}/checkouts [post]
// @Accept      json
// @Param				request body request.CheckoutRequest true "チェックアウト情報"
// @Produce     json
// @Success     200 {object} response.CheckoutResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) Checkout(ctx *gin.Context) {
	req := &request.CheckoutRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	res := &response.CheckoutResponse{}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     支払い状態取得
// @Description 支払い状態を取得します。
// @Tags        Checkout
// @Router      /facilities/{facilityId}/checkouts/{transactionId} [get]
// @Param       transactionId path string true "取引ID"
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} response.CheckoutStateResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     404 {object} util.ErrorResponse "取引が見つからない"
func (h *handler) GetCheckoutState(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.CheckoutStateResponse{}
	ctx.JSON(http.StatusOK, res)
}
