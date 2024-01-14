package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestOrderPaymentSummary(t *testing.T) {
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
		params    *NewOrderPaymentSummaryParams
		expect    *OrderPaymentSummary
		expectErr error
	}{
		{
			name: "success with shipping free",
			params: &NewOrderPaymentSummaryParams{
				PrefectureCode: 13,
				Baskets: CartBaskets{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: []*CartItem{
							{
								ProductID: "product-id01",
								Quantity:  1,
							},
							{
								ProductID: "product-id02",
								Quantity:  2,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				Products: []*Product{
					{
						ID:   "product-id01",
						Name: "じゃがいも",
						ProductRevision: ProductRevision{
							ID:        1,
							ProductID: "product-id01",
							Price:     500,
						},
					},
					{
						ID:   "product-id02",
						Name: "人参",
						ProductRevision: ProductRevision{
							ID:        2,
							ProductID: "product-id02",
							Price:     1980,
						},
					},
				},
				Shipping: &Shipping{
					ID:            "coordinator-id",
					CoordinatorID: "coordinator-id",
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
				Promotion: nil,
			},
			expect: &OrderPaymentSummary{
				Subtotal:    4460,
				Discount:    0,
				ShippingFee: 0,
				Tax:         405,
				TaxRate:     10,
				Total:       4460,
			},
			expectErr: nil,
		},
		{
			name: "success with discount",
			params: &NewOrderPaymentSummaryParams{
				PrefectureCode: 13,
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: []*CartItem{
							{
								ProductID: "product-id01",
								Quantity:  1,
							},
							{
								ProductID: "product-id02",
								Quantity:  2,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				Products: []*Product{
					{
						ID:   "product-id01",
						Name: "じゃがいも",
						ProductRevision: ProductRevision{
							ID:        1,
							ProductID: "product-id01",
							Price:     500,
						},
					},
					{
						ID:   "product-id02",
						Name: "人参",
						ProductRevision: ProductRevision{
							ID:        2,
							ProductID: "product-id02",
							Price:     1980,
						},
					},
				},
				Shipping: &Shipping{
					ID:            "coordinator-id",
					CoordinatorID: "coordinator-id",
					ShippingRevision: ShippingRevision{
						ShippingID:      "coordinator-id",
						Box60Rates:      rates,
						Box60Frozen:     800,
						Box80Rates:      rates,
						Box80Frozen:     800,
						Box100Rates:     rates,
						Box100Frozen:    800,
						HasFreeShipping: false,
					},
				},
				Promotion: &Promotion{
					Title:        "プロモーションタイトル",
					Description:  "プロモーションの詳細です。",
					Public:       true,
					PublishedAt:  jst.Date(2022, 8, 9, 18, 30, 0, 0),
					DiscountType: DiscountTypeRate,
					DiscountRate: 10,
					Code:         "excode01",
					CodeType:     PromotionCodeTypeAlways,
					StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			expect: &OrderPaymentSummary{
				Subtotal:    4460,
				Discount:    446,
				ShippingFee: 500,
				Tax:         410,
				TaxRate:     10,
				Total:       4514,
			},
			expectErr: nil,
		},
		{
			name: "failed to calc total price",
			params: &NewOrderPaymentSummaryParams{
				PrefectureCode: 13,
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: []*CartItem{
							{
								ProductID: "product-id01",
								Quantity:  1,
							},
							{
								ProductID: "product-id02",
								Quantity:  2,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				Products:  []*Product{},
				Shipping:  nil,
				Promotion: nil,
			},
			expect:    nil,
			expectErr: errNotFoundProduct,
		},
		{
			name: "failed to calc shipping fee",
			params: &NewOrderPaymentSummaryParams{
				PrefectureCode: 13,
				Baskets: []*CartBasket{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						Items: []*CartItem{
							{
								ProductID: "product-id01",
								Quantity:  1,
							},
							{
								ProductID: "product-id02",
								Quantity:  2,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				Products: []*Product{
					{
						ID:   "product-id01",
						Name: "じゃがいも",
						ProductRevision: ProductRevision{
							ID:        1,
							ProductID: "product-id01",
							Price:     500,
						},
					},
					{
						ID:   "product-id02",
						Name: "人参",
						ProductRevision: ProductRevision{
							ID:        2,
							ProductID: "product-id02",
							Price:     1980,
						},
					},
				},
				Shipping: &Shipping{
					ID:            "coordinator-id",
					CoordinatorID: "coordinator-id",
					ShippingRevision: ShippingRevision{
						ShippingID:      "coordinator-id",
						HasFreeShipping: false,
					},
				},
				Promotion: nil,
			},
			expect:    nil,
			expectErr: errNotFoundShippingRate,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewOrderPaymentSummary(tt.params)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
