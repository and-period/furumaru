package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/gin-gonic/gin"
)

// @tag.name        Order
// @tag.description 注文関連
func (h *handler) orderRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/orders", h.authentication)

	r.GET("", h.ListOrders)
	r.GET("/:orderId", h.GetOrder)
}

// @Summary     注文一覧取得
// @Description 注文の一覧を取得します。
// @Tags        Order
// @Router      /facilities/{facilityId}/orders [get]
// @Security    bearerauth
// @Param       limit query int64 false "取得件数" default(20)
// @Param       offset query int64 false "取得開始位置" default(0)
// @Param       types query []int32 false "注文ステータス" collectionFormat(csv)
// @Produce     json
// @Success     200 {object} response.OrdersResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) ListOrders(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.OrdersResponse{}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     注文詳細取得
// @Description 注文の詳細を取得します。
// @Tags        Order
// @Router      /facilities/{facilityId}/orders/{orderId} [get]
// @Param       orderId path string true "注文ID"
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} response.OrderResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     404 {object} util.ErrorResponse "注文が見つからない"
func (h *handler) GetOrder(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.OrderResponse{}
	ctx.JSON(http.StatusOK, res)
}
