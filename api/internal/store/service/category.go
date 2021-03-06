package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListCategories(
	ctx context.Context, in *store.ListCategoriesInput,
) (entity.Categories, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	params := &database.ListCategoriesParams{
		Name:   in.Name,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		categories entity.Categories
		total      int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		categories, err = s.db.Category.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Category.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return categories, total, nil
}

func (s *service) MultiGetCategories(
	ctx context.Context, in *store.MultiGetCategoriesInput,
) (entity.Categories, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	categories, err := s.db.Category.MultiGet(ctx, in.CategoryIDs)
	return categories, exception.InternalError(err)
}

func (s *service) CreateCategory(ctx context.Context, in *store.CreateCategoryInput) (*entity.Category, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	category := entity.NewCategory(in.Name)
	if err := s.db.Category.Create(ctx, category); err != nil {
		return nil, exception.InternalError(err)
	}
	return category, nil
}

func (s *service) UpdateCategory(ctx context.Context, in *store.UpdateCategoryInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Category.Update(ctx, in.CategoryID, in.Name)
	return exception.InternalError(err)
}

func (s *service) DeleteCategory(ctx context.Context, in *store.DeleteCategoryInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Category.Delete(ctx, in.CategoryID)
	return exception.InternalError(err)
}
