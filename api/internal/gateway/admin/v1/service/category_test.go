package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
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
				Category: &response.Category{
					ID:        "category-id",
					Name:      "野菜",
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
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
		expect   *response.Category
	}{
		{
			name: "success",
			category: &Category{
				Category: &response.Category{
					ID:        "category-id",
					Name:      "野菜",
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
			expect: &response.Category{
				ID:        "category-id",
				Name:      "野菜",
				CreatedAt: 1640962800,
				UpdatedAt: 1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
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
					Category: &response.Category{
						ID:        "category-id",
						Name:      "野菜",
						CreatedAt: 1640962800,
						UpdatedAt: 1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewCategories(tt.categories))
		})
	}
}

func TestCategories_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		categories Categories
		expect     []*response.Category
	}{
		{
			name: "success",
			categories: Categories{
				{
					Category: &response.Category{
						ID:        "category-id",
						Name:      "野菜",
						CreatedAt: 1640962800,
						UpdatedAt: 1640962800,
					},
				},
			},
			expect: []*response.Category{
				{
					ID:        "category-id",
					Name:      "野菜",
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.categories.Response())
		})
	}
}
