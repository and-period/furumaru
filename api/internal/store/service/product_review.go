package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) GetProductReview(ctx context.Context, in *store.GetProductReviewInput) (*entity.ProductReview, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	review, err := s.db.ProductReview.Get(ctx, in.ReviewID)
	return review, internalError(err)
}

func (s *service) CreateProductReview(ctx context.Context, in *store.CreateProductReviewInput) (*entity.ProductReview, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	if _, err := s.db.Product.Get(ctx, in.ProductID); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewProductReviewParams{
		ProductID: in.ProductID,
		UserID:    in.UserID,
		Rate:      in.Rate,
		Title:     in.Title,
		Comment:   in.Comment,
	}
	review := entity.NewProductReview(params)
	if err := s.db.ProductReview.Create(ctx, review); err != nil {
		return nil, internalError(err)
	}
	return review, nil
}

func (s *service) UpdateProductReview(ctx context.Context, in *store.UpdateProductReviewInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateProductReviewParams{
		Rate:    in.Rate,
		Title:   in.Title,
		Comment: in.Comment,
	}
	err := s.db.ProductReview.Update(ctx, in.ReviewID, params)
	return internalError(err)
}

func (s *service) DeleteProductReview(ctx context.Context, in *store.DeleteProductReviewInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.ProductReview.Delete(ctx, in.ReviewID)
	return internalError(err)
}
