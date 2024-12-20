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
	"github.com/gin-gonic/gin"
)

func (h *handler) productReviewRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products/:productId/reviews")

	r.GET("/:reviewId", h.GetProductReview)
	r.POST("", h.authentication, h.CreateProductReview)
	r.PATCH("/:reviewId", h.authentication, h.UpdateProductReview)
}

func (h *handler) GetProductReview(ctx *gin.Context) {
	review, err := h.getProductReview(ctx, util.GetParam(ctx, "reviewId"))
	if err != nil {
		h.httpError(ctx, err)
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
