package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/stretchr/testify/assert"
)

func TestShipping_Response(t *testing.T) {
	t.Parallel()
	rates := []*response.ShippingRate{
		{
			Number:      1,
			Name:        "四国",
			Price:       250,
			Prefectures: []string{"tokushima", "kagawa", "ehime", "kochi"},
		},
	}
	tests := []struct {
		name     string
		shipping *Shipping
		expect   *response.Shipping
	}{
		{
			name: "success",
			shipping: &Shipping{
				Shipping: response.Shipping{
					ID:                 "shipping-id",
					Name:               "デフォルト配送設定",
					Box60Rates:         rates,
					Box60Refrigerated:  500,
					Box60Frozen:        800,
					Box80Rates:         rates,
					Box80Refrigerated:  500,
					Box80Frozen:        800,
					Box100Rates:        rates,
					Box100Refrigerated: 500,
					Box100Frozen:       800,
					HasFreeShipping:    true,
					FreeShippingRates:  3000,
					CreatedAt:          1640962800,
					UpdatedAt:          1640962800,
				},
			},
			expect: &response.Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60Rates:         rates,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         rates,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        rates,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          1640962800,
				UpdatedAt:          1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shipping.Response())
		})
	}
}

func TestShippings_Response(t *testing.T) {
	t.Parallel()
	rates := []*response.ShippingRate{
		{
			Number:      1,
			Name:        "四国",
			Price:       250,
			Prefectures: []string{"tokushima", "kagawa", "ehime", "kochi"},
		},
	}
	tests := []struct {
		name      string
		shippings Shippings
		expect    []*response.Shipping
	}{
		{
			name: "success",
			shippings: Shippings{
				{
					Shipping: response.Shipping{
						ID:                 "shipping-id",
						Name:               "デフォルト配送設定",
						Box60Rates:         rates,
						Box60Refrigerated:  500,
						Box60Frozen:        800,
						Box80Rates:         rates,
						Box80Refrigerated:  500,
						Box80Frozen:        800,
						Box100Rates:        rates,
						Box100Refrigerated: 500,
						Box100Frozen:       800,
						HasFreeShipping:    true,
						FreeShippingRates:  3000,
						CreatedAt:          1640962800,
						UpdatedAt:          1640962800,
					},
				},
			},
			expect: []*response.Shipping{
				{
					ID:                 "shipping-id",
					Name:               "デフォルト配送設定",
					Box60Rates:         rates,
					Box60Refrigerated:  500,
					Box60Frozen:        800,
					Box80Rates:         rates,
					Box80Refrigerated:  500,
					Box80Frozen:        800,
					Box100Rates:        rates,
					Box100Refrigerated: 500,
					Box100Frozen:       800,
					HasFreeShipping:    true,
					FreeShippingRates:  3000,
					CreatedAt:          1640962800,
					UpdatedAt:          1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shippings.Response())
		})
	}
}
