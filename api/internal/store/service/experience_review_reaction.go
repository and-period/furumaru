package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) UpsertExperienceReviewReaction(
	ctx context.Context, in *store.UpsertExperienceReviewReactionInput,
) (*entity.ExperienceReviewReaction, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	review, err := s.db.ExperienceReview.Get(ctx, in.ReviewID)
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewExperienceReviewReactionParams{
		ReviewID:     review.ID,
		UserID:       in.UserID,
		ReactionType: in.ReactionType,
	}
	reaction := entity.NewExperienceReviewReaction(params)
	if err := s.db.ExperienceReviewReaction.Upsert(ctx, reaction); err != nil {
		return nil, internalError(err)
	}
	return reaction, nil
}

func (s *service) DeleteExperienceReviewReaction(
	ctx context.Context,
	in *store.DeleteExperienceReviewReactionInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.ExperienceReviewReaction.Delete(ctx, in.ReviewID, in.UserID)
	return internalError(err)
}

func (s *service) GetUserExperienceReviewReactions(
	ctx context.Context, in *store.GetUserExperienceReviewReactionsInput,
) (entity.ExperienceReviewReactions, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	reactions, err := s.db.ExperienceReviewReaction.GetUserReactions(
		ctx,
		in.ExperienceID,
		in.UserID,
	)
	return reactions, internalError(err)
}
