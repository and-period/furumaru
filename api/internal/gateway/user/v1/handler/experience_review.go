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

func (h *handler) experienceReviewRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/experiences/:experienceId/reviews")

	r.GET("", h.ListExperienceReviews)
	r.GET("/:reviewId", h.GetExperienceReview)
	r.POST("", h.authentication, h.CreateExperienceReview)
	r.PATCH("/:reviewId", h.authentication, h.UpdateExperienceReview)
	r.POST("/:reviewId/reactions", h.authentication, h.UpsertExperienceReviewReaction)
	r.DELETE("/:reviewId/reactions", h.authentication, h.DeleteExperienceReviewReaction)

	auth := rg.Group("/users/me/experiences/:experienceId/reviews", h.authentication)
	auth.GET("", h.ListUserExperienceReviews)
}

func (h *handler) ListExperienceReviews(ctx *gin.Context) {
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

	in := &store.ListExperienceReviewsInput{
		ExperienceID: util.GetParam(ctx, "experienceId"),
		UserID:       util.GetQuery(ctx, "userId", ""),
		Rates:        rates,
		Limit:        limit,
		NextToken:    util.GetQuery(ctx, "nextToken", ""),
	}
	reviews, nextToken, err := h.store.ListExperienceReviews(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(reviews) == 0 {
		res := &response.ExperienceReviewsResponse{
			Reviews: []*response.ExperienceReview{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	users, err := h.multiGetUsers(ctx, reviews.UserIDs())
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ExperienceReviewsResponse{
		Reviews:   service.NewExperienceReviews(reviews, users.Map()).Response(),
		NextToken: nextToken,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetExperienceReview(ctx *gin.Context) {
	review, err := h.getExperienceReview(ctx, util.GetParam(ctx, "reviewId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if review.ExperienceID != util.GetParam(ctx, "experienceId") {
		h.notFound(ctx, errors.New("handler: review not found"))
		return
	}
	res := &response.ExperienceReviewResponse{
		Review: review.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateExperienceReview(ctx *gin.Context) {
	req := &request.CreateExperienceReviewRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	user, err := h.getMember(ctx, h.getUserID(ctx))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	in := &store.CreateExperienceReviewInput{
		ExperienceID: util.GetParam(ctx, "experienceId"),
		UserID:       user.ID,
		Rate:         req.Rate,
		Title:        req.Title,
		Comment:      req.Comment,
	}
	review, err := h.store.CreateExperienceReview(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ExperienceReviewResponse{
		Review: service.NewExperienceReview(review, user).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateExperienceReview(ctx *gin.Context) {
	req := &request.UpdateExperienceReviewRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	review, err := h.getExperienceReview(ctx, util.GetParam(ctx, "reviewId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if review.UserID != h.getUserID(ctx) {
		h.forbidden(ctx, errors.New("handler: user is not owner"))
		return
	}
	in := &store.UpdateExperienceReviewInput{
		ReviewID: review.ID,
		Rate:     req.Rate,
		Title:    req.Title,
		Comment:  req.Comment,
	}
	if err := h.store.UpdateExperienceReview(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteExperienceReview(ctx *gin.Context) {
	review, err := h.getExperienceReview(ctx, util.GetParam(ctx, "reviewId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if review.UserID != h.getUserID(ctx) {
		h.forbidden(ctx, errors.New("handler: user is not owner"))
		return
	}
	in := &store.DeleteExperienceReviewInput{
		ReviewID: review.ID,
	}
	if err := h.store.DeleteExperienceReview(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) UpsertExperienceReviewReaction(ctx *gin.Context) {
	req := &request.UpsertExperienceReviewReactionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}
	reactionType, ok := service.NewExperienceReviewReactionTypeFromRequest(req.ReactionType)
	if !ok {
		h.badRequest(ctx, errors.New("handler: invalid reaction type"))
		return
	}
	in := &store.UpsertExperienceReviewReactionInput{
		ReviewID:     util.GetParam(ctx, "reviewId"),
		UserID:       h.getUserID(ctx),
		ReactionType: reactionType.StoreEntity(),
	}
	if _, err := h.store.UpsertExperienceReviewReaction(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteExperienceReviewReaction(ctx *gin.Context) {
	in := &store.DeleteExperienceReviewReactionInput{
		ReviewID: util.GetParam(ctx, "reviewId"),
		UserID:   h.getUserID(ctx),
	}
	if err := h.store.DeleteExperienceReviewReaction(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *handler) ListUserExperienceReviews(ctx *gin.Context) {
	experienceID := util.GetParam(ctx, "experienceId")

	user, err := h.getMember(ctx, h.getUserID(ctx))
	if err != nil {
		h.httpError(ctx, err)
	}

	var (
		reviews   service.ExperienceReviews
		reactions service.ExperienceReviewReactions
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		in := &store.ListExperienceReviewsInput{
			ExperienceID: experienceID,
			UserID:       user.ID,
			NoLimit:      true,
		}
		sreviews, _, err := h.store.ListExperienceReviews(ectx, in)
		if err != nil {
			return err
		}
		users := map[string]*entity.User{user.ID: user}
		reviews = service.NewExperienceReviews(sreviews, users)
		return nil
	})
	eg.Go(func() error {
		in := &store.GetUserExperienceReviewReactionsInput{
			ExperienceID: experienceID,
			UserID:       user.ID,
		}
		sreactions, err := h.store.GetUserExperienceReviewReactions(ectx, in)
		if err != nil {
			return err
		}
		reactions = service.NewExperienceReviewReactions(sreactions)
		return nil
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.UserExperienceReviewsResponse{
		Reviews:   reviews.Response(),
		Reactions: reactions.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) getExperienceReview(
	ctx context.Context,
	reviewID string,
) (*service.ExperienceReview, error) {
	in := &store.GetExperienceReviewInput{
		ReviewID: reviewID,
	}
	review, err := h.store.GetExperienceReview(ctx, in)
	if err != nil {
		return nil, err
	}
	user, err := h.getMember(ctx, review.UserID)
	if err != nil {
		return nil, err
	}
	return service.NewExperienceReview(review, user), nil
}

func (h *handler) aggregateExperienceRates(
	ctx context.Context,
	experienceIDs ...string,
) (service.ExperienceRates, error) {
	in := &store.AggregateExperienceReviewsInput{
		ExperienceIDs: experienceIDs,
	}
	reviews, err := h.store.AggregateExperienceReviews(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewExperienceRates(reviews), nil
}
