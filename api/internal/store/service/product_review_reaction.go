package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) UpsertProductReviewReaction(
	ctx context.Context, in *store.UpsertProductReviewReactionInput,
) (*entity.ProductReviewReaction, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	review, err := s.db.ProductReview.Get(ctx, in.ReviewID)
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewProductReviewReactionParams{
		ReviewID:     review.ID,
		UserID:       in.UserID,
		ReactionType: in.ReactionType,
	}
	reaction := entity.NewProductReviewReaction(params)
	if err := s.db.ProductReviewReaction.Upsert(ctx, reaction); err != nil {
		return nil, internalError(err)
	}
	return reaction, nil
}

func (s *service) DeleteProductReviewReaction(
	ctx context.Context,
	in *store.DeleteProductReviewReactionInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.ProductReviewReaction.Delete(ctx, in.ReviewID, in.UserID)
	return internalError(err)
}

func (s *service) GetUserProductReviewReactions(
	ctx context.Context, in *store.GetUserProductReviewReactionsInput,
) (entity.ProductReviewReactions, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	reactions, err := s.db.ProductReviewReaction.GetUserReactions(ctx, in.ProductID, in.UserID)
	return reactions, internalError(err)
}
