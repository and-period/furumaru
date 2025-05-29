package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestShippingRevision(t *testing.T) {
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
		params *NewShippingRevisionParams
		expect *ShippingRevision
	}{
		{
			name: "success",
			params: &NewShippingRevisionParams{
				ShippingID:        "shipping-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			expect: &ShippingRevision{
				ShippingID:        "shipping-id",
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
			actual := NewShippingRevision(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShippingRevision_Fill(t *testing.T) {
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
	additional := ShippingRates{
		{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: []string{"徳島県", "香川県"}, PrefectureCodes: pref1},
		{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: []string{"愛媛県", "高知県"}, PrefectureCodes: pref2},
	}
	tests := []struct {
		name   string
		params *ShippingRevision
		expect *ShippingRevision
	}{
		{
			name: "success",
			params: &ShippingRevision{
				Box60Rates:  rates,
				Box80Rates:  rates,
				Box100Rates: rates,
			},
			expect: &ShippingRevision{
				Box60Rates:  additional,
				Box80Rates:  additional,
				Box100Rates: additional,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.params.Fill()
			assert.Equal(t, tt.expect, tt.params)
		})
	}
}

func TestShippingRevisions_ShippingIDs(t *testing.T) {
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
		revisions ShippingRevisions
		expect    []string
	}{
		{
			name: "success",
			revisions: ShippingRevisions{
				{
					ID:                1,
					ShippingID:        "shipping-id",
					Box60Rates:        rates,
					Box60Frozen:       800,
					Box80Rates:        rates,
					Box80Frozen:       800,
					Box100Rates:       rates,
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
					CreatedAt:         jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt:         jst.Date(2022, 7, 3, 18, 30, 0, 0),
				},
			},
			expect: []string{"shipping-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.revisions.ShippingIDs())
		})
	}
}

func TestShippingRevisions_MapByShippingID(t *testing.T) {
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
		revisions ShippingRevisions
		expect    map[string]*ShippingRevision
	}{
		{
			name: "success",
			revisions: ShippingRevisions{
				{
					ID:                1,
					ShippingID:        "shipping-id",
					Box60Rates:        rates,
					Box60Frozen:       800,
					Box80Rates:        rates,
					Box80Frozen:       800,
					Box100Rates:       rates,
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
					CreatedAt:         jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt:         jst.Date(2022, 7, 3, 18, 30, 0, 0),
				},
			},
			expect: map[string]*ShippingRevision{
				"shipping-id": {
					ID:                1,
					ShippingID:        "shipping-id",
					Box60Rates:        rates,
					Box60Frozen:       800,
					Box80Rates:        rates,
					Box80Frozen:       800,
					Box100Rates:       rates,
					Box100Frozen:      800,
					HasFreeShipping:   true,
					FreeShippingRates: 3000,
					CreatedAt:         jst.Date(2022, 7, 3, 18, 30, 0, 0),
					UpdatedAt:         jst.Date(2022, 7, 3, 18, 30, 0, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.revisions.MapByShippingID())
		})
	}
}

func TestShippingRevisions_Merge(t *testing.T) {
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
		name      string
		revisions ShippingRevisions
		shippings map[string]*Shipping
		expect    Shippings
		hasErr    bool
	}{
		{
			name: "success",
			revisions: ShippingRevisions{
				{
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
				{
					ID:                2,
					ShippingID:        "shipping-id02",
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
			shippings: map[string]*Shipping{
				"shipping-id01": {
					ID:            "shipping-id01",
					CoordinatorID: "coordinator-id",
				},
			},
			expect: Shippings{
				{
					ID:            "shipping-id01",
					CoordinatorID: "coordinator-id",
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
					ID: "shipping-id02",
					ShippingRevision: ShippingRevision{
						ID:                2,
						ShippingID:        "shipping-id02",
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
			hasErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.revisions.Merge(tt.shippings)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShippingRate(t *testing.T) {
	t.Parallel()
	type input struct {
		num   int64
		name  string
		price int64
		prefs []int32
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
				prefs: []int32{
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
				PrefectureCodes: []int32{
					codes.PrefectureValues["tokushima"],
					codes.PrefectureValues["kagawa"],
					codes.PrefectureValues["ehime"],
					codes.PrefectureValues["kochi"],
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewShippingRate(tt.input.num, tt.input.name, tt.input.price, tt.input.prefs)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShippingRates_Find(t *testing.T) {
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
	tests := []struct {
		name           string
		rates          ShippingRates
		prefectureCode int32
		expect         *ShippingRate
		expectErr      error
	}{
		{
			name: "success",
			rates: ShippingRates{
				{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
				{Number: 2, Name: "四国(以外)", Price: 500, PrefectureCodes: others},
			},
			prefectureCode: codes.PrefectureValues["kagawa"],
			expect: &ShippingRate{
				Number:          1,
				Name:            "四国",
				Price:           250,
				PrefectureCodes: shikoku,
			},
			expectErr: nil,
		},
		{
			name:           "not found",
			rates:          nil,
			prefectureCode: codes.PrefectureValues["kagawa"],
			expect:         nil,
			expectErr:      errNotFoundShippingRate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.rates.Find(tt.prefectureCode)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestShippingRates_Validate(t *testing.T) {
	t.Parallel()
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for val := range codes.PrefectureNames {
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
				{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: shikoku},
				{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: others},
			},
			hasErr: false,
		},
		{
			name: "number is checked min",
			rates: ShippingRates{
				{Number: 0, Name: "四国(東部)", Price: 250, PrefectureCodes: shikoku},
				{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: others},
			},
			hasErr: true,
		},
		{
			name: "price is checked min",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: -1, PrefectureCodes: shikoku},
				{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: others},
			},
			hasErr: true,
		},
		{
			name: "number is not unique",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: shikoku},
				{Number: 1, Name: "四国(西部)", Price: 500, PrefectureCodes: others},
			},
			hasErr: true,
		},
		{
			name: "not exists prefecture values",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: shikoku},
				{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: []int32{0}},
			},
			hasErr: true,
		},
		{
			name: "prefectures is umnatch length",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: shikoku},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.rates.Validate()
			assert.Equal(t, tt.hasErr, tt.rates.Validate() != nil, err)
		})
	}
}

func TestShippingRates_Fill(t *testing.T) {
	t.Parallel()
	pref1 := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int32{
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	tests := []struct {
		name   string
		rates  ShippingRates
		expect ShippingRates
	}{
		{
			name: "success",
			rates: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: pref1},
				{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: pref2},
			},
			expect: ShippingRates{
				{Number: 1, Name: "四国(東部)", Price: 250, Prefectures: []string{"徳島県", "香川県"}, PrefectureCodes: pref1},
				{Number: 2, Name: "四国(西部)", Price: 500, Prefectures: []string{"愛媛県", "高知県"}, PrefectureCodes: pref2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.rates.Fill()
			assert.Equal(t, tt.expect, tt.rates)
		})
	}
}

func TestShippingRates_Marshal(t *testing.T) {
	t.Parallel()
	pref1 := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
	}
	pref2 := []int32{
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
				{Number: 1, Name: "四国(東部)", Price: 250, PrefectureCodes: pref1},
				{Number: 2, Name: "四国(西部)", Price: 500, PrefectureCodes: pref2},
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.rates.Marshal()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
