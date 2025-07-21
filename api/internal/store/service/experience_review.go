package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) ListExperienceReviews(
	ctx context.Context, in *store.ListExperienceReviewsInput,
) (entity.ExperienceReviews, string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, "", internalError(err)
	}
	params := &database.ListExperienceReviewsParams{
		ExperienceID: in.ExperienceID,
		UserID:       in.UserID,
		Rates:        in.Rates,
		Limit:        in.Limit,
		NextToken:    in.NextToken,
	}
	reviews, token, err := s.db.ExperienceReview.List(ctx, params)
	return reviews, token, internalError(err)
}

func (s *service) GetExperienceReview(
	ctx context.Context,
	in *store.GetExperienceReviewInput,
) (*entity.ExperienceReview, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	review, err := s.db.ExperienceReview.Get(ctx, in.ReviewID)
	return review, internalError(err)
}

func (s *service) CreateExperienceReview(
	ctx context.Context,
	in *store.CreateExperienceReviewInput,
) (*entity.ExperienceReview, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	if _, err := s.db.Experience.Get(ctx, in.ExperienceID); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewExperienceReviewParams{
		ExperienceID: in.ExperienceID,
		UserID:       in.UserID,
		Rate:         in.Rate,
		Title:        in.Title,
		Comment:      in.Comment,
	}
	review := entity.NewExperienceReview(params)
	if err := s.db.ExperienceReview.Create(ctx, review); err != nil {
		return nil, internalError(err)
	}
	return review, nil
}

func (s *service) UpdateExperienceReview(
	ctx context.Context,
	in *store.UpdateExperienceReviewInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateExperienceReviewParams{
		Rate:    in.Rate,
		Title:   in.Title,
		Comment: in.Comment,
	}
	err := s.db.ExperienceReview.Update(ctx, in.ReviewID, params)
	return internalError(err)
}

func (s *service) DeleteExperienceReview(
	ctx context.Context,
	in *store.DeleteExperienceReviewInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.ExperienceReview.Delete(ctx, in.ReviewID)
	return internalError(err)
}

func (s *service) AggregateExperienceReviews(
	ctx context.Context, in *store.AggregateExperienceReviewsInput,
) (entity.AggregatedExperienceReviews, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.AggregateExperienceReviewsParams{
		ExperienceIDs: in.ExperienceIDs,
	}
	aggregation, err := s.db.ExperienceReview.Aggregate(ctx, params)
	return aggregation, internalError(err)
}
