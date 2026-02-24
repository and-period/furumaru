package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestProductOrderPaymentSummary(t *testing.T) {
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
		params    *NewProductOrderPaymentSummaryParams
		expect    *OrderPaymentSummary
		expectErr error
	}{
		{
			name: "success with shipping free",
			params: &NewProductOrderPaymentSummaryParams{
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
			params: &NewProductOrderPaymentSummaryParams{
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
			name: "success free",
			params: &NewProductOrderPaymentSummaryParams{
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
							Price:     0,
						},
					},
					{
						ID:   "product-id02",
						Name: "人参",
						ProductRevision: ProductRevision{
							ID:        2,
							ProductID: "product-id02",
							Price:     0,
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
					DiscountType: DiscountTypeRate,
					DiscountRate: 10,
					Code:         "excode01",
					CodeType:     PromotionCodeTypeAlways,
					StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			expect: &OrderPaymentSummary{
				Subtotal:    0,
				Discount:    0,
				ShippingFee: 0,
				Tax:         0,
				TaxRate:     10,
				Total:       0,
			},
			expectErr: nil,
		},
		{
			name: "success with pickup",
			params: &NewProductOrderPaymentSummaryParams{
				PrefectureCode: 13,
				Pickup:         true,
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
			name: "failed to calc total price",
			params: &NewProductOrderPaymentSummaryParams{
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
			params: &NewProductOrderPaymentSummaryParams{
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
			expectErr: ErrNotFoundShippingRate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewProductOrderPaymentSummary(tt.params)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceOrderPaymentSummary(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		params *NewExperienceOrderPaymentSummaryParams
		expect *OrderPaymentSummary
	}{
		{
			name: "success with no promotion",
			params: &NewExperienceOrderPaymentSummaryParams{
				Experience: &Experience{
					ID:            "experience-id",
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
					TypeID:        "experience-type-id",
					Title:         "じゃがいも収穫",
					Description:   "じゃがいもを収穫する体験",
					Public:        true,
					SoldOut:       false,
					Status:        ExperienceStatusAccepting,
					ThumbnailURL:  "http://example.com/thumbnail.png",
					Media: MultiExperienceMedia{
						{
							URL:         "http://example.com/thumbnail.png",
							IsThumbnail: true,
						},
					},
					RecommendedPoints: []string{
						"ポイント1",
						"ポイント2",
					},
					PromotionVideoURL:  "http://example.com/promotion.mp4",
					Duration:           60,
					Direction:          "彦根駅から徒歩10分",
					BusinessOpenTime:   "1000",
					BusinessCloseTime:  "1800",
					HostPostalCode:     "5220061",
					HostPrefecture:     "滋賀県",
					HostPrefectureCode: 25,
					HostCity:           "彦根市",
					HostAddressLine1:   "金亀町１−１",
					HostAddressLine2:   "",
					HostLongitude:      136.251739,
					HostLatitude:       35.276833,
					ExperienceRevision: ExperienceRevision{
						ID:                    1,
						ExperienceID:          "experience-id",
						PriceAdult:            1000,
						PriceJuniorHighSchool: 500,
						PriceElementarySchool: 300,
						PricePreschool:        0,
						PriceSenior:           200,
					},
					StartAt:   now.AddDate(0, 0, -1),
					EndAt:     now.AddDate(0, 0, 1),
					CreatedAt: now,
					UpdatedAt: now,
				},
				Promotion:             nil,
				AdultCount:            2,
				JuniorHighSchoolCount: 0,
				ElementarySchoolCount: 0,
				PreschoolCount:        0,
				SeniorCount:           0,
			},
			expect: &OrderPaymentSummary{
				Subtotal:    2000,
				Discount:    0,
				ShippingFee: 0,
				Tax:         181,
				TaxRate:     10,
				Total:       2000,
			},
		},
		{
			name: "success with promotion",
			params: &NewExperienceOrderPaymentSummaryParams{
				Experience: &Experience{
					ID:            "experience-id",
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
					TypeID:        "experience-type-id",
					Title:         "じゃがいも収穫",
					Description:   "じゃがいもを収穫する体験",
					Public:        true,
					SoldOut:       false,
					Status:        ExperienceStatusAccepting,
					ThumbnailURL:  "http://example.com/thumbnail.png",
					Media: MultiExperienceMedia{
						{
							URL:         "http://example.com/thumbnail.png",
							IsThumbnail: true,
						},
					},
					RecommendedPoints: []string{
						"ポイント1",
						"ポイント2",
					},
					PromotionVideoURL:  "http://example.com/promotion.mp4",
					Duration:           60,
					Direction:          "彦根駅から徒歩10分",
					BusinessOpenTime:   "1000",
					BusinessCloseTime:  "1800",
					HostPostalCode:     "5220061",
					HostPrefecture:     "滋賀県",
					HostPrefectureCode: 25,
					HostCity:           "彦根市",
					HostAddressLine1:   "金亀町１−１",
					HostAddressLine2:   "",
					HostLongitude:      136.251739,
					HostLatitude:       35.276833,
					ExperienceRevision: ExperienceRevision{
						ID:                    1,
						ExperienceID:          "experience-id",
						PriceAdult:            1000,
						PriceJuniorHighSchool: 500,
						PriceElementarySchool: 300,
						PricePreschool:        0,
						PriceSenior:           200,
					},
					StartAt:   now.AddDate(0, 0, -1),
					EndAt:     now.AddDate(0, 0, 1),
					CreatedAt: now,
					UpdatedAt: now,
				},
				Promotion: &Promotion{
					Title:        "プロモーションタイトル",
					Description:  "プロモーションの詳細です。",
					Public:       true,
					DiscountType: DiscountTypeAmount,
					DiscountRate: 500,
					Code:         "excode02",
					CodeType:     PromotionCodeTypeAlways,
					StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
				AdultCount:            1,
				JuniorHighSchoolCount: 2,
				ElementarySchoolCount: 3,
				PreschoolCount:        0,
				SeniorCount:           0,
			},
			expect: &OrderPaymentSummary{
				Subtotal:    2900,
				Discount:    500,
				ShippingFee: 0,
				Tax:         218,
				TaxRate:     10,
				Total:       2400,
			},
		},
		{
			name: "success free",
			params: &NewExperienceOrderPaymentSummaryParams{
				Experience: &Experience{
					ID:            "experience-id",
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
					TypeID:        "experience-type-id",
					Title:         "じゃがいも収穫",
					Description:   "じゃがいもを収穫する体験",
					Public:        true,
					SoldOut:       false,
					Status:        ExperienceStatusAccepting,
					ThumbnailURL:  "http://example.com/thumbnail.png",
					Media: MultiExperienceMedia{
						{
							URL:         "http://example.com/thumbnail.png",
							IsThumbnail: true,
						},
					},
					RecommendedPoints: []string{
						"ポイント1",
						"ポイント2",
					},
					PromotionVideoURL:  "http://example.com/promotion.mp4",
					Duration:           60,
					Direction:          "彦根駅から徒歩10分",
					BusinessOpenTime:   "1000",
					BusinessCloseTime:  "1800",
					HostPostalCode:     "5220061",
					HostPrefecture:     "滋賀県",
					HostPrefectureCode: 25,
					HostCity:           "彦根市",
					HostAddressLine1:   "金亀町１−１",
					HostAddressLine2:   "",
					HostLongitude:      136.251739,
					HostLatitude:       35.276833,
					ExperienceRevision: ExperienceRevision{
						ID:                    1,
						ExperienceID:          "experience-id",
						PriceAdult:            0,
						PriceJuniorHighSchool: 0,
						PriceElementarySchool: 0,
						PricePreschool:        0,
						PriceSenior:           0,
					},
					StartAt:   now.AddDate(0, 0, -1),
					EndAt:     now.AddDate(0, 0, 1),
					CreatedAt: now,
					UpdatedAt: now,
				},
				Promotion: &Promotion{
					Title:        "プロモーションタイトル",
					Description:  "プロモーションの詳細です。",
					Public:       true,
					DiscountType: DiscountTypeAmount,
					DiscountRate: 500,
					Code:         "excode02",
					CodeType:     PromotionCodeTypeAlways,
					StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
				AdultCount:            1,
				JuniorHighSchoolCount: 2,
				ElementarySchoolCount: 3,
				PreschoolCount:        0,
				SeniorCount:           0,
			},
			expect: &OrderPaymentSummary{
				Subtotal:    0,
				Discount:    0,
				ShippingFee: 0,
				Tax:         0,
				TaxRate:     10,
				Total:       0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewExperienceOrderPaymentSummary(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
