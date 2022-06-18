package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListCategories(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	params := &database.ListCategoriesParams{
		Name:   "野菜",
		Limit:  30,
		Offset: 0,
	}
	categories := entity.Categories{
		{
			ID:        "category-id",
			Name:      "野菜",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListCategoriesInput
		expect    entity.Categories
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().List(ctx, params).Return(categories, nil)
			},
			input: &store.ListCategoriesInput{
				Name:   "野菜",
				Limit:  30,
				Offset: 0,
			},
			expect:    categories,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.ListCategoriesInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().List(ctx, params).Return(nil, errmock)
			},
			input: &store.ListCategoriesInput{
				Name:   "野菜",
				Limit:  30,
				Offset: 0,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			actual, err := service.ListCategories(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetCategories(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	categories := entity.Categories{
		{
			ID:        "category-id",
			Name:      "野菜",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetCategoriesInput
		expect    entity.Categories
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().MultiGet(ctx, []string{"category-id"}).Return(categories, nil)
			},
			input: &store.MultiGetCategoriesInput{
				CategoryIDs: []string{"category-id"},
			},
			expect:    categories,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetCategoriesInput{
				CategoryIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().MultiGet(ctx, []string{"category-id"}).Return(nil, errmock)
			},
			input: &store.MultiGetCategoriesInput{
				CategoryIDs: []string{"category-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			actual, err := service.MultiGetCategories(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestCreateCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateCategoryInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, category *entity.Category) error {
						expect := &entity.Category{
							ID:   category.ID, // ignore
							Name: "野菜",
						}
						assert.Equal(t, expect, category)
						return nil
					})
			},
			input: &store.CreateCategoryInput{
				Name: "野菜",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateCategoryInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &store.CreateCategoryInput{
				Name: "野菜",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			_, err := service.CreateCategory(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateCategoryInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().Update(ctx, "category-id", "野菜").Return(nil)
			},
			input: &store.UpdateCategoryInput{
				CategoryID: "category-id",
				Name:       "野菜",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateCategoryInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().Update(ctx, "category-id", "野菜").Return(errmock)
			},
			input: &store.UpdateCategoryInput{
				CategoryID: "category-id",
				Name:       "野菜",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			err := service.UpdateCategory(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteCategoryInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().Delete(ctx, "category-id").Return(nil)
			},
			input: &store.DeleteCategoryInput{
				CategoryID: "category-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteCategoryInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Category.EXPECT().Delete(ctx, "category-id").Return(errmock)
			},
			input: &store.DeleteCategoryInput{
				CategoryID: "category-id",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			err := service.DeleteCategory(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
