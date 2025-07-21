package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestShippingType_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		shippingType ShippingType
		expect       string
	}{
		{
			name:         "normal",
			shippingType: ShippingTypeNormal,
			expect:       "通常配送",
		},
		{
			name:         "frozen",
			shippingType: ShippingTypeFrozen,
			expect:       "クール配送",
		},
		{
			name:         "unknown",
			shippingType: ShippingTypeUnknown,
			expect:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shippingType.String())
		})
	}
}

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
				ShopID:            "shop-id",
				CoordinatorID:     "coordinator-id",
				Name:              "配送設定",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expect: &Shipping{
				ID:            "coordinator-id",
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				Name:          "配送設定",
				ShippingRevision: ShippingRevision{
					ShippingID:        "coordinator-id",
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
		},
	}
	for _, tt := range tests {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.shipping.IsDefault())
		})
	}
}

func TestShipping_CalcShippingFee(t *testing.T) {
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
		name           string
		shipping       *Shipping
		shippingSize   ShippingSize
		shippingType   ShippingType
		total          int64
		prefectureCode int32
		expect         int64
		expectErr      error
	}{
		{
			name: "success free shipping",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			shippingSize:   ShippingSize60,
			shippingType:   ShippingTypeNormal,
			total:          3000,
			prefectureCode: codes.PrefectureValues["kochi"],
			expect:         0,
			expectErr:      nil,
		},
		{
			name: "success normal box 60",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			shippingSize:   ShippingSize60,
			shippingType:   ShippingTypeNormal,
			total:          2980,
			prefectureCode: codes.PrefectureValues["kochi"],
			expect:         500,
			expectErr:      nil,
		},
		{
			name: "success frozen box 60",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			shippingSize:   ShippingSize60,
			shippingType:   ShippingTypeFrozen,
			total:          2980,
			prefectureCode: codes.PrefectureValues["kochi"],
			expect:         1300,
			expectErr:      nil,
		},
		{
			name: "success normal box 80",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			shippingSize:   ShippingSize80,
			shippingType:   ShippingTypeNormal,
			total:          2980,
			prefectureCode: codes.PrefectureValues["kochi"],
			expect:         500,
			expectErr:      nil,
		},
		{
			name: "success frozen box 80",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			shippingSize:   ShippingSize80,
			shippingType:   ShippingTypeFrozen,
			total:          2980,
			prefectureCode: codes.PrefectureValues["kochi"],
			expect:         1300,
			expectErr:      nil,
		},
		{
			name: "success normal box 100",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			shippingSize:   ShippingSize100,
			shippingType:   ShippingTypeNormal,
			total:          2980,
			prefectureCode: codes.PrefectureValues["kochi"],
			expect:         500,
			expectErr:      nil,
		},
		{
			name: "success frozen box 100",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			shippingSize:   ShippingSize100,
			shippingType:   ShippingTypeFrozen,
			total:          2980,
			prefectureCode: codes.PrefectureValues["kochi"],
			expect:         1300,
			expectErr:      nil,
		},
		{
			name: "unknown shipping size",
			shipping: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
			shippingSize:   ShippingSizeUnknown,
			shippingType:   ShippingTypeUnknown,
			total:          2980,
			prefectureCode: codes.PrefectureValues["kochi"],
			expect:         0,
			expectErr:      errUnknownShippingSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.shipping.CalcShippingFee(
				tt.shippingSize,
				tt.shippingType,
				tt.total,
				tt.prefectureCode,
			)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
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
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expect: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
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
				Box60Rates:        nil,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expect: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
					Box60Rates:        nil,
					Box60Frozen:       800,
					Box80Rates:        rates,
					Box80Frozen:       800,
					Box100Rates:       rates,
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
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
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        nil,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expect: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
					Box60Rates:        rates,
					Box60Frozen:       800,
					Box80Rates:        nil,
					Box80Frozen:       800,
					Box100Rates:       rates,
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
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
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       nil,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expect: &Shipping{
				ID:        "shipping-id",
				CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
				ShippingRevision: ShippingRevision{
					Box60Rates:        rates,
					Box60Frozen:       800,
					Box80Rates:        rates,
					Box80Frozen:       800,
					Box100Rates:       nil,
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
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
			},
			expect: []string{"shipping-id"},
		},
	}
	for _, tt := range tests {
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
			},
			expect: []string{"coordinator-id"},
		},
	}
	for _, tt := range tests {
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
			},
			expect: map[string]*Shipping{
				"shipping-id": {
					ID:            "shipping-id",
					CoordinatorID: "coordinator-id",
					ShippingRevision: ShippingRevision{
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
			},
		},
	}
	for _, tt := range tests {
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
					ID:                1,
					ShippingID:        "shipping-id01",
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
			expect: Shippings{
				{
					ID:        "shipping-id01",
					CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					ShippingRevision: ShippingRevision{
						ID:                1,
						ShippingID:        "shipping-id01",
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
				{
					ID:        "shipping-id02",
					CreatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt: jst.Date(2022, 7, 3, 18, 30, 0, 0),
					ShippingRevision: ShippingRevision{
						ShippingID: "shipping-id02",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.shippings.Fill(tt.revisions)
			assert.Equal(t, tt.expect, tt.shippings)
		})
	}
}
