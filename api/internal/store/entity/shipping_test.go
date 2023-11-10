package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestShipping(t *testing.T) {
	t.Parallel()
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for _, val := range codes.PrefectureValues {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	tests := []struct {
		name   string
		params *NewShippingParams
		expect *Shipping
	}{
		{
			name: "success",
			params: &NewShippingParams{
				CoordinatorID:      "coordinator-id",
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
				ID:            "coordinator-id",
				CoordinatorID: "coordinator-id",
				ShippingRevision: ShippingRevision{
					ShippingID:         "coordinator-id",
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
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShipping(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShipping_IsDefault(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		shipping *Shipping
		expect   bool
	}{
		{
			name: "default",
			shipping: &Shipping{
				ID:            "default",
				CoordinatorID: "",
			},
			expect: true,
		},
		{
			name: "not default",
			shipping: &Shipping{
				ID:            "shipping-id",
				CoordinatorID: "coordinator-id",
			},
			expect: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shipping.IsDefault())
		})
	}
}

func TestShipping_Fill(t *testing.T) {
	t.Parallel()
	pref1 := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int32{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: pref1},
		{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: pref2},
	}
	tests := []struct {
		name     string
		shipping *Shipping
		revision *ShippingRevision
		expect   *Shipping
		hasErr   bool
	}{
		{
			name: "success",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			revision: &ShippingRevision{
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
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			hasErr: false,
		},
		{
			name: "success Box60Rates is nil",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			revision: &ShippingRevision{
				Box60Rates:         nil,
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
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
					Box60Rates:         nil,
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
			hasErr: false,
		},
		{
			name: "success Box80Rates is nil",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			revision: &ShippingRevision{
				Box60Rates:         rates,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         nil,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        rates,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
			},
			expect: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
					Box60Rates:         rates,
					Box60Refrigerated:  500,
					Box60Frozen:        800,
					Box80Rates:         nil,
					Box80Refrigerated:  500,
					Box80Frozen:        800,
					Box100Rates:        rates,
					Box100Refrigerated: 500,
					Box100Frozen:       800,
					HasFreeShipping:    true,
					FreeShippingRates:  3000,
				},
			},
			hasErr: false,
		},
		{
			name: "success Box100Rates is nil",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
			},
			revision: &ShippingRevision{
				Box60Rates:         rates,
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         rates,
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        nil,
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
			},
			expect: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
					Box60Rates:         rates,
					Box60Refrigerated:  500,
					Box60Frozen:        800,
					Box80Rates:         rates,
					Box80Refrigerated:  500,
					Box80Frozen:        800,
					Box100Rates:        nil,
					Box100Refrigerated: 500,
					Box100Frozen:       800,
					HasFreeShipping:    true,
					FreeShippingRates:  3000,
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.shipping.Fill(tt.revision)
			assert.Equal(t, tt.expect, tt.shipping)
		})
	}
}

func TestShippings_IDs(t *testing.T) {
	t.Parallel()
	pref1 := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int32{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: pref1},
		{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: pref2},
	}
	tests := []struct {
		name      string
		shippings Shippings
		expect    []string
	}{
		{
			name: "success",
			shippings: Shippings{
				{
					ID:            "shipping-id",
					CoordinatorID: "coordinator-id",
					ShippingRevision: ShippingRevision{
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
			},
			expect: []string{"shipping-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shippings.IDs())
		})
	}
}

func TestShippings_CoordinatorIDs(t *testing.T) {
	t.Parallel()
	pref1 := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int32{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: pref1},
		{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: pref2},
	}
	tests := []struct {
		name      string
		shippings Shippings
		expect    []string
	}{
		{
			name: "success",
			shippings: Shippings{
				{
					ID:            "shipping-id",
					CoordinatorID: "coordinator-id",
					ShippingRevision: ShippingRevision{
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
			},
			expect: []string{"coordinator-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shippings.CoordinatorIDs())
		})
	}
}

func TestShippings_Map(t *testing.T) {
	t.Parallel()
	pref1 := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int32{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: pref1},
		{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: pref2},
	}
	tests := []struct {
		name      string
		shippings Shippings
		expect    map[string]*Shipping
	}{
		{
			name: "success",
			shippings: Shippings{
				{
					ID:            "shipping-id",
					CoordinatorID: "coordinator-id",
					ShippingRevision: ShippingRevision{
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
			},
			expect: map[string]*Shipping{
				"shipping-id": {
					ID:            "shipping-id",
					CoordinatorID: "coordinator-id",
					ShippingRevision: ShippingRevision{
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

func TestShippings_Fill(t *testing.T) {
	t.Parallel()
	pref1 := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int32{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	rates := ShippingRates{
		{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: pref1},
		{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: pref2},
	}
	tests := []struct {
		name      string
		shippings Shippings
		revisions map[string]*ShippingRevision
		expect    Shippings
	}{
		{
			name: "success",
			shippings: Shippings{
				{
					ID:        "shipping-id01",
					CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				},
				{
					ID:        "shipping-id02",
					CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				},
			},
			revisions: map[string]*ShippingRevision{
				"shipping-id01": {
					ID:                 1,
					ShippingID:         "shipping-id01",
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
			expect: Shippings{
				{
					ID:        "shipping-id01",
					CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					ShippingRevision: ShippingRevision{
						ID:                 1,
						ShippingID:         "shipping-id01",
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
				{
					ID:        "shipping-id02",
					CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.shippings.Fill(tt.revisions)
			assert.Equal(t, tt.expect, tt.shippings)
		})
	}
}
