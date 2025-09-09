package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        ProductReview
// @tag.description 商品レビュー関連
func (h *handler) productReviewRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products/:productId/reviews")

	r.GET("", h.ListProductReviews)
	r.GET("/:reviewId", h.GetProductReview)
	r.POST("", h.authentication, h.CreateProductReview)
	r.PATCH("/:reviewId", h.authentication, h.UpdateProductReview)
	r.POST("/:reviewId/reactions", h.authentication, h.UpsertProductReviewReaction)
	r.DELETE("/:reviewId/reactions", h.authentication, h.DeleteProductReviewReaction)

	auth := rg.Group("/users/me/products/:productId/reviews", h.authentication)
	auth.GET("", h.ListUserProductReviews)
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

// @Summary     商品レビュー作成
// @Description 商品のレビューを作成します。
// @Tags        ProductReview
// @Router      /products/{productId}/reviews [post]
// @Security    bearerauth
// @Param       productId path string true "商品ID"
// @Accept      json
// @Param       request body types.CreateProductReviewRequest true "レビュー作成"
// @Success     204 "作成成功"
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "商品が存在しない"
func (h *handler) CreateProductReview(ctx *gin.Context) {
	req := &types.CreateProductReviewRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	user, err := h.getMember(ctx, h.getUserID(ctx))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	in := &store.CreateProductReviewInput{
		ProductID: util.GetParam(ctx, "productId"),
		UserID:    user.ID,
		Rate:      req.Rate,
		Title:     req.Title,
		Comment:   req.Comment,
	}
	review, err := h.store.CreateProductReview(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.ProductReviewResponse{
		Review: service.NewProductReview(review, user).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     商品レビュー更新
// @Description 商品レビューの内容を更新します。
// @Tags        ProductReview
// @Router      /products/{productId}/reviews/{reviewId} [patch]
// @Security    bearerauth
// @Param       productId path string true "商品ID"
// @Param       reviewId path string true "レビューID"
// @Accept      json
// @Param       request body types.UpdateProductReviewRequest true "レビュー更新"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     404 {object} util.ErrorResponse "レビューが見つからない"
func (h *handler) UpdateProductReview(ctx *gin.Context) {
	req := &types.UpdateProductReviewRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	review, err := h.getProductReview(ctx, util.GetParam(ctx, "reviewId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if review.UserID != h.getUserID(ctx) {
		h.forbidden(ctx, errors.New("handler: user is not owner"))
		return
	}
	in := &store.UpdateProductReviewInput{
		ReviewID: review.ID,
		Rate:     req.Rate,
		Title:    req.Title,
		Comment:  req.Comment,
	}
	if err := h.store.UpdateProductReview(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteProductReview(ctx *gin.Context) {
	review, err := h.getProductReview(ctx, util.GetParam(ctx, "reviewId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if review.UserID != h.getUserID(ctx) {
		h.forbidden(ctx, errors.New("handler: user is not owner"))
		return
	}
	in := &store.DeleteProductReviewInput{
		ReviewID: review.ID,
	}
	if err := h.store.DeleteProductReview(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     商品レビューリアクション登録/更新
// @Description 商品レビューに対するリアクション（いいねなど）を登録または更新します。
// @Tags        ProductReview
// @Router      /products/{productId}/reviews/{reviewId}/reactions [post]
// @Security    bearerauth
// @Param       productId path string true "商品ID"
// @Param       reviewId path string true "レビューID"
// @Accept      json
// @Param       request body types.UpsertProductReviewReactionRequest true "リアクション登録/更新"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) UpsertProductReviewReaction(ctx *gin.Context) {
	req := &types.UpsertProductReviewReactionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	reactionType, ok := service.NewProductReviewReactionTypeFromRequest(req.ReactionType)
	if !ok {
		h.badRequest(ctx, errors.New("handler: invalid reaction type"))
		return
	}
	in := &store.UpsertProductReviewReactionInput{
		ReviewID:     util.GetParam(ctx, "reviewId"),
		UserID:       h.getUserID(ctx),
		ReactionType: reactionType.StoreEntity(),
	}
	if _, err := h.store.UpsertProductReviewReaction(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     商品レビューリアクション削除
// @Description 商品レビューに対するリアクションを削除します。
// @Tags        ProductReview
// @Router      /products/{productId}/reviews/{reviewId}/reactions [delete]
// @Security    bearerauth
// @Param       productId path string true "商品ID"
// @Param       reviewId path string true "レビューID"
// @Produce     json
// @Success     204
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) DeleteProductReviewReaction(ctx *gin.Context) {
	in := &store.DeleteProductReviewReactionInput{
		ReviewID: util.GetParam(ctx, "reviewId"),
		UserID:   h.getUserID(ctx),
	}
	if err := h.store.DeleteProductReviewReaction(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary     ユーザー商品レビュー一覧取得
// @Description ログイン中のユーザーが投稿した商品レビューの一覧を取得します。
// @Tags        ProductReview
// @Router      /users/me/products/{productId}/reviews [get]
// @Security    bearerauth
// @Param       productId path string true "商品ID"
// @Produce     json
// @Success     200 {object} types.ProductReviewsResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
func (h *handler) ListUserProductReviews(ctx *gin.Context) {
	productID := util.GetParam(ctx, "productId")

	user, err := h.getMember(ctx, h.getUserID(ctx))
	if err != nil {
		h.httpError(ctx, err)
	}

	var (
		reviews   service.ProductReviews
		reactions service.ProductReviewReactions
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		in := &store.ListProductReviewsInput{
			ProductID: productID,
			UserID:    user.ID,
			NoLimit:   true,
		}
		sreviews, _, err := h.store.ListProductReviews(ectx, in)
		if err != nil {
			return err
		}
		users := map[string]*entity.User{user.ID: user}
		reviews = service.NewProductReviews(sreviews, users)
		return nil
	})
	eg.Go(func() error {
		in := &store.GetUserProductReviewReactionsInput{
			ProductID: productID,
			UserID:    user.ID,
		}
		sreactions, err := h.store.GetUserProductReviewReactions(ectx, in)
		if err != nil {
			return err
		}
		reactions = service.NewProductReviewReactions(sreactions)
		return nil
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.UserProductReviewsResponse{
		Reviews:   reviews.Response(),
		Reactions: reactions.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getProductReview(ctx context.Context, reviewID string) (*service.ProductReview, error) {
	in := &store.GetProductReviewInput{
		ReviewID: reviewID,
	}
	review, err := h.store.GetProductReview(ctx, in)
	if err != nil {
		return nil, err
	}
	user, err := h.getMember(ctx, review.UserID)
	if err != nil {
		return nil, err
	}
	return service.NewProductReview(review, user), nil
}

func (h *handler) aggregateProductRates(ctx context.Context, productIDs ...string) (service.ProductRates, error) {
	if len(productIDs) == 0 {
		return service.ProductRates{}, nil
	}
	in := &store.AggregateProductReviewsInput{
		ProductIDs: productIDs,
	}
	reviews, err := h.store.AggregateProductReviews(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewProductRates(reviews), nil
}
