package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
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
				ProductTag: response.ProductTag{
					ID:   "product-tag-id",
					Name: "野菜",
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
		expect     *response.ProductTag
	}{
		{
			name: "success",
			productTag: &ProductTag{
				ProductTag: response.ProductTag{
					ID:   "product-tag-id",
					Name: "野菜",
				},
			},
			expect: &response.ProductTag{
				ID:   "product-tag-id",
				Name: "野菜",
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
					ProductTag: response.ProductTag{
						ID:   "product-tag-id",
						Name: "野菜",
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
		expect      []*response.ProductTag
	}{
		{
			name: "success",
			productTags: ProductTags{
				{
					ProductTag: response.ProductTag{
						ID:   "product-tag-id",
						Name: "野菜",
					},
				},
			},
			expect: []*response.ProductTag{
				{
					ID:   "product-tag-id",
					Name: "野菜",
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
