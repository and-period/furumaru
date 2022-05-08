package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShop_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		shop   *Shop
		expect string
	}{
		{
			name: "success",
			shop: &Shop{
				Lastname:  "&.",
				Firstname: "スタッフ",
			},
			expect: "&. スタッフ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.shop.Name())
		})
	}
}

func TestShops_Map(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		shops  Shops
		expect map[string]*Shop
	}{
		{
			name: "success",
			shops: Shops{
				{
					ID:        "shop-id",
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
			},
			expect: map[string]*Shop{
				"shop-id": {
					ID:        "shop-id",
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.shops.Map())
		})
	}
}
