package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		category *entity.Category
		expect   *Category
	}{
		{
			name: "success",
			category: &entity.Category{
				ID:        "category-id",
				Name:      "野菜",
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Category{
				Category: types.Category{
					ID:   "category-id",
					Name: "野菜",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewCategory(tt.category))
		})
	}
}

func TestCategory_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		category *Category
		expect   *types.Category
	}{
		{
			name: "success",
			category: &Category{
				Category: types.Category{
					ID:   "category-id",
					Name: "野菜",
				},
			},
			expect: &types.Category{
				ID:   "category-id",
				Name: "野菜",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.category.Response())
		})
	}
}

func TestCategories(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		categories entity.Categories
		expect     Categories
	}{
		{
			name: "success",
			categories: entity.Categories{
				{
					ID:        "category-id",
					Name:      "野菜",
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Categories{
				{
					Category: types.Category{
						ID:   "category-id",
						Name: "野菜",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewCategories(tt.categories))
		})
	}
}

func TestCategories_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		categories Categories
		expect     map[string]*Category
	}{
		{
			name: "success",
			categories: Categories{
				{
					Category: types.Category{
						ID:   "category-id",
						Name: "野菜",
					},
				},
			},
			expect: map[string]*Category{
				"category-id": {
					Category: types.Category{
						ID:   "category-id",
						Name: "野菜",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.categories.Map())
		})
	}
}

func TestCategories_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		categories Categories
		expect     []*types.Category
	}{
		{
			name: "success",
			categories: Categories{
				{
					Category: types.Category{
						ID:   "category-id",
						Name: "野菜",
					},
				},
			},
			expect: []*types.Category{
				{
					ID:   "category-id",
					Name: "野菜",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.categories.Response())
		})
	}
}
