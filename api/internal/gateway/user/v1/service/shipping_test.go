package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestShipping(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		shipping *entity.Shipping
		expect   *Shipping
	}{
		{
			name: "success",
			shipping: &entity.Shipping{
				ID:            "shipping-id",
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				Name:          "配送設定",
				CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				ShippingRevision: entity.ShippingRevision{
					ID:         1,
					ShippingID: "shipping-id",
					Box60Rates: entity.ShippingRates{
						{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
					},
					Box60Frozen: 800,
					Box80Rates: entity.ShippingRates{
						{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
					},
					Box80Frozen: 800,
					Box100Rates: entity.ShippingRates{
						{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
					},
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
				},
			},
			expect: &Shipping{
				Shipping: types.Shipping{
					ID: "shipping-id",
					Box60Rates: []*types.ShippingRate{
						{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
					},
					Box60Frozen: 800,
					Box80Rates: []*types.ShippingRate{
						{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
					},
					Box80Frozen: 800,
					Box100Rates: []*types.ShippingRate{
						{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
					},
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShipping(tt.shipping)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShipping_Response(t *testing.T) {
	t.Parallel()
	rates := []*types.ShippingRate{
		{
			Number:          1,
			Name:            "四国",
			Price:           250,
			PrefectureCodes: []int32{36, 37, 38, 39},
		},
	}
	tests := []struct {
		name     string
		shipping *Shipping
		expect   *types.Shipping
	}{
		{
			name: "success",
			shipping: &Shipping{
				Shipping: types.Shipping{
					ID:                "shipping-id",
					Box60Rates:        rates,
					Box60Frozen:       800,
					Box80Rates:        rates,
					Box80Frozen:       800,
					Box100Rates:       rates,
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
				},
			},
			expect: &types.Shipping{
				ID:                "shipping-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shipping.Response())
		})
	}
}

func TestShippingRate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rate   *entity.ShippingRate
		expect *ShippingRate
	}{
		{
			name: "success",
			rate: &entity.ShippingRate{
				Number:          1,
				Name:            "東京都",
				Price:           1200,
				PrefectureCodes: []int32{13},
			},
			expect: &ShippingRate{
				ShippingRate: types.ShippingRate{
					Number:          1,
					Name:            "東京都",
					Price:           1200,
					PrefectureCodes: []int32{13},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShippingRate(tt.rate)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShippingRate_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rate   *ShippingRate
		expect *types.ShippingRate
	}{
		{
			name: "success",
			rate: &ShippingRate{
				ShippingRate: types.ShippingRate{
					Number:          1,
					Name:            "東京都",
					Price:           1200,
					PrefectureCodes: []int32{13},
				},
			},
			expect: &types.ShippingRate{
				Number:          1,
				Name:            "東京都",
				Price:           1200,
				PrefectureCodes: []int32{13},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.rate.Response())
		})
	}
}

func TestShippingRates(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rates  entity.ShippingRates
		expect ShippingRates
	}{
		{
			name: "success",
			rates: entity.ShippingRates{
				{
					Number:          1,
					Name:            "東京都",
					Price:           1200,
					PrefectureCodes: []int32{13},
				},
			},
			expect: ShippingRates{
				{
					ShippingRate: types.ShippingRate{
						Number:          1,
						Name:            "東京都",
						Price:           1200,
						PrefectureCodes: []int32{13},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShippingRates(tt.rates)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShippingRates_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rates  ShippingRates
		expect []*types.ShippingRate
	}{
		{
			name: "success",
			rates: ShippingRates{
				{
					ShippingRate: types.ShippingRate{
						Number:          1,
						Name:            "東京都",
						Price:           1200,
						PrefectureCodes: []int32{13},
					},
				},
			},
			expect: []*types.ShippingRate{
				{
					Number:          1,
					Name:            "東京都",
					Price:           1200,
					PrefectureCodes: []int32{13},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.rates.Response())
		})
	}
}
