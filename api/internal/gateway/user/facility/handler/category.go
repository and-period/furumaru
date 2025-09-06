package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/service"
	"github.com/and-period/furumaru/api/internal/store"
)

func (h *handler) multiGetCategories(ctx context.Context, categoryIDs []string) (service.Categories, error) {
	if len(categoryIDs) == 0 {
		return service.Categories{}, nil
	}
	in := &store.MultiGetCategoriesInput{
		CategoryIDs: categoryIDs,
	}
	categories, err := h.store.MultiGetCategories(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCategories(categories), nil
}

func (h *handler) getCategory(ctx context.Context, categoryID string) (*service.Category, error) {
	in := &store.GetCategoryInput{
		CategoryID: categoryID,
	}
	category, err := h.store.GetCategory(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewCategory(category), nil
}
