package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

// @tag.name        ProductReview
// @tag.description 商品レビュー関連
func (h *handler) productReviewRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products/:productId/reviews", h.authentication)

	r.GET("", h.ListProductReviews)
	r.GET("/:reviewId", h.GetProductReview)
	r.POST("", h.CreateDummyProductReview)
}

// @Summary     商品レビュー一覧取得
// @Description 指定した商品のレビュー一覧を取得します。
// @Tags        ProductReview
// @Router      /products/{productId}/reviews [get]
// @Param       productId path string true "商品ID"
// @Param       limit query int64 false "取得件数" default(20)
// @Param       offset query int64 false "取得開始位置" default(0)
// @Produce     json
// @Success     200 {object} types.ProductReviewsResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) ListProductReviews(ctx *gin.Context) {
	const defaultLimit = 20
	rates, err := util.GetQueryInt64s(ctx, "rates")
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &store.ListProductReviewsInput{
		ProductID: util.GetParam(ctx, "productId"),
		UserID:    util.GetQuery(ctx, "userId", ""),
		Rates:     rates,
		Limit:     limit,
		NextToken: util.GetQuery(ctx, "nextToken", ""),
	}
	reviews, nextToken, err := h.store.ListProductReviews(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(reviews) == 0 {
		res := &types.ProductReviewsResponse{
			Reviews: []*types.ProductReview{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	users, err := h.multiGetUsers(ctx, reviews.UserIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.ProductReviewsResponse{
		Reviews:   service.NewProductReviews(reviews, users.Map()).Response(),
		NextToken: nextToken,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     商品レビュー詳細取得
// @Description 商品レビューの詳細情報を取得します。
// @Tags        ProductReview
// @Router      /products/{productId}/reviews/{reviewId} [get]
// @Param       productId path string true "商品ID"
// @Param       reviewId path string true "レビューID"
// @Produce     json
// @Success     200 {object} types.ProductReviewResponse
// @Failure     404 {object} util.ErrorResponse "レビューが見つからない"
func (h *handler) GetProductReview(ctx *gin.Context) {
	review, err := h.getProductReview(ctx, util.GetParam(ctx, "reviewId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if review.ProductID != util.GetParam(ctx, "productId") {
		h.notFound(ctx, errors.New("handler: review not found"))
		return
	}
	res := &types.ProductReviewResponse{
		Review: review.Response(),
	}
	ctx.JSON(http.StatusOK, res)
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

func (h *handler) getProductReview(ctx context.Context, reviewID string) (*service.ProductReview, error) {
	in := &store.GetProductReviewInput{
		ReviewID: reviewID,
	}
	review, err := h.store.GetProductReview(ctx, in)
	if err != nil {
		return nil, err
	}
	user, err := h.getUser(ctx, review.UserID)
	if err != nil {
		return nil, err
	}
	return service.NewProductReview(review, user), nil
}
