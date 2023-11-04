package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
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
				ID:   "shipping-id",
				Name: "デフォルト配送設定",
				Box60Rates: entity.ShippingRates{
					{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
				},
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         entity.ShippingRates{},
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        entity.ShippingRates{},
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Shipping{
				response.Shipping{
					ID:   "shipping-id",
					Name: "デフォルト配送設定",
					Box60Rates: []*response.ShippingRate{
						{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
					},
					Box60Refrigerated:  500,
					Box60Frozen:        800,
					Box80Rates:         []*response.ShippingRate{},
					Box80Refrigerated:  500,
					Box80Frozen:        800,
					Box100Rates:        []*response.ShippingRate{},
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
			actual := NewShipping(tt.shipping)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShipping_Response(t *testing.T) {
	t.Parallel()
	rates := []*response.ShippingRate{
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

func TestShippings_Map(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		shippings Shippings
		expect    map[string]*Shipping
	}{
		{
			name: "success",
			shippings: Shippings{
				{
					Shipping: response.Shipping{
						ID:        "shipping-id",
						Name:      "デフォルト配送設定",
						CreatedAt: 1640962800,
						UpdatedAt: 1640962800,
					},
				},
			},
			expect: map[string]*Shipping{
				"shipping-id": {
					Shipping: response.Shipping{
						ID:        "shipping-id",
						Name:      "デフォルト配送設定",
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
			assert.Equal(t, tt.expect, tt.shippings.Map())
		})
	}
}

func TestShippings(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		shippings entity.Shippings
		expect    Shippings
	}{
		{
			name: "success",
			shippings: entity.Shippings{
				{
					ID:   "shipping-id",
					Name: "デフォルト配送設定",
					Box60Rates: entity.ShippingRates{
						{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
					},
					Box60Refrigerated:  500,
					Box60Frozen:        800,
					Box80Rates:         entity.ShippingRates{},
					Box80Refrigerated:  500,
					Box80Frozen:        800,
					Box100Rates:        entity.ShippingRates{},
					Box100Refrigerated: 500,
					Box100Frozen:       800,
					HasFreeShipping:    true,
					FreeShippingRates:  3000,
					CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Shippings{
				{
					response.Shipping{
						ID:   "shipping-id",
						Name: "デフォルト配送設定",
						Box60Rates: []*response.ShippingRate{
							{Number: 1, Name: "東京都", Price: 0, PrefectureCodes: []int32{13}},
						},
						Box60Refrigerated:  500,
						Box60Frozen:        800,
						Box80Rates:         []*response.ShippingRate{},
						Box80Refrigerated:  500,
						Box80Frozen:        800,
						Box100Rates:        []*response.ShippingRate{},
						Box100Refrigerated: 500,
						Box100Frozen:       800,
						HasFreeShipping:    true,
						FreeShippingRates:  3000,
						CreatedAt:          1640962800,
						UpdatedAt:          1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShippings(tt.shippings)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShippings_Response(t *testing.T) {
	t.Parallel()
	rates := []*response.ShippingRate{
		{
			Number:          1,
			Name:            "四国",
			Price:           250,
			PrefectureCodes: []int32{36, 37, 38, 39},
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
				ShippingRate: response.ShippingRate{
					Number:          1,
					Name:            "東京都",
					Price:           1200,
					PrefectureCodes: []int32{13},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
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
		expect *response.ShippingRate
	}{
		{
			name: "success",
			rate: &ShippingRate{
				ShippingRate: response.ShippingRate{
					Number:          1,
					Name:            "東京都",
					Price:           1200,
					PrefectureCodes: []int32{13},
				},
			},
			expect: &response.ShippingRate{
				Number:          1,
				Name:            "東京都",
				Price:           1200,
				PrefectureCodes: []int32{13},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
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
					ShippingRate: response.ShippingRate{
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
		tt := tt
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
		expect []*response.ShippingRate
	}{
		{
			name: "success",
			rates: ShippingRates{
				{
					ShippingRate: response.ShippingRate{
						Number:          1,
						Name:            "東京都",
						Price:           1200,
						PrefectureCodes: []int32{13},
					},
				},
			},
			expect: []*response.ShippingRate{
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.rates.Response())
		})
	}
}
