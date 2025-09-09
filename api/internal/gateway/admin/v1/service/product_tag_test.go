package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestProductTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		productTag *entity.ProductTag
		expect     *ProductTag
	}{
		{
			name: "success",
			productTag: &entity.ProductTag{
				ID:        "product-tag-id",
				Name:      "野菜",
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &ProductTag{
				ProductTag: types.ProductTag{
					ID:        "product-tag-id",
					Name:      "野菜",
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProductTag(tt.productTag))
		})
	}
}

func TestProductTag_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		productTag *ProductTag
		expect     *types.ProductTag
	}{
		{
			name: "success",
			productTag: &ProductTag{
				ProductTag: types.ProductTag{
					ID:        "product-tag-id",
					Name:      "野菜",
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
			expect: &types.ProductTag{
				ID:        "product-tag-id",
				Name:      "野菜",
				CreatedAt: 1640962800,
				UpdatedAt: 1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.productTag.Response())
		})
	}
}

func TestProductTags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		productTags entity.ProductTags
		expect      ProductTags
	}{
		{
			name: "success",
			productTags: entity.ProductTags{
				{
					ID:        "product-tag-id",
					Name:      "野菜",
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: ProductTags{
				{
					ProductTag: types.ProductTag{
						ID:        "product-tag-id",
						Name:      "野菜",
						CreatedAt: 1640962800,
						UpdatedAt: 1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProductTags(tt.productTags))
		})
	}
}

func TestProductTags_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		productTags ProductTags
		expect      []*types.ProductTag
	}{
		{
			name: "success",
			productTags: ProductTags{
				{
					ProductTag: types.ProductTag{
						ID:        "product-tag-id",
						Name:      "野菜",
						CreatedAt: 1640962800,
						UpdatedAt: 1640962800,
					},
				},
			},
			expect: []*types.ProductTag{
				{
					ID:        "product-tag-id",
					Name:      "野菜",
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.productTags.Response())
		})
	}
}
