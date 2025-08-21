package handler

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/gin-gonic/gin"
)

// @tag.name        Cart
// @tag.description カート関連
func (h *handler) cartRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/carts", h.authentication)

	r.GET("", h.GetCart)
	r.POST("/-/items", h.AddCartItem)
	r.DELETE("/-/items/:productId", h.RemoveCartItem)
}

// @Summary     カート取得
// @Description カートの内容を取得します。
// @Tags        Cart
// @Router      /facilities/{facilityId}/carts [get]
// @Security    bearerauth
// @Produce     json
// @Success     200 {object} response.CartResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) GetCart(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.CartResponse{
		Carts:        []*response.Cart{},
		Coordinators: []*response.Coordinator{},
		Products:     []*response.Product{},
	}
	ctx.JSON(200, res)
}

// @Summary     カートに追加
// @Description カートに商品を追加します。
// @Tags        Cart
// @Router      /facilities/{facilityId}/carts/-/items [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body request.AddCartItemRequest true "カートに追加リクエスト"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) AddCartItem(ctx *gin.Context) {
	req := &request.AddCartItemRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	ctx.Status(204)
}

// @Summary     カートから削除
// @Description カートから商品を削除します。
// @Tags        Cart
// @Router      /facilities/{facilityId}/carts/-/items/:productId [delete]
// @Security    bearerauth
// @Produce     json
// @Success     204
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) RemoveCartItem(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.Status(204)
}
