package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

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
		res := &response.ProductReviewsResponse{
			Reviews: []*response.ProductReview{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	users, err := h.multiGetUsers(ctx, reviews.UserIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProductReviewsResponse{
		Reviews:   service.NewProductReviews(reviews, users.Map()).Response(),
		NextToken: nextToken,
	}
	ctx.JSON(http.StatusOK, res)
}

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
	res := &response.ProductReviewResponse{
		Review: review.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProductReview(ctx *gin.Context) {
	req := &request.CreateProductReviewRequest{}
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
	res := &response.ProductReviewResponse{
		Review: service.NewProductReview(review, user).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProductReview(ctx *gin.Context) {
	req := &request.UpdateProductReviewRequest{}
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

func (h *handler) UpsertProductReviewReaction(ctx *gin.Context) {
	req := &request.UpsertProductReviewReactionRequest{}
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

	res := &response.UserProductReviewsResponse{
		Reviews:   reviews.Response(),
		Reactions: reactions.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getProductReview(
	ctx context.Context,
	reviewID string,
) (*service.ProductReview, error) {
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

func (h *handler) aggregateProductRates(
	ctx context.Context,
	productIDs ...string,
) (service.ProductRates, error) {
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
