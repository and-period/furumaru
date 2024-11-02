package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestPaymentStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status komoju.PaymentStatus
		expect PaymentStatus
	}{
		{
			name:   "pending",
			status: komoju.PaymentStatusPending,
			expect: PaymentStatusPending,
		},
		{
			name:   "authorized",
			status: komoju.PaymentStatusAuthorized,
			expect: PaymentStatusAuthorized,
		},
		{
			name:   "captured",
			status: komoju.PaymentStatusCaptured,
			expect: PaymentStatusCaptured,
		},
		{
			name:   "refunded",
			status: komoju.PaymentStatusRefunded,
			expect: PaymentStatusRefunded,
		},
		{
			name:   "cancelled",
			status: komoju.PaymentStatusCancelled,
			expect: PaymentStatusCanceled,
		},
		{
			name:   "expired",
			status: komoju.PaymentStatusExpired,
			expect: PaymentStatusFailed,
		},
		{
			name:   "failed",
			status: komoju.PaymentStatusFailed,
			expect: PaymentStatusFailed,
		},
		{
			name:   "failed",
			status: "",
			expect: PaymentStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewPaymentStatus(tt.status))
		})
	}
}

func TestKomojuPaymentTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		methodType PaymentMethodType
		expect     []komoju.PaymentType
	}{
		{
			name:       "success",
			methodType: PaymentMethodTypeCash,
			expect:     []komoju.PaymentType{},
		},
		{
			name:       "credit card",
			methodType: PaymentMethodTypeCreditCard,
			expect:     []komoju.PaymentType{komoju.PaymentTypeCreditCard},
		},
		{
			name:       "konbini",
			methodType: PaymentMethodTypeKonbini,
			expect:     []komoju.PaymentType{komoju.PaymentTypeKonbini},
		},
		{
			name:       "bank transfer",
			methodType: PaymentMethodTypeBankTransfer,
			expect:     []komoju.PaymentType{komoju.PaymentTypeBankTransfer},
		},
		{
			name:       "paypay",
			methodType: PaymentMethodTypePayPay,
			expect:     []komoju.PaymentType{komoju.PaymentTypePayPay},
		},
		{
			name:       "line pay",
			methodType: PaymentMethodTypeLinePay,
			expect:     []komoju.PaymentType{komoju.PaymentTypeLinePay},
		},
		{
			name:       "merpay",
			methodType: PaymentMethodTypeMerpay,
			expect:     []komoju.PaymentType{komoju.PaymentTypeMerpay},
		},
		{
			name:       "rakuten pay",
			methodType: PaymentMethodTypeRakutenPay,
			expect:     []komoju.PaymentType{komoju.PaymentTypeRakutenPay},
		},
		{
			name:       "au pay",
			methodType: PaymentMethodTypeAUPay,
			expect:     []komoju.PaymentType{komoju.PaymentTypeAUPay},
		},
		{
			name:       "unknown",
			methodType: PaymentMethodTypeUnknown,
			expect:     []komoju.PaymentType{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewKomojuPaymentTypes(tt.methodType)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestPaymentMethodType_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		methodType PaymentMethodType
		expect     string
	}{
		{
			name:       "cash",
			methodType: PaymentMethodTypeCash,
			expect:     "代引支払い",
		},
		{
			name:       "credit card",
			methodType: PaymentMethodTypeCreditCard,
			expect:     "クレジットカード決済",
		},
		{
			name:       "konbini",
			methodType: PaymentMethodTypeKonbini,
			expect:     "コンビニ決済",
		},
		{
			name:       "bank transfer",
			methodType: PaymentMethodTypeBankTransfer,
			expect:     "銀行振込決済",
		},
		{
			name:       "paypay",
			methodType: PaymentMethodTypePayPay,
			expect:     "QR決済（PayPay）",
		},
		{
			name:       "line pay",
			methodType: PaymentMethodTypeLinePay,
			expect:     "QR決済（LINE Pay）",
		},
		{
			name:       "merpay",
			methodType: PaymentMethodTypeMerpay,
			expect:     "QR決済（メルペイ）",
		},
		{
			name:       "rakuten pay",
			methodType: PaymentMethodTypeRakutenPay,
			expect:     "QR決済（楽天ペイ）",
		},
		{
			name:       "au pay",
			methodType: PaymentMethodTypeAUPay,
			expect:     "QR決済（au PAY）",
		},
		{
			name:       "none",
			methodType: PaymentMethodTypeNone,
			expect:     "決済なし",
		},
		{
			name:       "unknown",
			methodType: PaymentMethodTypeUnknown,
			expect:     "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.methodType.String())
		})
	}
}

func TestProductOrderPayment(t *testing.T) {
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
		params    *NewProductOrderPaymentParams
		expect    *OrderPayment
		expectErr error
	}{
		{
			name: "success with shipping free",
			params: &NewProductOrderPaymentParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				MethodType: PaymentMethodTypeCreditCard,
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
			expect: &OrderPayment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            PaymentStatusPending,
				TransactionID:     "",
				MethodType:        PaymentMethodTypeCreditCard,
				Subtotal:          4460,
				Discount:          0,
				ShippingFee:       0,
				Tax:               405,
				Total:             4460,
			},
			expectErr: nil,
		},
		{
			name: "success with discount",
			params: &NewProductOrderPaymentParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				MethodType: PaymentMethodTypeCreditCard,
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
			expect: &OrderPayment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            PaymentStatusPending,
				TransactionID:     "",
				MethodType:        PaymentMethodTypeCreditCard,
				Subtotal:          4460,
				Discount:          446,
				ShippingFee:       500,
				Tax:               410,
				Total:             4514,
			},
			expectErr: nil,
		},
		{
			name: "failed to calc total price",
			params: &NewProductOrderPaymentParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				MethodType: PaymentMethodTypeCreditCard,
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
			params: &NewProductOrderPaymentParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				MethodType: PaymentMethodTypeCreditCard,
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
		{
			name: "empty address",
			params: &NewProductOrderPaymentParams{
				OrderID:    "order-id",
				Address:    nil,
				MethodType: PaymentMethodTypeCreditCard,
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
			expectErr: errNotFoundAddress,
		},
		{
			name: "invalid prefecture code",
			params: &NewProductOrderPaymentParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: -1,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				MethodType: PaymentMethodTypeCreditCard,
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
			expectErr: codes.ErrUnknownPrefecture,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewProductOrderPayment(tt.params)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceOrderPayment(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		params *NewExperienceOrderPaymentParams
		expect *OrderPayment
		hasErr bool
	}{
		{
			name: "success",
			params: &NewExperienceOrderPaymentParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				MethodType: PaymentMethodTypeCreditCard,
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
					MediaJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/thumbnail.png","isThumbnail":true}]`)),
					RecommendedPoints: []string{
						"ポイント1",
						"ポイント2",
					},
					RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
					PromotionVideoURL:     "http://example.com/promotion.mp4",
					Duration:              60,
					Direction:             "彦根駅から徒歩10分",
					BusinessOpenTime:      "1000",
					BusinessCloseTime:     "1800",
					HostPostalCode:        "5220061",
					HostPrefecture:        "滋賀県",
					HostPrefectureCode:    25,
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
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
			expect: &OrderPayment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            PaymentStatusPending,
				TransactionID:     "",
				MethodType:        PaymentMethodTypeCreditCard,
				Subtotal:          2000,
				Discount:          0,
				ShippingFee:       0,
				Tax:               181,
				Total:             2000,
			},
			hasErr: false,
		},
		{
			name: "empty address",
			params: &NewExperienceOrderPaymentParams{
				OrderID:    "order-id",
				Address:    nil,
				MethodType: PaymentMethodTypeCreditCard,
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
					MediaJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/thumbnail.png","isThumbnail":true}]`)),
					RecommendedPoints: []string{
						"ポイント1",
						"ポイント2",
					},
					RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
					PromotionVideoURL:     "http://example.com/promotion.mp4",
					Duration:              60,
					Direction:             "彦根駅から徒歩10分",
					BusinessOpenTime:      "1000",
					BusinessCloseTime:     "1800",
					HostPostalCode:        "5220061",
					HostPrefecture:        "滋賀県",
					HostPrefectureCode:    25,
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
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
			expect: nil,
			hasErr: true,
		},
		{
			name: "invalid prefecture code",
			params: &NewExperienceOrderPaymentParams{
				OrderID: "order-id",
				Address: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: -1,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				MethodType: PaymentMethodTypeCreditCard,
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
					MediaJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/thumbnail.png","isThumbnail":true}]`)),
					RecommendedPoints: []string{
						"ポイント1",
						"ポイント2",
					},
					RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
					PromotionVideoURL:     "http://example.com/promotion.mp4",
					Duration:              60,
					Direction:             "彦根駅から徒歩10分",
					BusinessOpenTime:      "1000",
					BusinessCloseTime:     "1800",
					HostPostalCode:        "5220061",
					HostPrefecture:        "滋賀県",
					HostPrefectureCode:    25,
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
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
			expect: nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewExperienceOrderPayment(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestOrderPayment_IsCompleted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status PaymentStatus
		expect bool
	}{
		{
			name:   "pending",
			status: PaymentStatusPending,
			expect: false,
		},
		{
			name:   "authorized",
			status: PaymentStatusAuthorized,
			expect: false,
		},
		{
			name:   "captured",
			status: PaymentStatusCaptured,
			expect: true,
		},
		{
			name:   "canceled",
			status: PaymentStatusCanceled,
			expect: true,
		},
		{
			name:   "refunded",
			status: PaymentStatusRefunded,
			expect: true,
		},
		{
			name:   "failed",
			status: PaymentStatusFailed,
			expect: true,
		},
		{
			name:   "unknown",
			status: PaymentStatusUnknown,
			expect: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			payment := &OrderPayment{Status: tt.status}
			assert.Equal(t, tt.expect, payment.IsCompleted())
		})
	}
}

func TestOrderPayment_IsCanceled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status PaymentStatus
		expect bool
	}{
		{
			name:   "pending",
			status: PaymentStatusPending,
			expect: false,
		},
		{
			name:   "authorized",
			status: PaymentStatusAuthorized,
			expect: false,
		},
		{
			name:   "captured",
			status: PaymentStatusCaptured,
			expect: false,
		},
		{
			name:   "canceled",
			status: PaymentStatusCanceled,
			expect: true,
		},
		{
			name:   "refunded",
			status: PaymentStatusRefunded,
			expect: true,
		},
		{
			name:   "failed",
			status: PaymentStatusFailed,
			expect: false,
		},
		{
			name:   "unknown",
			status: PaymentStatusUnknown,
			expect: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			payment := &OrderPayment{Status: tt.status}
			assert.Equal(t, tt.expect, payment.IsCanceled())
		})
	}
}

func TestOrderPayment_SetTransactionID(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name          string
		payment       *OrderPayment
		transactionID string
		now           time.Time
		expect        *OrderPayment
	}{
		{
			name: "success",
			payment: &OrderPayment{
				Total: 1000,
			},
			transactionID: "transaction-id",
			now:           now,
			expect: &OrderPayment{
				Total:         1000,
				TransactionID: "transaction-id",
				OrderedAt:     now,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.payment.SetTransactionID(tt.transactionID, tt.now)
			assert.Equal(t, tt.expect, tt.payment)
		})
	}
}

func TestOrderPayments_AddressRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		payments OrderPayments
		expect   []int64
	}{
		{
			name: "success",
			payments: OrderPayments{
				{
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id01",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          1980,
					Discount:          0,
					ShippingFee:       550,
					Tax:               253,
					Total:             2530,
				},
				{
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id02",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          3000,
					Discount:          0,
					ShippingFee:       0,
					Tax:               300,
					Total:             3000,
				},
			},
			expect: []int64{1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.payments.AddressRevisionIDs())
		})
	}
}

func TestOrderPayments_MapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		payments OrderPayments
		expect   map[string]*OrderPayment
	}{
		{
			name: "success",
			payments: OrderPayments{
				{
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id01",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          1980,
					Discount:          0,
					ShippingFee:       550,
					Tax:               253,
					Total:             2530,
				},
				{
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id02",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          3000,
					Discount:          0,
					ShippingFee:       0,
					Tax:               300,
					Total:             3000,
				},
			},
			expect: map[string]*OrderPayment{
				"order-id01": {
					OrderID:           "order-id01",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id01",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          1980,
					Discount:          0,
					ShippingFee:       550,
					Tax:               253,
					Total:             2530,
				},
				"order-id02": {
					OrderID:           "order-id02",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id02",
					Status:            PaymentStatusCaptured,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          3000,
					Discount:          0,
					ShippingFee:       0,
					Tax:               300,
					Total:             3000,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.payments.MapByOrderID())
		})
	}
}
