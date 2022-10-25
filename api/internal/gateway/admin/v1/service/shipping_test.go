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
		hasErr   bool
	}{
		{
			name: "success",
			shipping: &entity.Shipping{
				ID:   "shipping-id",
				Name: "デフォルト配送設定",
				Box60Rates: entity.ShippingRates{
					{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{13}},
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
						{Number: 1, Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
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
			hasErr: false,
		},
		{
			name: "failed to create box 60 rates",
			shipping: &entity.Shipping{
				ID:   "shipping-id",
				Name: "デフォルト配送設定",
				Box60Rates: entity.ShippingRates{
					{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{0}},
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
			expect: nil,
			hasErr: true,
		},
		{
			name: "failed to create box 80 rates",
			shipping: &entity.Shipping{
				ID:                "shipping-id",
				Name:              "デフォルト配送設定",
				Box60Rates:        entity.ShippingRates{},
				Box60Refrigerated: 500,
				Box60Frozen:       800,
				Box80Rates: entity.ShippingRates{
					{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{0}},
				},
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
			expect: nil,
			hasErr: true,
		},
		{
			name: "failed to create box 100 rates",
			shipping: &entity.Shipping{
				ID:                "shipping-id",
				Name:              "デフォルト配送設定",
				Box60Rates:        entity.ShippingRates{},
				Box60Refrigerated: 500,
				Box60Frozen:       800,
				Box80Rates:        entity.ShippingRates{},
				Box80Refrigerated: 500,
				Box80Frozen:       800,
				Box100Rates: entity.ShippingRates{
					{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{0}},
				},
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewShipping(tt.shipping)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

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
		hasErr    bool
	}{
		{
			name: "success",
			shippings: entity.Shippings{
				{
					ID:   "shipping-id",
					Name: "デフォルト配送設定",
					Box60Rates: entity.ShippingRates{
						{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{13}},
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
							{Number: 1, Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
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
			hasErr: false,
		},
		{
			name: "failed to create",
			shippings: entity.Shippings{
				{
					ID:   "shipping-id",
					Name: "デフォルト配送設定",
					Box60Rates: entity.ShippingRates{
						{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{0}},
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
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewShippings(tt.shippings)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
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

func TestShippingRate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rate   *entity.ShippingRate
		expect *ShippingRate
		hasErr bool
	}{
		{
			name: "success",
			rate: &entity.ShippingRate{
				Number:      1,
				Name:        "東京都",
				Price:       1200,
				Prefectures: []int64{13},
			},
			expect: &ShippingRate{
				ShippingRate: response.ShippingRate{
					Number:      1,
					Name:        "東京都",
					Price:       1200,
					Prefectures: []string{"tokyo"},
				},
			},
			hasErr: false,
		},
		{
			name: "failed to unknown prefecture",
			rate: &entity.ShippingRate{
				Number:      1,
				Name:        "東京都",
				Price:       1200,
				Prefectures: []int64{0},
			},
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewShippingRate(tt.rate)
			assert.Equal(t, tt.hasErr, err != nil, err)
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
					Number:      1,
					Name:        "東京都",
					Price:       1200,
					Prefectures: []string{"tokyo"},
				},
			},
			expect: &response.ShippingRate{
				Number:      1,
				Name:        "東京都",
				Price:       1200,
				Prefectures: []string{"tokyo"},
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
		hasErr bool
	}{
		{
			name: "success",
			rates: entity.ShippingRates{
				{
					Number:      1,
					Name:        "東京都",
					Price:       1200,
					Prefectures: []int64{13},
				},
			},
			expect: ShippingRates{
				{
					ShippingRate: response.ShippingRate{
						Number:      1,
						Name:        "東京都",
						Price:       1200,
						Prefectures: []string{"tokyo"},
					},
				},
			},
			hasErr: false,
		},
		{
			name: "failed to unknown prefecture",
			rates: entity.ShippingRates{
				{
					Number:      1,
					Name:        "東京都",
					Price:       1200,
					Prefectures: []int64{0},
				},
			},
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewShippingRates(tt.rates)
			assert.Equal(t, tt.hasErr, err != nil, err)
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
						Number:      1,
						Name:        "東京都",
						Price:       1200,
						Prefectures: []string{"tokyo"},
					},
				},
			},
			expect: []*response.ShippingRate{
				{
					Number:      1,
					Name:        "東京都",
					Price:       1200,
					Prefectures: []string{"tokyo"},
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
