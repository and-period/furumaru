package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

// @tag.name        ProductReview
// @tag.description 商品レビュー関連
func (h *handler) productReviewRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products/:productId/reviews", h.authentication)

	r.POST("", h.CreateDummyProductReview)
}

// @Summary     ダミー商品レビュー作成
// @Description ダミーの商品レビューを作成します。
// @Tags        ProductReview
// @Router      /v1/products/{productId}/reviews [post]
// @Security    bearerauth
// @Param       productId path string true "商品ID" example("product_12345678")
// @Produce     json
// @Param       request body types.CreateProductReviewRequest true "商品レビュー作成リクエスト"
// @Success     204 "作成成功"
// @Failure     400 {object} util.ErrorResponse "リクエスト不正"
// @Failure     404 {object} util.ErrorResponse "商品が存在しない"
func (h *handler) CreateDummyProductReview(ctx *gin.Context) {
	req := &types.CreateProductReviewRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	guest, err := h.user.GetDummyGuest(ctx, &user.GetDummyGuestInput{})
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &store.CreateProductReviewInput{
		ProductID: ctx.Param("productId"),
		UserID:    guest.UserID,
		Rate:      req.Rate,
		Title:     req.Title,
		Comment:   req.Comment,
	}
	if _, err := h.store.CreateProductReview(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
