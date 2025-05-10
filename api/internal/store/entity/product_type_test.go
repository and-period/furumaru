package entity

import (
	"testing"

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductType(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
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
					CategoryID: "category-id",
				},
			},
			expect: []string{"category-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.productTypes.CategoryIDs())
		})
	}
}
