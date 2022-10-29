package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/store/codes"
	"github.com/and-period/furumaru/api/pkg/jst"
	set "github.com/and-period/furumaru/api/pkg/set/v2"
	"github.com/stretchr/testify/assert"
)

func TestShipping(t *testing.T) {
	t.Parallel()
	shikoku := []int64{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int64, 0, 47-len(shikoku))
	for _, val := range codes.PrefectureValues {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国", Price: 250, Prefectures: shikoku},
		{Number: 2, Name: "その他", Price: 500, Prefectures: others},
	}
	tests := []struct {
		name   string
		params *NewShippingParams
		expect *Shipping
	}{
		{
			name: "success",
			params: &NewShippingParams{
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
			},
			expect: &Shipping{
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
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShipping(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShipping_Fill(t *testing.T) {
	t.Parallel()
	pref1 := []int64{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int64{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: pref1},
		{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: pref2},
	}
	buf := []byte(`[{"number":1,"name":"四国(東部)","price":250,"prefectures":[36,37]},{"number":2,"name":"四国(西部)","price":500,"prefectures":[38,39]}]`)
	tests := []struct {
		name     string
		shipping *Shipping
		expect   *Shipping
		hasErr   bool
	}{
		{
			name: "success",
			shipping: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60RatesJSON:     buf,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80RatesJSON:     buf,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100RatesJSON:    buf,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			expect: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60Rates:         rates,
				Box60RatesJSON:     buf,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         rates,
				Box80RatesJSON:     buf,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        rates,
				Box100RatesJSON:    buf,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			hasErr: false,
		},
		{
			name: "success Box60Rates is nil",
			shipping: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60RatesJSON:     nil,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80RatesJSON:     buf,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100RatesJSON:    buf,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			expect: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60Rates:         ShippingRates{},
				Box60RatesJSON:     nil,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         rates,
				Box80RatesJSON:     buf,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        rates,
				Box100RatesJSON:    buf,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			hasErr: false,
		},
		{
			name: "success Box80Rates is nil",
			shipping: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60RatesJSON:     buf,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80RatesJSON:     nil,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100RatesJSON:    buf,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			expect: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60Rates:         rates,
				Box60RatesJSON:     buf,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         ShippingRates{},
				Box80RatesJSON:     nil,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        rates,
				Box100RatesJSON:    buf,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			hasErr: false,
		},
		{
			name: "success Box100Rates is nil",
			shipping: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60RatesJSON:     buf,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80RatesJSON:     buf,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100RatesJSON:    nil,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			expect: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60Rates:         rates,
				Box60RatesJSON:     buf,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         rates,
				Box80RatesJSON:     buf,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        ShippingRates{},
				Box100RatesJSON:    nil,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.shipping.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.shipping)
		})
	}
}

func TestShipping_FillJSON(t *testing.T) {
	t.Parallel()
	pref1 := []int64{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int64{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: pref1},
		{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: pref2},
	}
	buf := []byte(`[{"number":1,"name":"四国(東部)","price":250,"prefectures":[36,37]},{"number":2,"name":"四国(西部)","price":500,"prefectures":[38,39]}]`)
	tests := []struct {
		name     string
		shipping *Shipping
		expect   *Shipping
		hasErr   bool
	}{
		{
			name: "success",
			shipping: &Shipping{
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
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			expect: &Shipping{
				ID:                 "shipping-id",
				Name:               "デフォルト配送設定",
				Box60Rates:         rates,
				Box60RatesJSON:     buf,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         rates,
				Box80RatesJSON:     buf,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        rates,
				Box100RatesJSON:    buf,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
				CreatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt:          jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.shipping.FillJSON()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.shipping)
		})
	}
}

func TestShippingRate(t *testing.T) {
	t.Parallel()
	type input struct {
		num   int64
		name  string
		price int64
		prefs []int64
	}
	tests := []struct {
		name   string
		input  input
		expect *ShippingRate
	}{
		{
			name: "success",
			input: input{
				num:   1,
				name:  "四国",
				price: 2000,
				prefs: []int64{
					codes.PrefectureValues["tokushima"],
					codes.PrefectureValues["kagawa"],
					codes.PrefectureValues["ehime"],
					codes.PrefectureValues["kochi"],
				},
			},
			expect: &ShippingRate{
				Number: 1,
				Name:   "四国",
				Price:  2000,
				Prefectures: []int64{
					codes.PrefectureValues["tokushima"],
					codes.PrefectureValues["kagawa"],
					codes.PrefectureValues["ehime"],
					codes.PrefectureValues["kochi"],
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShippingRate(tt.input.num, tt.input.name, tt.input.price, tt.input.prefs)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShippingRates(t *testing.T) {
	t.Parallel()
	shikoku := []int64{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int64, 0, 47-len(shikoku))
	for _, val := range codes.PrefectureValues {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	tests := []struct {
		name   string
		rates  ShippingRates
		hasErr bool
	}{
		{
			name: "success",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: shikoku},
				{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: others},
			},
			hasErr: false,
		},
		{
			name: "number is checked min",
			rates: ShippingRates{
				{Number: 0, Name: "四国(東部)", Price: 250, Prefectures: shikoku},
				{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: others},
			},
			hasErr: true,
		},
		{
			name: "price is checked min",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: -1, Prefectures: shikoku},
				{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: others},
			},
			hasErr: true,
		},
		{
			name: "number is not unique",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: shikoku},
				{Number: 1, Name: "四国(西部)", Price: 500, Prefectures: others},
			},
			hasErr: true,
		},
		{
			name: "not exists prefecture values",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: shikoku},
				{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: []int64{0}},
			},
			hasErr: true,
		},
		{
			name: "prefectures is umnatch length",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: shikoku},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.rates.Validate()
			assert.Equal(t, tt.hasErr, tt.rates.Validate() != nil, err)
		})
	}
}

func TestShippingRates_Marshal(t *testing.T) {
	t.Parallel()
	pref1 := []int64{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int64{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	tests := []struct {
		name   string
		rates  ShippingRates
		expect []byte
		hasErr bool
	}{
		{
			name: "success",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: pref1},
				{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: pref2},
			},
			expect: []byte(`[{"number":1,"name":"四国(東部)","price":250,"prefectures":[36,37]},{"number":2,"name":"四国(西部)","price":500,"prefectures":[38,39]}]`),
			hasErr: false,
		},
		{
			name:   "shipping rate is empty",
			rates:  nil,
			expect: []byte{},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.rates.Marshal()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
