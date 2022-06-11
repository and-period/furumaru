package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *storeService) ListCategories(ctx context.Context, in *store.ListCategoriesInput) (entity.Categories, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &database.ListCategoriesParams{
		Name:   in.Name,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	categories, err := s.db.Category.List(ctx, params)
	return categories, exception.InternalError(err)
}

func (s *storeService) CreateCategory(ctx context.Context, in *store.CreateCategoryInput) (*entity.Category, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	category := entity.NewCategory(in.Name)
	if err := s.db.Category.Create(ctx, category); err != nil {
		return nil, exception.InternalError(err)
	}
	return category, nil
}

func (s *storeService) UpdateCategory(ctx context.Context, in *store.UpdateCategoryInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Category.Update(ctx, in.CategoryID, in.Name)
	return exception.InternalError(err)
}

func (s *storeService) DeleteCategory(ctx context.Context, in *store.DeleteCategoryInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Category.Delete(ctx, in.CategoryID)
	return exception.InternalError(err)
}
