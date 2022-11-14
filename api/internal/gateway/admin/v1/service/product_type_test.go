package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestProductType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		productType *entity.ProductType
		expect      *ProductType
	}{
		{
			name: "success",
			productType: &entity.ProductType{
				ID:         "product-type-id",
				CategoryID: "category-id",
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
				Icons: common.Images{
					{URL: "https://and-period.jp/icon_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/icon_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/icon_900.png", Size: common.ImageSizeLarge},
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &ProductType{
				ProductType: response.ProductType{
					ID:           "product-type-id",
					CategoryID:   "category-id",
					CategoryName: "",
					Name:         "じゃがいも",
					IconURL:      "https://and-period.jp/icon.png",
					Icons: []*response.Image{
						{URL: "https://and-period.jp/icon_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/icon_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/icon_900.png", Size: int32(ImageSizeLarge)},
					},
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
			assert.Equal(t, tt.expect, NewProductType(tt.productType))
		})
	}
}

func TestProductType_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		productType *ProductType
		category    *Category
		expect      *ProductType
	}{
		{
			name: "success",
			productType: &ProductType{
				ProductType: response.ProductType{
					ID:           "product-type-id",
					CategoryID:   "category-id",
					CategoryName: "",
					Name:         "じゃがいも",
					IconURL:      "https://and-period.jp/icon.png",
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
			category: &Category{
				Category: response.Category{
					ID:        "category-id",
					Name:      "野菜",
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
			expect: &ProductType{
				ProductType: response.ProductType{
					ID:           "product-type-id",
					CategoryID:   "category-id",
					CategoryName: "野菜",
					Name:         "じゃがいも",
					IconURL:      "https://and-period.jp/icon.png",
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.productType.Fill(tt.category)
			assert.Equal(t, tt.expect, tt.productType)
		})
	}
}

func TestProductType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		productType *ProductType
		expect      *response.ProductType
	}{
		{
			name: "success",
			productType: &ProductType{
				ProductType: response.ProductType{
					ID:           "product-type-id",
					CategoryID:   "category-id",
					CategoryName: "野菜",
					Name:         "じゃがいも",
					IconURL:      "https://and-period.jp/icon.png",
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
			expect: &response.ProductType{
				ID:           "product-type-id",
				CategoryID:   "category-id",
				CategoryName: "野菜",
				Name:         "じゃがいも",
				IconURL:      "https://and-period.jp/icon.png",
				CreatedAt:    1640962800,
				UpdatedAt:    1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.productType.Response())
		})
	}
}

func TestProductTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		productTypes entity.ProductTypes
		expect       ProductTypes
	}{
		{
			name: "success",
			productTypes: entity.ProductTypes{
				{
					ID:         "product-type-id",
					CategoryID: "category-id",
					Name:       "じゃがいも",
					IconURL:    "https://and-period.jp/icon.png",
					CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: ProductTypes{
				{
					ProductType: response.ProductType{
						ID:           "product-type-id",
						CategoryID:   "category-id",
						CategoryName: "",
						Name:         "じゃがいも",
						IconURL:      "https://and-period.jp/icon.png",
						Icons:        []*response.Image{},
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProductTypes(tt.productTypes))
		})
	}
}

func TestProductTypes_CategoryIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		productTypes ProductTypes
		expect       []string
	}{
		{
			name: "success",
			productTypes: ProductTypes{
				{
					ProductType: response.ProductType{
						ID:           "product-type-id",
						CategoryID:   "category-id",
						CategoryName: "野菜",
						Name:         "じゃがいも",
						IconURL:      "https://and-period.jp/icon.png",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			expect: []string{"category-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.productTypes.CategoryIDs())
		})
	}
}

func TestProductTypes_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		productTypes ProductTypes
		expect       map[string]*ProductType
	}{
		{
			name: "success",
			productTypes: ProductTypes{
				{
					ProductType: response.ProductType{
						ID:           "product-type-id",
						CategoryID:   "category-id",
						CategoryName: "野菜",
						Name:         "じゃがいも",
						IconURL:      "https://and-period.jp/icon.png",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			expect: map[string]*ProductType{
				"product-type-id": {
					ProductType: response.ProductType{
						ID:           "product-type-id",
						CategoryID:   "category-id",
						CategoryName: "野菜",
						Name:         "じゃがいも",
						IconURL:      "https://and-period.jp/icon.png",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.productTypes.Map())
		})
	}
}

func TestProductTypes_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		productTypes ProductTypes
		categories   map[string]*Category
		expect       ProductTypes
	}{
		{
			name: "success",
			productTypes: ProductTypes{
				{
					ProductType: response.ProductType{
						ID:           "product-type-id",
						CategoryID:   "category-id",
						CategoryName: "",
						Name:         "じゃがいも",
						IconURL:      "https://and-period.jp/icon.png",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
				{
					ProductType: response.ProductType{
						ID:           "other-id",
						CategoryID:   "other-id",
						CategoryName: "",
						Name:         "ほうれん草",
						IconURL:      "https://and-period.jp/icon.png",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			categories: map[string]*Category{
				"category-id": {
					Category: response.Category{
						ID:        "category-id",
						Name:      "野菜",
						CreatedAt: 1640962800,
						UpdatedAt: 1640962800,
					},
				},
			},
			expect: ProductTypes{
				{
					ProductType: response.ProductType{
						ID:           "product-type-id",
						CategoryID:   "category-id",
						CategoryName: "野菜",
						Name:         "じゃがいも",
						IconURL:      "https://and-period.jp/icon.png",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
				{
					ProductType: response.ProductType{
						ID:           "other-id",
						CategoryID:   "other-id",
						CategoryName: "",
						Name:         "ほうれん草",
						IconURL:      "https://and-period.jp/icon.png",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.productTypes.Fill(tt.categories)
			assert.Equal(t, tt.expect, tt.productTypes)
		})
	}
}

func TestProductTypes_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		productTypes ProductTypes
		expect       []*response.ProductType
	}{
		{
			name: "success",
			productTypes: ProductTypes{
				{
					ProductType: response.ProductType{
						ID:           "product-type-id",
						CategoryID:   "category-id",
						CategoryName: "野菜",
						Name:         "じゃがいも",
						IconURL:      "https://and-period.jp/icon.png",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			expect: []*response.ProductType{
				{
					ID:           "product-type-id",
					CategoryID:   "category-id",
					CategoryName: "野菜",
					Name:         "じゃがいも",
					IconURL:      "https://and-period.jp/icon.png",
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.productTypes.Response())
		})
	}
}
