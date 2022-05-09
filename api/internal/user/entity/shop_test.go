package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		shopID        string
		cognitoID     string
		lastname      string
		firstname     string
		lastnameKana  string
		firstnameKana string
		email         string
		expect        *Shop
	}{
		{
			name:          "success",
			shopID:        "shop-id",
			cognitoID:     "cognito-id",
			lastname:      "&.",
			firstname:     "スタッフ",
			lastnameKana:  "あんどどっと",
			firstnameKana: "すたっふ",
			email:         "test-shop@and-period.jp",
			expect: &Shop{
				ID:            "shop-id",
				CognitoID:     "cognito-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-shop@and-period.jp",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := NewShop(
				tt.shopID, tt.cognitoID,
				tt.lastname, tt.firstname,
				tt.lastnameKana, tt.firstnameKana,
				tt.email,
			)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

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
