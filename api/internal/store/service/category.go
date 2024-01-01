package service

import (
	"context"
	"fmt"

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
		return nil, 0, internalError(err)
	}
	orders := make([]*database.ListCategoriesOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListCategoriesOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListCategoriesParams{
		Name:   in.Name,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
		Orders: orders,
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
		return nil, 0, internalError(err)
	}
	return categories, total, nil
}

func (s *service) MultiGetCategories(
	ctx context.Context, in *store.MultiGetCategoriesInput,
) (entity.Categories, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	categories, err := s.db.Category.MultiGet(ctx, in.CategoryIDs)
	return categories, internalError(err)
}

func (s *service) GetCategory(ctx context.Context, in *store.GetCategoryInput) (*entity.Category, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	category, err := s.db.Category.Get(ctx, in.CategoryID)
	return category, internalError(err)
}

func (s *service) CreateCategory(ctx context.Context, in *store.CreateCategoryInput) (*entity.Category, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewCategoryParams{
		Name: in.Name,
	}
	category := entity.NewCategory(params)
	if err := s.db.Category.Create(ctx, category); err != nil {
		return nil, internalError(err)
	}
	return category, nil
}

func (s *service) UpdateCategory(ctx context.Context, in *store.UpdateCategoryInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Category.Update(ctx, in.CategoryID, in.Name)
	return internalError(err)
}

func (s *service) DeleteCategory(ctx context.Context, in *store.DeleteCategoryInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.ListProductTypesParams{
		CategoryID: in.CategoryID,
	}
	total, err := s.db.ProductType.Count(ctx, params)
	if err != nil {
		return internalError(err)
	}
	if total > 0 {
		return fmt.Errorf("service: associated with product type: %w", exception.ErrFailedPrecondition)
	}
	err = s.db.Category.Delete(ctx, in.CategoryID)
	return internalError(err)
}
