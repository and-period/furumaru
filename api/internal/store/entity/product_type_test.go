package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestProductType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		params     *NewProductTypeParams
		input      string
		iconURL    string
		categoryID string
		expect     *ProductType
	}{
		{
			name: "success",
			params: &NewProductTypeParams{
				CategoryID: "category-id",
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
			},
			expect: &ProductType{
				CategoryID: "category-id",
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductType(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductType_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		productType *ProductType
		expect      *ProductType
		hasErr      bool
	}{
		{
			name: "success",
			productType: &ProductType{
				ID:         "product-type-id",
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
				IconsJSON:  []byte(`[{"url":"http://example.com/media.png","size":1}]`),
				CategoryID: "category-id",
			},
			expect: &ProductType{
				ID:      "product-type-id",
				Name:    "じゃがいも",
				IconURL: "https://and-period.jp/icon.png",
				Icons: common.Images{
					{
						Size: common.ImageSizeSmall,
						URL:  "http://example.com/media.png",
					},
				},
				IconsJSON:  []byte(`[{"url":"http://example.com/media.png","size":1}]`),
				CategoryID: "category-id",
			},
			hasErr: false,
		},
		{
			name: "failed to marshal json",
			productType: &ProductType{
				ID:         "product-type-id",
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
				IconsJSON:  []byte(`{{{`),
				CategoryID: "category-id",
			},
			expect: &ProductType{
				ID:         "product-type-id",
				Name:       "じゃがいも",
				IconURL:    "https://and-period.jp/icon.png",
				IconsJSON:  []byte(`{{{`),
				CategoryID: "category-id",
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.productType.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.productType)
		})
	}
}

func TestProductTypes_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		productTypes ProductTypes
		expect       ProductTypes
		hasErr       bool
	}{
		{
			name: "success",
			productTypes: ProductTypes{
				{
					ID:         "product-type-id",
					Name:       "じゃがいも",
					IconURL:    "https://and-period.jp/icon.png",
					IconsJSON:  []byte(`[{"url":"http://example.com/media.png","size":1}]`),
					CategoryID: "category-id",
				},
			},
			expect: ProductTypes{
				{
					ID:      "product-type-id",
					Name:    "じゃがいも",
					IconURL: "https://and-period.jp/icon.png",
					Icons: common.Images{
						{
							Size: common.ImageSizeSmall,
							URL:  "http://example.com/media.png",
						},
					},
					IconsJSON:  []byte(`[{"url":"http://example.com/media.png","size":1}]`),
					CategoryID: "category-id",
				},
			},
			hasErr: false,
		},
		{
			name: "failed to marshal json",
			productTypes: ProductTypes{
				{
					ID:         "product-type-id",
					Name:       "じゃがいも",
					IconURL:    "https://and-period.jp/icon.png",
					IconsJSON:  []byte(`{{{`),
					CategoryID: "category-id",
				},
			},
			expect: ProductTypes{
				{
					ID:         "product-type-id",
					Name:       "じゃがいも",
					IconURL:    "https://and-period.jp/icon.png",
					IconsJSON:  []byte(`{{{`),
					CategoryID: "category-id",
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.productTypes.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.productTypes)
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
					ID:         "product-type-id",
					Name:       "じゃがいも",
					IconURL:    "https://and-period.jp/icon.png",
					IconsJSON:  []byte(`[{"url":"http://example.com/media.png","size":1}]`),
					CategoryID: "category-id",
				},
			},
			expect: []string{"category-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.productTypes.CategoryIDs())
		})
	}
}
